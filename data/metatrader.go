package data

import (
	"errors"
	"fmt"
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

// TradesMsg keeps all Metatrader data for each particular client
//easyjson:json
type TradesMsg struct {
	Page          string    `json:"-"`
	ClientVersion string    `json:"-"`
	Updated       time.Time `json:"-"`
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
	// Use Ticket as a map key
	Orders map[string]OrderType `json:"orders"`
}

// OrderType represent one Metatrader order
// Tickets are always sent, other values sent only if not null AND changed since last update
type OrderType struct {
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

// Validate incoming TradesMsg
func (t *TradesMsg) Validate() error {
	if err := validPage(t.Page); err != nil {
		return err
	}
	if err := validString(t.ClientVersion, "ClientVersion"); err != nil {
		return err
	}
	if err := validString(t.UpdateFreq, "UpdateFreq"); err != nil {
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
	for k, v := range t.Orders {
		if err := validNumber(k, "Ticket"); err != nil {
			return nil
		}
		if err := validNumber(v.Symbol, "Symbol"); err != nil {
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
		if (b < '0' || b > '9') && (b < 'a' || b > 'z') && (b != '_' && b != '-') {
			return errors.New("'Page' may only contain lowercase latin letters, digits and following symbols '_-'")
		}
	}

	return nil
}

func validString(bt string, fn string) error {
	if len(bt) > 32 {
		return errors.New("'" + fn + "' field is limited to 32 characters")
	}

	for _, b := range bt {
		if (b < '0' || b > '9') && (b < 'a' || b > 'z') && (b < 'A' && b > 'Z') && (b != '_' && b != '-' && b != ' ' && b != '(' && b != ')' && b != '.' && b != ',') {
			return errors.New("'" + fn + "' field may only contain latin letters, digits and following symbols '_- ().,'")
		}
	}

	return nil
}

func validTime(bt string, fn string) error {
	if len(bt) > 32 {
		return errors.New("'" + fn + "' field is limited to 32 characters")
	}

	for _, b := range bt {
		if (b < '0' || b > '9') && (b != ' ' && b != '.' && b != ',' && b != ':') {
			return errors.New("'" + fn + "' field may only contain digits and following symbols ' .,:'")
		}
	}

	return nil
}

func validNumber(bt string, fn string) error {
	if len(bt) > 10 {
		return errors.New("'" + fn + "' field is limited to 10 characters")
	}

	for _, b := range bt {
		if (b < '0' || b > '9') && (b != '.' && b != ',' && b != '-') {
			return errors.New("'" + fn + "' field may only contain digits and following symbols '.,-'")
		}
	}

	return nil
}

func (t TradesMsg) String() string {
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

func (t OrderType) String() string {
	var ret string
	ret += fmt.Sprintf("{\"Symbol\":\"%s\",", string(t.Symbol))
	ret += fmt.Sprintf("\"TimeOpen\":\"%s\",", string(t.TimeOpen))
	ret += fmt.Sprintf("\"Type\":\"%s\",", string(t.Type))
	ret += fmt.Sprintf("\"InitVolume\":\"%s\",", string(t.InitVolume))
	ret += fmt.Sprintf("\"CurVolume\":\"%s\",", string(t.CurVolume))
	ret += fmt.Sprintf("\"PriceOpen\":\"%s\",", string(t.PriceOpen))
	ret += fmt.Sprintf("\"SL\":\"%s\",", string(t.SL))
	ret += fmt.Sprintf("\"TP\":\"%s\",", string(t.TP))
	ret += fmt.Sprintf("\"Swap\":\"%s\",", string(t.Swap))
	ret += fmt.Sprintf("\"PriceSL\":\"%s\",", string(t.PriceSL))
	ret += fmt.Sprintf("\"Profit\":\"%s\"}", string(t.Profit))
	return ret
}
