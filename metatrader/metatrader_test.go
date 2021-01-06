package metatrader

import (
	"encoding/gob"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Run this test as:
// go test -v -count=1 ./metatrader/

const testMetatraderPort string = ":8182"

// Use globals for simplicity
var enc *gob.Encoder
var dec *gob.Decoder

func TestMetatrader(t *testing.T) {
	t.Skip()
	t.Log("Server started on " + testMetatraderPort)

	var conn net.Conn
	var err error
	conn, enc, dec, err = startTestServer()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if t.Run("Connection", connection) {
		t.Run("Page", page)
		t.Run("ClientVersion", clientVersion)
		t.Run("UpdateFrequency", updateFreq)
		t.Run("Name", name)
		t.Run("Login", login)
		t.Run("Server", server)
		t.Run("Company", company)
		t.Run("Balance", balance)
		t.Run("Equity", equity)
		t.Run("Margin", margin)
		t.Run("FreeMargin", freemargin)
		t.Run("MarginLevel", marginlevel)
		t.Run("ProfitTotal", profittotal)

		t.Run("Ticket", ticket)
		t.Run("Symbol", symbol)
		t.Run("TimeOpen", timeopen)
		t.Run("OrderType", ordtype)
		t.Run("OrderInitVolume", initvolume)
		t.Run("OrderCurVolume", curvolume)
		t.Run("OrderPriceOpen", priceopen)
		t.Run("OrderSL", stoploss)
		t.Run("OrderTP", takeprofit)
		t.Run("OrderSwap", swap)
		t.Run("OrderPriceSL", pricesl)
		t.Run("OrderProfit", ordprofit)

		t.Run("HeavyMessage", heavyMessage)
	}
}

func connection(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg)
	}
}

func page(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Equal(t, "page", msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func clientVersion(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Equal(t, "1.0", msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func updateFreq(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Equal(t, "second", msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func name(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Equal(t, "name", msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func login(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Equal(t, "login", msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func server(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Equal(t, "server", msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func company(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Equal(t, "company", msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func balance(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Equal(t, "balance", msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func equity(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Equal(t, "equity", msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func margin(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Equal(t, "margin", msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func freemargin(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Equal(t, "freemargin", msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func marginlevel(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Equal(t, "marginlevel", msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func profittotal(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Equal(t, "profittotal", msg.ProfitTotal)
		assert.Empty(t, msg.OrdersCount)
		assert.Empty(t, msg.Orders)
		assert.Zero(t, msg.OrdersCount)
	}
}

func ticket(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"])
			}
		}
	}
}

func symbol(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Equal(t, "market", msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func timeopen(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Equal(t, "timeopen", msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func ordtype(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Equal(t, "type", msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func initvolume(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Equal(t, "initvolume", msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func curvolume(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Equal(t, "curvolume", msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func priceopen(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Equal(t, "priceopen", msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func stoploss(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Equal(t, "stoploss", msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func takeprofit(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Equal(t, "takeprofit", msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func swap(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Equal(t, "swap", msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func pricesl(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Equal(t, "pricesl", msg.Orders["111"].PriceSL)
				assert.Empty(t, msg.Orders["111"].Profit)
			}
		}
	}
}

func ordprofit(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Empty(t, msg.Page)
		assert.Empty(t, msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Empty(t, msg.UpdateFreq)
		assert.Empty(t, msg.Name)
		assert.Empty(t, msg.Login)
		assert.Empty(t, msg.Server)
		assert.Empty(t, msg.Company)
		assert.Empty(t, msg.Balance)
		assert.Empty(t, msg.Equity)
		assert.Empty(t, msg.Margin)
		assert.Empty(t, msg.FreeMargin)
		assert.Empty(t, msg.MarginLevel)
		assert.Empty(t, msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		if assert.Equal(t, 1, len(msg.Orders)) {
			if assert.Contains(t, msg.Orders, "111") {
				assert.Empty(t, msg.Orders["111"].Symbol)
				assert.Empty(t, msg.Orders["111"].TimeOpen)
				assert.Empty(t, msg.Orders["111"].Type)
				assert.Empty(t, msg.Orders["111"].InitVolume)
				assert.Empty(t, msg.Orders["111"].CurVolume)
				assert.Empty(t, msg.Orders["111"].PriceOpen)
				assert.Empty(t, msg.Orders["111"].SL)
				assert.Empty(t, msg.Orders["111"].TP)
				assert.Empty(t, msg.Orders["111"].Swap)
				assert.Empty(t, msg.Orders["111"].PriceSL)
				assert.Equal(t, "profit", msg.Orders["111"].Profit)
			}
		}
	}
}

func heavyMessage(t *testing.T) {
	msg := new(Message)
	err := dec.Decode(msg)
	if assert.Nil(t, err) {
		assert.Equal(t, "page", msg.Page)
		assert.Equal(t, "clientversion", msg.ClientVersion)
		assert.Empty(t, msg.Started)
		assert.Empty(t, msg.Updated)
		assert.Equal(t, "second", msg.UpdateFreq)
		assert.Equal(t, "name", msg.Name)
		assert.Equal(t, "login", msg.Login)
		assert.Equal(t, "server", msg.Server)
		assert.Equal(t, "company", msg.Company)
		assert.Equal(t, "balance", msg.Balance)
		assert.Equal(t, "equity", msg.Equity)
		assert.Equal(t, "margin", msg.Margin)
		assert.Equal(t, "freemargin", msg.FreeMargin)
		assert.Equal(t, "marginlevel", msg.MarginLevel)
		assert.Equal(t, "profittotal", msg.ProfitTotal)
		assert.Zero(t, msg.OrdersCount)

		cnt := 10
		if assert.Equal(t, cnt, len(msg.Orders)) {
			for i := 0; i < cnt; i++ {
				n := strconv.Itoa(i)
				tick := n + n + n
				if assert.Contains(t, msg.Orders, tick) {
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].Symbol)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].TimeOpen)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].Type)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].InitVolume)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].CurVolume)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].PriceOpen)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].SL)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].TP)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].Swap)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].PriceSL)
					assert.Equal(t, tick, msg.Orders[OrderTicket(tick)].Profit)
				}
			}
		}
	}
}

func startTestServer() (net.Conn, *gob.Encoder, *gob.Decoder, error) {
	l, err := net.Listen("tcp", testMetatraderPort)
	if err != nil {
		return nil, nil, nil, err
	}

	conn, err := l.Accept()
	if err != nil {
		return nil, nil, nil, err
	}

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	return conn, enc, dec, nil
}
