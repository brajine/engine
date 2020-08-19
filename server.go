package main

// Rebuild json access methods for all structs in file
// easyjson -all <file>.go

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
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

// apiHandler is a handler for server api
func apiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	if page == "stats" {
		type stats struct {
			Online    int `json:"online"`
			Online24h int `json:"online24h"`
		}
		st := stats{
			Online:    Storage.ClientsNum(),
			Online24h: 5,
		}
		bt, _ := json.Marshal(st)
		w.Write(bt)
		return
	}

	// Not found
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// ViewHandler is a handler for Account /view page
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// Check if page exists
	if acc := Storage.PageExist(page); acc == nil {
		fmt.Fprintln(w, "PAGE DOES NOT EXIST")
		return
	}

	// Send view template
	tmpl := template.Must(template.ParseFiles("templates/view.htm"))
	tmpl.Execute(w, &struct{ Page string }{page})
}

// WsHandler is serving WebSocket connections
func WsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	// Check if page exists
	acc := Storage.PageExist(page)
	if acc == nil {
		fmt.Fprintln(w, "PAGE DOES NOT EXIST")
		return
	}

	// Upgrade http connection to WebSockets type
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
		Trace.Println("Accepted connection from", conn.RemoteAddr())
		go mt.Handle(conn)
	}
}

// Handle MetaTrader connection
func (mt *TCPServer) Handle(conn net.Conn) {
	defer func() {
		Trace.Printf("Disconnected (%d online): %s", Storage.ClientsNum(), conn.RemoteAddr())
		Storage.RemoveClient(conn)
		conn.Close()
	}()

	// Limit maximum amount of data to read
	var bb bytes.Buffer
	dec := gob.NewDecoder(&bb)
	enc := gob.NewEncoder(conn)

	// Read first message and allocate new account
	var resp data.ResponseMsg
	msg, err := readMessage(conn, dec, &bb, true)
	if err == nil {
		if exist := Storage.PageExist(msg.Page); exist != nil {
			resp.Error = "Page address " + string(msg.Page) + " is already in use. Please try another."
		}
	} else {
		resp.Error = err.Error()
	}
	err = writeMessage(enc, resp)
	if resp.Error != "" || err != nil {
		return
	}
	acc := Storage.AddClient(msg, conn)

	// Messaging loop
	var num, sec int
	for {
		// Read & process message, "page" value is ignored here
		if msg, err = readMessage(conn, dec, &bb, false); err == nil {
			acc.Update(msg)
			acc.BCastUpdate()
		}

		// Check for messages frequency
		if msgFrequency(&num, &sec) > 3 {
			err = errors.New("Update frequency limit exceeded")
		}

		// Send back response
		if err != nil {
			resp.Error = err.Error()
			Error.Println(err)
		}
		if err = writeMessage(enc, resp); err != nil {
			Error.Println("writeMessage():", err, conn.RemoteAddr())
			return
		}
		if resp.Error != "" {
			return
		}
	}
}

// Read, decode and validate new message from Metatrader connection
func readMessage(conn net.Conn, dec *gob.Decoder, bb *bytes.Buffer, header bool) (*data.TradesMsg, error) {
	msg := &data.TradesMsg{}
	bb.Reset()
	readsz := 0

	if header {
		// Read message headers
		for i := 0; i < 3; i++ {
			if err := readChunk(conn, bb, &readsz); err != nil {
				return nil, err
			}
		}
	}

	// Read body
	if err := readChunk(conn, bb, &readsz); err != nil {
		return nil, err
	}

	var err error
	if err = dec.Decode(msg); err == nil {
		if err = msg.Validate(); err == nil {
			return msg, nil
		}
		Error.Println("Validation Error:", err)
	} else {
		Error.Println("Decode Error:", err)
		err = errors.New("Message is not valid")
	}

	return nil, err
}

func readChunk(conn net.Conn, bb *bytes.Buffer, readsz *int) error {
	var tmp [3]byte
	var buffer [3]byte
	t := tmp[:]
	buf := buffer[:]

	if _, err := conn.Read(buf); err == nil {
		bb.Write(buf)
		r := bytes.NewReader(buf)
		if msgLen, w, err := decodeUintReader(r, t); err == nil {
			if *readsz+w+int(msgLen) <= data.MaxMsgSize {
				io.CopyN(bb, conn, int64(msgLen-uint64((3-w))))
				*readsz += w + int(msgLen)
			} else {
				return errors.New("Message is too long")
			}
		} else {
			return errors.New("Can't determine message length")
		}
	} else {
		return errors.New("Message is not valid")
	}

	return nil
}

// decodeUintReader reads an encoded unsigned integer from an io.Reader.
// Used only by the Decoder to read the message length.
func decodeUintReader(r io.Reader, buf []byte) (x uint64, width int, err error) {
	n, err := io.ReadFull(r, buf[0:1])
	if n == 0 {
		err = errors.New("Message is empty")
		return
	}
	if buf[0] <= 0x7f {
		// If a number < 128 then it is single-byte coded
		return uint64(buf[0]), 1, nil
	}

	// FF 82 = 130
	// FE 01 00 = 256
	sz := ^buf[0] + 1
	// FIX for server scan panics
	if int(sz)+1 > len(buf) {
		return 0, 0, errors.New("Wrong message length")
	}
	width, err = io.ReadFull(r, buf[1:sz+1])
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return
	}
	// Could check that the high byte is zero but it's not worth it.
	for _, b := range buf[1 : width+1] {
		x = x<<8 | uint64(b)
	}
	width++ // +1 for length byte
	return
}

// Write a response to Metatrader client
func writeMessage(enc *gob.Encoder, resp data.ResponseMsg) error {
	return enc.Encode(resp)
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
	r.HandleFunc("/api/{page}", apiHandler)
	r.HandleFunc("/accounts/{page}/view", ViewHandler)
	r.HandleFunc("/accounts/{page}/ws", WsHandler)

	// Running GO app as a service
	// https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/

	// setup signal catching
	sigs := make(chan os.Signal, 1)

	// catch all signals since not explicitly listing
	signal.Notify(sigs)
	// signal.Notify(sigs,syscall.SIGQUIT)

	// method invoked upon seeing signal
	go func() {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s", s)
		os.Exit(1)
	}()

	err := http.ListenAndServe(":8182", r)
	if err != nil {
		fmt.Println(err)
	}
}

// Return number of messages per second
func msgFrequency(num, sec *int) int {
	// Message frequency limit
	if *sec != time.Now().Second() {
		*sec = time.Now().Second()
		*num = 1
	} else {
		*num++
	}
	return *num
}
