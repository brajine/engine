package main

// Rebuild json access methods for all structs in file
// easyjson -all <file>.go

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"../engine/data"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Loggers
var (
	Trace *log.Logger
	Error *log.Logger
)

// WebSockets
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// TCPServer is to manage tcp listener service
type TCPServer struct {
	addr string
}

// MTListener implements Metatrader tcp listener service
var MTListener TCPServer

// Storage is primary data structure
var Storage data.MainStorage

// statsAPIHandler is a handler for server state api
func statsAPIHandler(w http.ResponseWriter, r *http.Request) {
	st := Storage.ExportState()
	bt, _ := json.Marshal(st)
	w.Write(bt)
}

// restAPIHandler is serving REST API calls
func restAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// Check if page exists
	acc := Storage.PageExist(page)
	if acc == nil {
		http.NotFound(w, r)
		return
	}

	w.Write(acc.ToJSON())
}

// wsAPIHandler is serving WebSocket connections
func wsAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// Check if page exists
	acc := Storage.PageExist(page)
	if acc == nil {
		fmt.Fprintln(w, "PAGE DOES NOT EXIST")
		return
	}

	// Upgrade http connection to that of WebSockets type
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			Error.Println(err)
		}
		return
	}

	// Add new connection to page viewers pool, and send him update
	acc.AddView(conn)
}

// Run our MetaTrader listener service
func (mt *TCPServer) Run(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		Error.Println("Can not create tcp listener.", err)
		return
	}

	for {
		// Wait for connection
		conn, err := ln.Accept()
		if err != nil {
			Error.Println("Error accepting connection, err #", err)
			continue
		}

		// Proceed with connection
		go mt.Handle(conn)
	}
}

// Handle MetaTrader connection
func (mt *TCPServer) Handle(conn net.Conn) {
	var connPage string
	Trace.Println("Accepted connection from", conn.RemoteAddr())
	defer func() {
		if connPage != "" {
			Trace.Println("Client disconnected: \"" + connPage + "\"")
		} else {
			Trace.Println("Disconnected", conn.RemoteAddr())
		}
		Storage.RemoveClient(conn)
		conn.Close()
	}()

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	// Read first message
	var msg data.TradesMsg
	if err := dec.Decode(&msg); err != nil {
		Error.Println("Message decode error:", err)
		return
	}

	// Validate first message
	if err := msg.Validate(); err != nil {
		Error.Println("Message validate error:", err)
		resp := data.ResponseMsg{
			Error: err.Error(),
		}
		writeMessage(enc, resp)
		return
	}

	// Page address already exists
	if Storage.PageExist(msg.Page) != nil {
		resp := data.ResponseMsg{
			Error: "Page address " + string(msg.Page) + " is already in use. Please try another.",
		}
		writeMessage(enc, resp)
		return
	}

	// Add new client
	acc := Storage.AddClient(&msg, conn)
	writeOkMessage(enc)
	connPage = acc.Page()
	Trace.Println("New client registered: \"" + connPage + "\"")

	// Messaging loop
	var num, sec int
	timeout := 5 * time.Second
	if msg.UpdateFreq != "second" {
		timeout = 1*time.Minute + 5*time.Second
	}
	for {
		// Check for messages frequency
		if checkFrequency(&num, &sec) > 5 {
			writeMessage(enc, data.ResponseMsg{
				Error: "Update frequency limit exceeded",
			})
			return
		}

		// Read message
		conn.SetReadDeadline(time.Now().Add(timeout))
		if err := dec.Decode(&msg); err != nil {
			Error.Println("Message decode error:", err)
			return
		}

		// Validate message
		if err := msg.Validate(); err != nil {
			writeMessage(enc, data.ResponseMsg{
				Error: err.Error(),
			})
			continue
		}

		// Process message
		acc.Update(&msg)
		acc.BCastUpdate()

		// Send Ok
		writeOkMessage(enc)
	}
}

// Write a response to Metatrader client
func writeMessage(enc *gob.Encoder, resp data.ResponseMsg) error {
	if resp.Error != "" {
		Error.Println(resp.Error)
	}
	return enc.Encode(resp)
}

func writeOkMessage(enc *gob.Encoder) error {
	return enc.Encode(data.ResponseMsg{})
}

// Return number of messages per second
func checkFrequency(num, sec *int) int {
	// Message frequency limit
	if *sec != time.Now().Second() {
		*sec = time.Now().Second()
		*num = 0
	} else {
		*num++
	}
	return *num
}

func main() {
	// Log subsystem: ioutil.Discard for nuldev
	Trace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)

	// Main MT data storage use Token as a Key
	Storage.Initialize()

	// TCP server on 8181 to listen MT clients
	go MTListener.Run(":8181")
	Trace.Println("MetaTrader listener is up and running on :8181")

	// HTTP server to serve data
	r := mux.NewRouter()
	r.HandleFunc("/api/stats", statsAPIHandler)
	r.HandleFunc("/api/rest/{page}", restAPIHandler)
	r.HandleFunc("/api/ws/{page}", wsAPIHandler)

	// Running GO app as a service
	// https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/

	// setup signal catching
	sigs := make(chan os.Signal, 1)

	// catch all signals since not explicitly listing
	signal.Notify(sigs)

	// method invoked upon seeing signal
	go func() {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s", s)
		os.Exit(1)
	}()

	err := http.ListenAndServeTLS(":8182", "/etc/letsencrypt/live/metatrader.live/cert.pem", "/etc/letsencrypt/live/metatrader.live/privkey.pem", r)
	if err != nil {
		fmt.Println(err)
	}
}
