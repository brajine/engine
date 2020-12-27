package metatrader

// Rebuild json access methods for all structs in file
// easyjson -all <file>.go

import (
	"encoding/gob"
	"errors"
	"net"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// MaxFreeOrders ...
const (
	MaxFreeOrders int = 30   // orders limit for free accounts
	MaxMsgSize    int = 2343 // maximum theretical incoming message
	// MaxAwaitingSeconds int = 3    // awaiting an Account update. Drop the connection if exceeded
	// MaxUpdateRate      int = 5    // Max updates per second. Disconnect if exceeded
)

// Factory is to manage Metatrader service
type Factory struct {
	addr     string
	accounts map[string]*Account // page as a key
	log      *zap.SugaredLogger
	sync.RWMutex
}

// New create Metatrader factory
func New(addr string, log *zap.SugaredLogger) *Factory {
	return &Factory{
		log:      log,
		addr:     addr,
		accounts: make(map[string]*Account),
	}
}

// Run our MetaTrader listener service
func (f *Factory) Run() {
	go f.startAPIServer(":8182")

	ln, err := net.Listen("tcp", f.addr)
	if err != nil {
		f.log.Error("Can not create tcp listener.", err)
		return
	}

	for {
		// Wait for connection
		conn, err := ln.Accept()
		if err != nil {
			f.log.Error("Error accepting connection, err #", err)
			continue
		}

		// Proceed with connection
		go f.Handle(conn)
	}
}

// Handle MetaTrader connection
func (f *Factory) Handle(conn net.Conn) {
	f.log.Info("Accepted connection from", conn.RemoteAddr())

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	f.ProcessMessages(enc, dec, conn.RemoteAddr().String())

	conn.Close()
}

// ProcessMessages from metatrader connection
func (f *Factory) ProcessMessages(enc *gob.Encoder, dec *gob.Decoder, logaddr string) {
	defer func() {
		f.log.Info("Connection is closed (", logaddr, ")")
	}()

	// Messaging loop
	var page string
	var acc *Account
	for {
		// Decode new message
		msg := new(Message)
		if err := dec.Decode(msg); err != nil {
			f.writeErrorMessage(enc, logaddr, page, "Failed to decode a message: "+err.Error())
			return
		}

		// Validate message
		if err := msg.Validate(); err != nil {
			f.writeErrorMessage(enc, logaddr, page, "Message is not valid: "+err.Error())
			return
		}

		// All Subsequent messages except first one
		if page != "" {
			acc.update(msg)
			acc.bcastUpdate()
			f.writeOkMessage(enc, "")
			continue
		}

		// First message
		if err := f.firstMessageCheck(msg); err != nil {
			f.writeErrorMessage(enc, logaddr, page, err.Error())
			return
		}
		page = msg.Page
		acc = f.createAccount(msg)
		defer func() {
			f.removeAccount(page)
			f.log.Info("Account disconnected: " + page + "")
		}()

		f.writeOkMessage(enc, "New account registered: "+page+"")
	}
}

// First message should contain mandatory fields - Page, UpdateFreq
func (f *Factory) firstMessageCheck(msg *Message) error {
	if msg.Page == "" {
		return errors.New("Page address is not provided")
	}
	if f.PageExist(msg.Page) != nil {
		return errors.New("Page address " + msg.Page + " is already in use")
	}
	freq := strings.ToLower(msg.UpdateFreq)
	if freq != "second" && freq != "minute" {
		return errors.New("Update frequency " + freq + " is not valid")
	}
	return nil
}

func (f *Factory) createAccount(msg *Message) *Account {
	acc := NewAccount(msg)
	f.Lock()
	f.accounts[msg.Page] = acc
	f.Unlock()

	return acc
}

func (f *Factory) removeAccount(page string) {
	f.Lock()
	defer f.Unlock()

	if acc := f.PageExist(page); acc != nil {
		delete(f.accounts, page)
		acc.close()
	}
}

// ExportState return slice of account pointers for Stats page
func (f *Factory) exportState() *StateData {
	f.RLock()
	defer f.RUnlock()

	st := StateData{Online: len(f.accounts)}
	for _, acc := range f.accounts {
		started := acc.Started.Format("2006-01-02 15:04:05")
		entry := StateEntry{
			Page:       acc.Page,
			Started:    started,
			UpdateFreq: acc.UpdateFreq,
		}
		st.Accounts = append(st.Accounts, entry)
	}
	return &st
}

// PageExist return account with page specified OR nil
func (f *Factory) PageExist(page string) *Account {
	if acc, ok := f.accounts[page]; ok {
		return acc
	}
	return nil
}

// NumAccounts ...
func (f *Factory) NumAccounts() int {
	return len(f.accounts)
}

// Write a response to Metatrader client
// If page is set, then output console message also
func (f *Factory) writeErrorMessage(enc *gob.Encoder, addr, page, text string) {
	f.log.Error(text, ". (", addr, ", ", page, ")")
	err := enc.Encode(ResponseMsg{
		Error: text,
	})
	if err != nil {
		f.log.Error("Failed to encode response message: ", err)
	}
}

func (f *Factory) writeOkMessage(enc *gob.Encoder, str string) error {
	if str != "" {
		f.log.Info(str)
	}
	return enc.Encode(ResponseMsg{})
}
