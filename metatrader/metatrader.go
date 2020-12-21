package metatrader

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// easyjson -no_std_marshalers data/metatrader.go

const (
	orderTypeBuy           int = 0
	orderTypeSell          int = 1
	orderTypeBuyLimit      int = 2
	orderTypeSellLimit     int = 3
	orderTypeBuyStop       int = 4
	orderTypeSellStop      int = 5
	orderTypeBuyStopLimit  int = 6
	orderTypeSellStopLimit int = 7
)

// Message keeps all Metatrader data for each particular client
//easyjson:json
type Message struct {
	Page          string    `json:"page,omitempty"`
	ClientVersion string    `json:"clientversion,omitempty"`
	Started       time.Time `json:"started,omitempty"` // Connection established
	Updated       time.Time `json:"updated,omitempty"` // Last time account was updated
	UpdateFreq    string    `json:"updatefreq,omitempty"`
	Name          string    `json:"name,omitempty"`
	Login         string    `json:"login,omitempty"`
	Server        string    `json:"server,omitempty"`
	Company       string    `json:"company,omitempty"`
	Balance       string    `json:"balance,omitempty"`
	Equity        string    `json:"equity,omitempty"`
	Margin        string    `json:"margin,omitempty"`
	FreeMargin    string    `json:"freemargin,omitempty"`
	MarginLevel   string    `json:"marginlevel,omitempty"`
	ProfitTotal   string    `json:"profittotal,omitempty"`
	OrdersCount   int       `json:"orderscount,omitempty"`
	// Use Ticket as a map key
	Orders map[string]Order `json:"orders,omitempty"`
}

// Order represent one Metatrader order
// Tickets are always sent, other values sent only if not null AND changed since last update
type Order struct {
	Symbol     string `json:"symbol,omitempty"`
	TimeOpen   string `json:"timeopen,omitempty"`
	Type       string `json:"type,omitempty"`
	InitVolume string `json:"initvolume,omitempty"`
	CurVolume  string `json:"curvolume,omitempty"`
	PriceOpen  string `json:"priceopen,omitempty"`
	SL         string `json:"sl,omitempty"`
	TP         string `json:"tp,omitempty"`
	Swap       string `json:"swap,omitempty"`
	PriceSL    string `json:"pricesl,omitempty"`
	Profit     string `json:"profit,omitempty"`
}

// ResponseMsg is sending to MetaTrader
type ResponseMsg struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// Validate incoming Message
func (t *Message) Validate() error {
	if err := validPage(t.Page); err != nil {
		return err
	}
	if err := validString(t.ClientVersion, "ClientVersion"); err != nil {
		return err
	}
	if err := validString(t.UpdateFreq, "UpdateFrequency"); err != nil {
		return err
	}
	if err := validString(t.Name, "Name"); err != nil {
		return err
	}
	if err := validString(t.Login, "Login"); err != nil {
		return err
	}
	if err := validString(t.Server, "Server"); err != nil {
		return err
	}
	if err := validString(t.Company, "Company"); err != nil {
		return err
	}
	if err := validNumber(t.Balance, "Balance"); err != nil {
		return err
	}
	if err := validNumber(t.Equity, "Equity"); err != nil {
		return err
	}
	if err := validNumber(t.Margin, "Margin"); err != nil {
		return err
	}
	if err := validNumber(t.FreeMargin, "FreeMargin"); err != nil {
		return err
	}
	if err := validNumber(t.MarginLevel, "MarginLevel"); err != nil {
		return err
	}
	if err := validNumber(t.ProfitTotal, "ProfitTotal"); err != nil {
		return err
	}
	if len(t.Orders) > MaxFreeOrders {
		return errors.New("Exceeded maximum orders number (" + strconv.Itoa(MaxFreeOrders) + ")")
	}
	for k, v := range t.Orders {
		if err := validNumber(k, "Ticket"); err != nil {
			return nil
		}
		if err := validString(v.Symbol, "Symbol"); err != nil {
			return nil
		}
		if err := validTime(v.TimeOpen, "TimeOpen"); err != nil {
			return nil
		}
		if err := validNumber(v.Type, "Type"); err != nil {
			return nil
		}
		if err := validNumber(v.InitVolume, "Number"); err != nil {
			return nil
		}
		if err := validNumber(v.CurVolume, "CurVolume"); err != nil {
			return nil
		}
		if err := validNumber(v.PriceOpen, "PriceOpen"); err != nil {
			return nil
		}
		if err := validNumber(v.SL, "SL"); err != nil {
			return nil
		}
		if err := validNumber(v.TP, "TP"); err != nil {
			return nil
		}
		if err := validNumber(v.Swap, "Swap"); err != nil {
			return nil
		}
		if err := validNumber(v.PriceSL, "PriceSl"); err != nil {
			return nil
		}
		if err := validNumber(v.Profit, "Profit"); err != nil {
			return nil
		}
	}

	return nil
}

func validPage(bt string) error {
	// Using simple lexer is faster than regexp's
	// https://commandcenter.blogspot.com/2011/08/regular-expressions-in-lexing-and.html
	if len(bt) > 32 {
		return errors.New("'Page' is limited to 32 characters")
	}

	for _, b := range bt {
		if (b >= '0' && b <= '9') || (b >= 'a' && b <= 'z') || b == '_' || b == '-' {
			continue
		}
		return errors.New("'Page' may only contain lowercase latin letters, digits and following symbols '_-'")
	}

	return nil
}

func validString(bt string, fn string) error {
	if len(bt) > 32 {
		return errors.New("'" + fn + "' field is limited to 32 characters")
	}

	for _, b := range bt {
		if (b >= '0' && b <= '9') || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || b == '-' || b == ' ' || b == '(' || b == ')' || b == '.' || b == ',' {
			continue
		}
		return errors.New("'" + fn + "' field may only contain latin letters, digits and following symbols '_- ().,'")
	}
	return nil
}

func validTime(bt string, fn string) error {
	if len(bt) > 32 {
		return errors.New("'" + fn + "' field is limited to 32 characters")
	}

	for _, b := range bt {
		if (b >= '0' && b <= '9') || b == ' ' || b == '.' || b == ',' || b == ':' {
			continue
		}
		return errors.New("'" + fn + "' field may only contain digits and following symbols ' .,:'")
	}

	return nil
}

func validNumber(bt string, fn string) error {
	if len(bt) > 10 {
		return errors.New("'" + fn + "' field is limited to 10 characters")
	}

	for _, b := range bt {
		if (b >= '0' && b <= '9') || b == '.' || b == ',' || b == '-' {
			continue
		}
		return errors.New("'" + fn + "' field may only contain digits and following symbols '.,-'")
	}

	return nil
}

func (t Message) String() string {
	var ret string
	ret += fmt.Sprintf("{\"Page\":\"%s\",", string(t.Page))
	ret += fmt.Sprintf("\"ClientVersion\":\"%s\",", string(t.ClientVersion))
	ret += fmt.Sprintf("\"UpdateFreq\":\"%s\",", string(t.UpdateFreq))
	ret += fmt.Sprintf("\"Name\":\"%s\",", string(t.Name))
	ret += fmt.Sprintf("\"Login\":\"%s\",", string(t.Login))
	ret += fmt.Sprintf("\"Server\":\"%s\",", string(t.Server))
	ret += fmt.Sprintf("\"Company\":\"%s\",", string(t.Company))
	ret += fmt.Sprintf("\"Balance\":\"%s\",", string(t.Balance))
	ret += fmt.Sprintf("\"Equity\":\"%s\",", string(t.Equity))
	ret += fmt.Sprintf("\"Margin\":\"%s\",", string(t.Margin))
	ret += fmt.Sprintf("\"FreeMargin\":\"%s\",", string(t.FreeMargin))
	ret += fmt.Sprintf("\"MarginLevel\":\"%s\",", string(t.MarginLevel))
	ret += fmt.Sprintf("\"ProfitTotal\":\"%s\",", string(t.ProfitTotal))
	ret += fmt.Sprintf("\"Orders\":[")

	for _, v := range t.Orders {
		ret += fmt.Sprint(v)
	}

	ret += "]"

	return ret
}

// UpdateWith order a with data from order b
func (a *Order) UpdateWith(b Order) {
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

// String returns string representation
func (a Order) String() string {
	var ret string
	ret += fmt.Sprintf("{\"Symbol\":\"%s\",", string(a.Symbol))
	ret += fmt.Sprintf("\"TimeOpen\":\"%s\",", string(a.TimeOpen))
	ret += fmt.Sprintf("\"Type\":\"%s\",", string(a.Type))
	ret += fmt.Sprintf("\"InitVolume\":\"%s\",", string(a.InitVolume))
	ret += fmt.Sprintf("\"CurVolume\":\"%s\",", string(a.CurVolume))
	ret += fmt.Sprintf("\"PriceOpen\":\"%s\",", string(a.PriceOpen))
	ret += fmt.Sprintf("\"SL\":\"%s\",", string(a.SL))
	ret += fmt.Sprintf("\"TP\":\"%s\",", string(a.TP))
	ret += fmt.Sprintf("\"Swap\":\"%s\",", string(a.Swap))
	ret += fmt.Sprintf("\"PriceSL\":\"%s\",", string(a.PriceSL))
	ret += fmt.Sprintf("\"Profit\":\"%s\"}", string(a.Profit))
	return ret
}
