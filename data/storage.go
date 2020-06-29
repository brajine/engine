package data

import (
	"net"
	"sync"
)

// MaxFreeOrders limit orders number for free accounts
// Don't forget to modify TradesMsg.Orders `validate max` amount
const MaxFreeOrders int = 15

// MaxMsgSize is a maximum theretical incoming metatrader message
const MaxMsgSize int = 2343 // check MT_Update.json for details

// MainStorage type describe main storage for all data
// clients map reference net.Conn -> Client account
// pages map reference page -> Client account
type MainStorage struct {
	sync.RWMutex // Used with general account operations
	clients      map[net.Conn]*TradesAccount
	pages        map[string]*TradesAccount
}

// Initialize main storage
func (s *MainStorage) Initialize() {
	s.clients = make(map[net.Conn]*TradesAccount)
	s.pages = make(map[string]*TradesAccount)
}

// ClientsNum return actual number of accounts in MainStorage
func (s *MainStorage) ClientsNum() int {
	return len(s.clients)
}

// AddClient adds new client to the main Storage
func (s *MainStorage) AddClient(msg *TradesMsg, conn net.Conn) *TradesAccount {
	s.Lock()
	defer s.Unlock()

	var acc TradesAccount
	acc.Init(msg)
	s.clients[conn] = &acc
	s.pages[string(msg.Page)] = &acc
	return &acc
}

// RemoveClient from storage
func (s *MainStorage) RemoveClient(conn net.Conn) {
	s.Lock()
	defer s.Unlock()

	if c, ok := s.clients[conn]; ok {
		// Close all WebSocket clients
		for v := range c.views {
			v.Close()
		}
		conn.Close()
		delete(s.pages, c.Page())
		delete(s.clients, conn)
	}
}

// ExportAccArray return slice of account pointers for Stats page
func (s *MainStorage) ExportAccArray() (accs []*TradesAccount) {
	s.RLock()
	defer s.RUnlock()

	for _, acc := range s.clients {
		accs = append(accs, acc)
	}
	return accs
}

// PageExist checks if account exist in main Storage
func (s *MainStorage) PageExist(page interface{}) *TradesAccount {
	var p string
	switch page.(type) {
	case string:
		p = page.(string)
	default:
		p = string(page.([]byte))
	}

	if client, ok := s.pages[p]; ok {
		return client
	}
	return nil
}
