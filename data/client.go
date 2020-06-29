package data

import (
	"time"

	"github.com/gorilla/websocket"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// TradesAccount represent connected MetaTrader Client
// views keep WebSocket clients array
type TradesAccount struct {
	acc   TradesMsg
	views map[*websocket.Conn]bool
}

// Init newly created Account
func (client *TradesAccount) Init(msg *TradesMsg) {
	client.views = make(map[*websocket.Conn]bool)
	client.acc.Orders = make(map[string]OrderType)
	client.Update(msg)

	// Only allow Page set upon first connect
	client.acc.Page = msg.Page
}

// Update update existing MT Client account with new data
func (client *TradesAccount) Update(upd *TradesMsg) {
	// Update Account data
	client.updateInfo(upd)
	client.acc.Updated = time.Now()

	// Remove closed orders
	// Metatrader sends entire ticket array in every message
	// If ticket array in new message doesn't contains one of Storage tickets, this means order was closed and should be removed from Storage
	for tick := range client.acc.Orders {
		if _, ok := upd.Orders[tick]; !ok {
			delete(client.acc.Orders, tick)
		}
	}

	// Add new && Update existing orders
	for tick, order := range upd.Orders {
		if ord, ok := client.acc.Orders[tick]; !ok {
			client.acc.Orders[tick] = order
		} else {
			ord.updateWith(order)
			client.acc.Orders[tick] = ord
		}
	}
}

// updateInfo updates only account parameters (not Orders)
func (client *TradesAccount) updateInfo(upd *TradesMsg) {
	if upd.ClientVersion != "" {
		client.acc.ClientVersion = upd.ClientVersion
	}

	if upd.UpdateFreq != "" {
		client.acc.UpdateFreq = upd.UpdateFreq
	}

	if upd.Name != "" {
		client.acc.Name = upd.Name
	}

	if upd.Login != "" {
		client.acc.Login = upd.Login
	}

	if upd.Server != "" {
		client.acc.Server = upd.Server
	}

	if upd.Company != "" {
		client.acc.Company = upd.Company
	}

	if upd.Balance != "" {
		client.acc.Balance = upd.Balance
	}

	if upd.Equity != "" {
		client.acc.Equity = upd.Equity
	}

	if upd.Margin != "" {
		client.acc.Margin = upd.Margin
	}

	if upd.FreeMargin != "" {
		client.acc.FreeMargin = upd.FreeMargin
	}

	if upd.MarginLevel != "" {
		client.acc.MarginLevel = upd.MarginLevel
	}

	if upd.ProfitTotal != "" {
		client.acc.ProfitTotal = upd.ProfitTotal
	}
}

// UpdateWith order a with data from order b
func (a *OrderType) updateWith(b OrderType) {
	if b.Symbol != "" {
		a.Symbol = b.Symbol
	}
	if b.TimeOpen != "" {
		a.TimeOpen = b.TimeOpen
	}
	if b.Type != "" {
		a.Type = b.Type
	}
	if b.InitVolume != "" {
		a.InitVolume = b.InitVolume
	}
	if b.CurVolume != "" {
		a.CurVolume = b.CurVolume
	}
	if b.PriceOpen != "" {
		a.PriceOpen = b.PriceOpen
	}
	if b.SL != "" {
		a.SL = b.SL
	}
	if b.TP != "" {
		a.TP = b.TP
	}
	if b.Swap != "" {
		a.Swap = b.Swap
	}
	if b.PriceSL != "" {
		a.PriceSL = b.PriceSL
	}
	if b.Profit != "" {
		a.Profit = b.Profit
	}
}

// Page return account page
func (client *TradesAccount) Page() string {
	return string(client.acc.Page)
}

// Updated return sting with last update time
func (client *TradesAccount) Updated() string {
	return client.acc.Updated.Format("2006-01-02 15:04:05")
}

// ToJSON create json message for WebSockets
func (client *TradesAccount) ToJSON() []byte {
	var w jwriter.Writer
	client.acc.MarshalEasyJSON(&w)
	b, _ := w.BuildBytes()
	return b
}

// func (e *Email) UnmarshalJSON(data []byte) error {
// 	l := jlexer.Lexer{Data: data}
// 	e.UnmarshalEasyJSON(&l)
// 	return l.Error()
// }

// AddView add new Websocket.Conn to the page viewers pool
// Also send him Update message
func (client *TradesAccount) AddView(view *websocket.Conn) {
	if msg := client.ToJSON(); msg != nil {
		if err := view.WriteMessage(websocket.TextMessage, msg); err == nil {
			client.views[view] = true
		}
	}
}

// RemoveView removes a connection from page vewers pool
func (client *TradesAccount) RemoveView(view *websocket.Conn) {
	view.Close()
	delete(client.views, view)
}

// BCastUpdate broadcating update message to all page vewers
func (client *TradesAccount) BCastUpdate() {
	if msg := client.ToJSON(); msg != nil {
		for view := range client.views {
			if err := view.WriteMessage(websocket.TextMessage, msg); err != nil {
				client.RemoveView(view)
			}
		}
	}
}
