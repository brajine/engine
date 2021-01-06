package metatrader

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Account represent connected MetaTrader Client
// broker keep WebSocket clients array
type Account struct {
	Message
	broker *BrokerFactory
}

// NewAccount ...
func NewAccount(msg *Message, log *zap.SugaredLogger) *Account {
	acc := new(Account) // Promoted fields may not initialized in list
	acc.broker = NewBroker(log)
	acc.Page = msg.Page
	acc.Started = time.Now()
	acc.UpdateFreq = msg.UpdateFreq
	acc.ClientVersion = msg.ClientVersion
	acc.Orders = make(map[OrderTicket]Order)
	acc.update(msg)
	return acc
}

// close all viewers and destroy account
func (a *Account) close() {
	a.Orders = nil
	a.broker.Stop()
}

// Update update existing MT Client account with new data
func (a *Account) update(upd *Message) {
	// Update Account data
	a.updateInfo(upd)
	a.Updated = time.Now()

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
func (a *Account) AddViewer(viewer *websocket.Conn) {
	a.broker.AddViewer(viewer)
}

// RemoveViewer removes a connection from page vewers pool
func (a *Account) RemoveViewer(viewer *websocket.Conn) {
	a.broker.RemoveViewer(viewer)
}

// SendUpdateToViewer ...
func (a *Account) SendUpdateToViewer(viewer *websocket.Conn) {
	msg, _ := a.ToJSON()
	a.broker.SendMessageToViewer(viewer, msg)
}

// SendUpdateToAllViewers ...
func (a *Account) SendUpdateToAllViewers() {
	msg, _ := a.ToJSON()
	a.broker.SendMessage(msg)
}
