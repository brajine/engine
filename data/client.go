package data

import (
	"sort"
	"strconv"
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

	// Set Page & Started on new connection
	client.acc.Page = msg.Page
	client.acc.Started = time.Now()
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

// Started return sting with account initialization time
func (client *TradesAccount) Started() string {
	return client.acc.Started.Format("2006-01-02 15:04:05")
}

// ToJSON create json message for WebSockets
func (client *TradesAccount) ToJSON() []byte {
	var w jwriter.Writer
	client.acc.MarshalEasyJSON(&w)
	b, _ := w.BuildBytes()
	return b
}

// MarshalEasyJSON overwrites EaseJSON standart handler
// This makes orders array [] instead of map {}
func (in TradesMsg) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	first := true
	_ = first
	if in.UpdateFreq != "" {
		const prefix string = ",\"updatefreq\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.UpdateFreq))
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.Login != "" {
		const prefix string = ",\"login\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Login))
	}
	if in.Server != "" {
		const prefix string = ",\"server\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Server))
	}
	if in.Company != "" {
		const prefix string = ",\"company\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Company))
	}
	if in.Balance != "" {
		const prefix string = ",\"balance\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Balance))
	}
	if in.Equity != "" {
		const prefix string = ",\"equity\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Equity))
	}
	if in.Margin != "" {
		const prefix string = ",\"margin\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Margin))
	}
	if in.FreeMargin != "" {
		const prefix string = ",\"freemargin\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.FreeMargin))
	}
	if in.MarginLevel != "" {
		const prefix string = ",\"marginlevel\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.MarginLevel))
	}
	if in.ProfitTotal != "" {
		const prefix string = ",\"profittotal\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ProfitTotal))
	}
	const prefix string = ",\"orderscount\":"
	if first {
		first = false
		out.RawString(prefix[1:])
	} else {
		out.RawString(prefix)
	}
	out.String(strconv.Itoa(len(in.Orders)))
	{
		const prefix string = ",\"orders\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Orders == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('[')
			v2First := true
			var sorted []string

			// Sort output by ticket
			for v2Name := range in.Orders {
				sorted = append(sorted, v2Name)
			}
			sort.Sort(sort.StringSlice(sorted))

			for _, v2Name := range sorted {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				easyMarshallOrder(out, v2Name, in.Orders[v2Name])
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// easyMarshallOrder overwrites marshalling of a single order
func easyMarshallOrder(out *jwriter.Writer, ticket string, in OrderType) {
	out.RawByte('{')
	first := true
	_ = first
	if ticket != "" {
		const prefix string = ",\"ticket\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(ticket))
	}
	if in.Symbol != "" {
		const prefix string = ",\"symbol\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Symbol))
	}
	if in.TimeOpen != "" {
		const prefix string = ",\"timeopen\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.TimeOpen))
	}
	if in.Type != "" {
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	if in.InitVolume != "" {
		const prefix string = ",\"initvolume\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.InitVolume))
	}
	if in.CurVolume != "" {
		const prefix string = ",\"curvolume\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.CurVolume))
	}
	if in.PriceOpen != "" {
		const prefix string = ",\"priceopen\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PriceOpen))
	}
	if in.SL != "" {
		const prefix string = ",\"sl\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SL))
	}
	if in.TP != "" {
		const prefix string = ",\"tp\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.TP))
	}
	if in.Swap != "" {
		const prefix string = ",\"swap\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Swap))
	}
	if in.PriceSL != "" {
		const prefix string = ",\"pricesl\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PriceSL))
	}
	if in.Profit != "" {
		const prefix string = ",\"profit\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Profit))
	}
	out.RawByte('}')
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
