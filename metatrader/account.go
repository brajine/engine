package metatrader

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

// Account represent connected MetaTrader Client
// viewers keep WebSocket clients array
type Account struct {
	Message
	viewers map[*websocket.Conn]bool
}

// NewAccount ...
func NewAccount(msg *Message) *Account {
	acc := new(Account) // Promoted fields may not initialized in list
	acc.viewers = make(map[*websocket.Conn]bool)
	acc.Page = msg.Page
	acc.Started = OrderTime(time.Now())
	acc.UpdateFreq = msg.UpdateFreq
	acc.ClientVersion = msg.ClientVersion
	acc.Orders = make(map[OrderTicket]Order)
	acc.update(msg)
	return acc
}

// close all viewers and destroy account
func (a *Account) close() {
	// Close all WebSocket clients
	for v := range a.viewers {
		v.Close()
	}
	a.viewers = nil
	a.Orders = nil
}

// Update update existing MT Client account with new data
func (a *Account) update(upd *Message) {
	// Update Account data
	a.updateInfo(upd)
	a.Updated = OrderTime(time.Now())

	// Remove closed orders
	// Metatrader sends entire ticket array in every message
	// If ticket array in new message doesn't contains one of Storage tickets, this means order was closed and should be removed from Storage
	for tick := range a.Orders {
		if _, ok := upd.Orders[tick]; !ok {
			delete(a.Orders, tick)
		}
	}

	// Add new && Update existing orders
	for tick, order := range upd.Orders {
		if ord, ok := a.Orders[tick]; !ok {
			a.Orders[tick] = order
		} else {
			ord.UpdateWith(order)
			a.Orders[tick] = ord
		}
	}

	a.OrdersCount = len(a.Orders)
}

func (a *Account) updateInfo(upd *Message) {
	if upd.Name != "" {
		a.Name = upd.Name
	}

	if upd.Login != "" {
		a.Login = upd.Login
	}

	if upd.Server != "" {
		a.Server = upd.Server
	}

	if upd.Company != "" {
		a.Company = upd.Company
	}

	if upd.Balance != "" {
		a.Balance = upd.Balance
	}

	if upd.Equity != "" {
		a.Equity = upd.Equity
	}

	if upd.Margin != "" {
		a.Margin = upd.Margin
	}

	if upd.FreeMargin != "" {
		a.FreeMargin = upd.FreeMargin
	}

	if upd.MarginLevel != "" {
		a.MarginLevel = upd.MarginLevel
	}

	if upd.ProfitTotal != "" {
		a.ProfitTotal = upd.ProfitTotal
	}
}

// ToJSON create json message for WebSockets
func (a *Account) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

// Marshall Easy JSON
// var w jwriter.Writer
// a.MarshalEasyJSON(&w)
// b, _ := w.BuildBytes()
// return b

// Unmarshall Easy JSON
// 	l := jlexer.Lexer{Data: data}
// 	e.UnmarshalEasyJSON(&l)
// 	return l.Error()

// AddViewer add new Websocket.Conn to the page viewers pool
// Also send him Update message
func (a *Account) AddViewer(viewer *websocket.Conn) error {
	msg, err := a.ToJSON()
	if err != nil {
		return err
	}
	err = viewer.WriteMessage(websocket.TextMessage, msg)
	if err == nil {
		a.viewers[viewer] = true
	}
	return nil
}

// RemoveViewer removes a connection from page vewers pool
func (a *Account) RemoveViewer(viewer *websocket.Conn) {
	viewer.Close()
	delete(a.viewers, viewer)
}

// BCastUpdate send update message to all page vewers
func (a *Account) bcastUpdate() {
	msg, err := a.ToJSON()
	if err == nil && msg != nil {
		for viewer := range a.viewers {
			if err := viewer.WriteMessage(websocket.TextMessage, msg); err != nil {
				a.RemoveViewer(viewer)
			}
		}
	}
}

// Return timeout in time.Duration
func (a *Account) getTimeout() time.Duration {
	if a.UpdateFreq == "second" {
		return time.Second
	}
	return time.Minute
}
