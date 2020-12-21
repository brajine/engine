package main

import (
	"encoding/gob"
	"engine/metatrader"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const TestTimeoutSeconds time.Duration = 3 * time.Second

type engineTestSuite struct {
	suite.Suite
	mt             *metatrader.Factory
	enc            *gob.Encoder
	dec            *gob.Decoder
	server, client net.Conn
	zapRecorder    *observer.ObservedLogs
	zapObserver    *zap.Logger
	testEcho       *echo.Echo
}

func TestServer(t *testing.T) {
	suite.Run(t, new(engineTestSuite))
}

func (e *engineTestSuite) SetupTest() {
	rand.Seed(time.Now().UnixNano())

	core, recorder := observer.New(zapcore.DebugLevel)
	e.zapRecorder = recorder
	e.zapObserver = zap.New(core)

	e.mt = metatrader.New("", e.zapObserver.Sugar())

	e.server, e.client = net.Pipe() // Emulate server connection
	enc := gob.NewEncoder(e.server)
	dec := gob.NewDecoder(e.server)
	go e.mt.ProcessMessages(enc, dec, "testconn")

	e.enc = gob.NewEncoder(e.client)
	e.dec = gob.NewDecoder(e.client)

	// API requests
	e.testEcho = echo.New()
}

func (e *engineTestSuite) TearDownTest() {
	e.enc = nil
	e.dec = nil
	e.client.Close()
	e.server.Close()
	e.testEcho.Close()

	ts := time.Now()
	for {
		for _, log := range e.zapRecorder.All() {
			if strings.Contains(log.Message, "Connection is closed") {
				e.printTestLogs()
				return
			}
		}
		if time.Since(ts) > (TestTimeoutSeconds) {
			e.printTestLogs()
			break
		}
		time.Sleep(time.Millisecond)
	}
}

func (e *engineTestSuite) TestRegister() {
	println("TestRegister started")
	msg := metatrader.Message{
		Page:          "test",
		ClientVersion: "0.1",
		UpdateFreq:    "second",
		Name:          "Name",
		Login:         "Login",
		Server:        "Server",
		Company:       "Company",
		Balance:       "0.1",
		Equity:        "0.1",
		Margin:        "0.1",
		FreeMargin:    "0.1",
		MarginLevel:   "1000",
		ProfitTotal:   "1000",
		Orders:        make(map[string]metatrader.Order),
	}
	ord := metatrader.Order{
		Symbol:     "EURGBP",
		TimeOpen:   "2020.12.25 10:08:23",
		Type:       "1",
		InitVolume: "0.1",
		CurVolume:  "0.1",
		PriceOpen:  "0.1",
		SL:         "0.1",
		TP:         "0.1",
		Swap:       "0",
		PriceSL:    "0",
		Profit:     "1",
	}
	msg.Orders["11111"] = ord

	resp, err := e.Push(&msg)
	if !e.NoError(err) {
		return
	}
	e.Empty(resp)

	acc := e.mt.PageExist("test")
	if e.NotNil(acc) {
		e.Equal("test", acc.Page)
		e.Equal("0.1", acc.ClientVersion)
		e.Equal("second", acc.UpdateFreq)
		e.Equal("Name", acc.Name)
		e.Equal("Login", acc.Login)
		e.Equal("Server", acc.Server)
		e.Equal("Company", acc.Company)
		e.Equal("0.1", acc.Balance)
		e.Equal("0.1", acc.Equity)
		e.Equal("0.1", acc.Margin)
		e.Equal("0.1", acc.FreeMargin)
		e.Equal("1000", acc.MarginLevel)
		e.Equal("1000", acc.ProfitTotal)
	}
	if e.Contains(acc.Orders, "11111") {
		ord := acc.Orders["11111"]
		e.Equal("EURGBP", ord.Symbol)
		e.Equal("2020.12.25 10:08:23", ord.TimeOpen)
		e.Equal("1", ord.Type)
		e.Equal("0.1", ord.InitVolume)
		e.Equal("0.1", ord.CurVolume)
		e.Equal("0.1", ord.PriceOpen)
		e.Equal("0.1", ord.SL)
		e.Equal("0.1", ord.TP)
		e.Equal("0", ord.Swap)
		e.Equal("0", ord.PriceSL)
		e.Equal("1", ord.Profit)
	}
}

func (e *engineTestSuite) TestPageExists() {
	println("TestPageExists started")

	// Register an account
	err := e.Send(&metatrader.Message{
		Page:       "test",
		UpdateFreq: "second",
	})
	e.Nil(err)

	resp, err := e.Recv()
	if e.NoError(err) {
		e.Empty(resp)
		e.NotNil(e.mt.PageExist("test"))
	}

	// Repeat from another connection
	done := make(chan bool)
	go func() {
		defer func() {
			done <- true
		}()
		server, client := net.Pipe()
		enc := gob.NewEncoder(server)
		dec := gob.NewDecoder(server)
		go e.mt.ProcessMessages(enc, dec, "testconn")

		enc2 := gob.NewEncoder(client)
		dec2 := gob.NewDecoder(client)

		err = enc2.Encode(metatrader.Message{
			Page:          "test",
			ClientVersion: "0.1",
			UpdateFreq:    "second",
		})
		e.Nil(err)

		// Wait for confirmation!
		resp := new(metatrader.ResponseMsg)
		err := dec2.Decode(&resp)
		if e.NoError(err) {
			e.NotEmpty(resp.Error)
		}
	}()
	<-done
}

func (e *engineTestSuite) TestAccountUpdate() {
	println("TestAccountUpdate started")

	// Register
	msg := metatrader.Message{
		Page:          "test",
		ClientVersion: "0.1",
		UpdateFreq:    "second",
	}

	resp, err := e.Push(&msg)
	if !e.NoError(err) {
		return
	}
	e.Empty(resp)

	acc := e.mt.PageExist("test")
	if e.NotNil(acc) {
		e.Equal("test", acc.Page)
		e.Equal("second", acc.UpdateFreq)
		e.Equal("0.1", acc.ClientVersion)
		e.Empty(acc.Name)
		e.Empty(acc.Login)
		e.Empty(acc.Server)
		e.Empty(acc.Company)
		e.Empty(acc.Balance)
		e.Empty(acc.Equity)
		e.Empty(acc.Margin)
		e.Empty(acc.FreeMargin)
		e.Empty(acc.MarginLevel)
		e.Empty(acc.ProfitTotal)
		e.Empty(acc.Orders)
	}

	// Update
	msg = metatrader.Message{
		Page:          "updated",
		ClientVersion: "0.2",
		UpdateFreq:    "minute",
		Name:          "Name",
		Login:         "Login",
		Server:        "Server",
		Company:       "Company",
		Balance:       "0.1",
		Equity:        "0.1",
		Margin:        "0.1",
		FreeMargin:    "0.1",
		MarginLevel:   "1000",
		ProfitTotal:   "1000",
	}

	resp, err = e.Push(&msg)
	if !e.NoError(err) {
		return
	}
	e.Empty(resp)

	e.Empty(e.mt.PageExist("updated"))
	acc = e.mt.PageExist("test")
	if e.NotNil(acc) {
		e.Equal("test", acc.Page, "Should not change after registration")
		e.Equal("0.1", acc.ClientVersion, "Should not change after registration")
		e.Equal("second", acc.UpdateFreq, "Should not change after registration")
		e.Equal("Name", acc.Name)
		e.Equal("Login", acc.Login)
		e.Equal("Server", acc.Server)
		e.Equal("Company", acc.Company)
		e.Equal("0.1", acc.Balance)
		e.Equal("0.1", acc.Equity)
		e.Equal("0.1", acc.Margin)
		e.Equal("0.1", acc.FreeMargin)
		e.Equal("1000", acc.MarginLevel)
		e.Equal("1000", acc.ProfitTotal)
	}
	e.Equal(0, len(acc.Orders))
}

func (e *engineTestSuite) TestOrdersUpdate() {
	println("TestOrdersUpdate started")

	// Register with 2 orders
	msg := metatrader.Message{
		Page:       "test",
		UpdateFreq: "second",
		Orders:     make(map[string]metatrader.Order),
	}
	ord1 := metatrader.Order{
		Symbol:     "EURGBP",
		TimeOpen:   "1111.11.11 11:11:11",
		Type:       "1",
		InitVolume: "0.1",
		CurVolume:  "0.1",
		PriceOpen:  "0.1",
		SL:         "0.1",
		TP:         "0.1",
		Swap:       "1",
		PriceSL:    "1",
		Profit:     "1",
	}
	ord2 := metatrader.Order{
		Symbol:     "GBPUSD",
		TimeOpen:   "2222.22.22 22:22:22",
		Type:       "2",
		InitVolume: "0.2",
		CurVolume:  "0.2",
		PriceOpen:  "0.2",
		SL:         "0.2",
		TP:         "0.2",
		Swap:       "2",
		PriceSL:    "2",
		Profit:     "2",
	}
	msg.Orders["11111"] = ord1
	msg.Orders["22222"] = ord2

	resp, err := e.Push(&msg)
	if !e.NoError(err) {
		return
	}
	e.Empty(resp)

	acc := e.mt.PageExist("test")
	if e.NotNil(acc) {
		e.Equal(2, len(acc.Orders))
		if e.Contains(acc.Orders, "11111") {
			ord := acc.Orders["11111"]
			e.Equal("EURGBP", ord.Symbol)
			e.Equal("1111.11.11 11:11:11", ord.TimeOpen)
			e.Equal("1", ord.Type)
			e.Equal("0.1", ord.InitVolume)
			e.Equal("0.1", ord.CurVolume)
			e.Equal("0.1", ord.PriceOpen)
			e.Equal("0.1", ord.SL)
			e.Equal("0.1", ord.TP)
			e.Equal("1", ord.Swap)
			e.Equal("1", ord.PriceSL)
			e.Equal("1", ord.Profit)
		}
		if e.Contains(acc.Orders, "22222") {
			ord := acc.Orders["22222"]
			e.Equal("GBPUSD", ord.Symbol)
			e.Equal("2222.22.22 22:22:22", ord.TimeOpen)
			e.Equal("2", ord.Type)
			e.Equal("0.2", ord.InitVolume)
			e.Equal("0.2", ord.CurVolume)
			e.Equal("0.2", ord.PriceOpen)
			e.Equal("0.2", ord.SL)
			e.Equal("0.2", ord.TP)
			e.Equal("2", ord.Swap)
			e.Equal("2", ord.PriceSL)
			e.Equal("2", ord.Profit)
		}
	}

	// Update with only 1 modified order (this will remove the first one)
	msg = metatrader.Message{
		Page:       "test",
		UpdateFreq: "second",
		Orders:     make(map[string]metatrader.Order),
	}
	ord2 = metatrader.Order{
		Symbol:     "GBPJPY",
		TimeOpen:   "3333.33.33 33:33:33",
		Type:       "3",
		InitVolume: "0.3",
		CurVolume:  "0.3",
		PriceOpen:  "0.3",
		SL:         "0.3",
		TP:         "0.3",
		Swap:       "3",
		PriceSL:    "3",
		Profit:     "3",
	}
	msg.Orders["22222"] = ord2

	resp, err = e.Push(&msg)
	if !e.NoError(err) {
		return
	}
	e.Empty(resp)

	acc = e.mt.PageExist("test")
	if e.NotNil(acc) {
		e.Equal(1, len(acc.Orders))
		if e.Contains(acc.Orders, "22222") {
			ord := acc.Orders["22222"]
			e.Equal("GBPJPY", ord.Symbol)
			e.Equal("3333.33.33 33:33:33", ord.TimeOpen)
			e.Equal("3", ord.Type)
			e.Equal("0.3", ord.InitVolume)
			e.Equal("0.3", ord.CurVolume)
			e.Equal("0.3", ord.PriceOpen)
			e.Equal("0.3", ord.SL)
			e.Equal("0.3", ord.TP)
			e.Equal("3", ord.Swap)
			e.Equal("3", ord.PriceSL)
			e.Equal("3", ord.Profit)
		}
	}

	// Update with new order (this will remove the previous)
	msg = metatrader.Message{
		Page:          "test",
		ClientVersion: "0.1",
		UpdateFreq:    "second",
		Orders:        make(map[string]metatrader.Order),
	}
	ord3 := metatrader.Order{
		Symbol: "GBPJPY",
	}
	msg.Orders["33333"] = ord3

	resp, err = e.Push(&msg)
	if !e.NoError(err) {
		return
	}
	e.Empty(resp)

	acc = e.mt.PageExist("test")
	if e.NotNil(acc) {
		e.Equal(1, len(acc.Orders))
		if e.Contains(acc.Orders, "33333") {
			ord := acc.Orders["33333"]
			e.Equal("GBPJPY", ord.Symbol)
			e.Empty(ord.TimeOpen)
			e.Empty(ord.Type)
			e.Empty(ord.InitVolume)
			e.Empty(ord.CurVolume)
			e.Empty(ord.PriceOpen)
			e.Empty(ord.SL)
			e.Empty(ord.TP)
			e.Empty(ord.Swap)
			e.Empty(ord.PriceSL)
			e.Empty(ord.Profit)
		}
	}

}

func (e *engineTestSuite) TestValidateAccountFailed() {
	println("TestValidateAccountFailed started")

	var fails []*metatrader.Message
	fails = append(fails, &metatrader.Message{})
	fails = append(fails, &metatrader.Message{Page: "'Page' may only contain lowercase latin letters, digits and following symbols '_-'"})
	fails = append(fails, &metatrader.Message{Page: "Page_may_only_contain_lowercase_latin_letters_digits_and_following_symbols"})
	fails = append(fails, &metatrader.Message{Page: "Page may only contain"})
	fails = append(fails, &metatrader.Message{Page: "%page&"})
	fails = append(fails, &metatrader.Message{Page: "Page"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "Simple update frequency"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "Hour"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", ClientVersion: "Version is a string so it is limited to 32 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", ClientVersion: "%%"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Name: "Simple name that is longer than 32 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Name: "$$"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Login: "Simple login that is longer than 32 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Login: "$$"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Server: "Simple server that is longer than 32 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Server: "$$"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Company: "Simple company that is longer than 32 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Company: "$$"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Balance: "Field is limited to 10 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Balance: "$$"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Equity: "Field is limited to 10 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Equity: "$$"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Margin: "Field is limited to 10 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", Margin: "$$"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", FreeMargin: "Field is limited to 10 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", FreeMargin: "$$"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", MarginLevel: "Field is limited to 10 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", MarginLevel: "$$"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", ProfitTotal: "Field is limited to 10 characters"})
	fails = append(fails, &metatrader.Message{Page: "page", UpdateFreq: "second", ProfitTotal: "$$"})

	for _, v := range fails {
		resp, err := e.PushToNewInstance(v)
		if e.Nil(err) {
			e.NotEmpty(resp.Error)
		}
	}

	println("Intermedate messages for this test was switched off")
	e.zapRecorder.TakeAll()
}

func (e *engineTestSuite) TestMaxFreeOrders() {
	println("TestMaxFreeOrders started")
	msg := metatrader.Message{
		Page:       "test",
		UpdateFreq: "second",
		Orders:     make(map[string]metatrader.Order),
	}

	for i := 0; i < metatrader.MaxFreeOrders; i++ {
		msg.Orders[strconv.Itoa(i)] = metatrader.Order{}
	}

	resp, err := e.Push(&msg)
	if e.NoError(err) {
		e.Empty(resp.Error)
	}

	// Add extra one order
	msg.Orders[strconv.Itoa(metatrader.MaxFreeOrders)] = metatrader.Order{}

	resp, err = e.Push(&msg)
	if e.NoError(err) {
		e.NotEmpty(resp.Error)
	}
}

func (e *engineTestSuite) TestStatsAPIHandler() {
	println("TestStatsAPIHandler started")

	// Push 2 accounts
	resp, err := e.Push(&metatrader.Message{
		Page:       "test",
		UpdateFreq: "second",
	})
	if e.NoError(err) {
		e.Empty(resp.Error)
	}

	resp, err = e.PushToNewInstance(&metatrader.Message{
		Page:       "test2",
		UpdateFreq: "second",
	})
	if e.NoError(err) {
		e.Empty(resp.Error)
	}

	// Read stats
	code, body, err := e.GetStats()
	if e.Nil(err) {
		e.Equal(200, code)
		e.Contains(body, "\"online\":2")
		e.Contains(body, "\"page\":\"test\"")
		e.Contains(body, "\"page\":\"test2\"")
	}
}

func (e *engineTestSuite) TestRestAPIHandler() {
	println("TestRestAPIHandler started")

	// Push account
	msg := metatrader.Message{
		Page:          "test",
		ClientVersion: "0.1",
		UpdateFreq:    "second",
		Name:          "Name",
		Login:         "Login",
		Server:        "Server",
		Company:       "Company",
		Balance:       "0.1",
		Equity:        "0.1",
		Margin:        "0.1",
		FreeMargin:    "0.1",
		MarginLevel:   "1000",
		ProfitTotal:   "1000",
		Orders:        make(map[string]metatrader.Order),
	}
	ord := metatrader.Order{
		Symbol:     "EURGBP",
		TimeOpen:   "2020.12.25 10:08:23",
		Type:       "1",
		InitVolume: "0.1",
		CurVolume:  "0.1",
		PriceOpen:  "0.1",
		SL:         "0.1",
		TP:         "0.1",
		Swap:       "0",
		PriceSL:    "0",
		Profit:     "1",
	}
	msg.Orders["11111"] = ord

	resp, err := e.Push(&msg)
	if !e.NoError(err) {
		return
	}
	e.Empty(resp)

	// Read stats
	code, body, err := e.GetRest("test")
	if e.Nil(err) {
		e.Equal(200, code)
		e.Contains(body, "\"updatefreq\":\"second\"")
		e.Contains(body, "\"name\":\"Name\"")
		e.Contains(body, "\"login\":\"Login\"")
		e.Contains(body, "\"server\":\"Server\"")
		e.Contains(body, "\"company\":\"Company\"")
		e.Contains(body, "\"balance\":\"0.1\"")
		e.Contains(body, "\"equity\":\"0.1\"")
		e.Contains(body, "\"margin\":\"0.1\"")
		e.Contains(body, "\"freemargin\":\"0.1\"")
		e.Contains(body, "\"marginlevel\":\"1000\"")
		e.Contains(body, "\"profittotal\":\"1000\"")

		e.Contains(body, "\"orderscount\":1")
		e.Contains(body, "\"11111\":{")
		e.Contains(body, "\"symbol\":\"EURGBP\"")
		e.Contains(body, "\"timeopen\":\"2020.12.25 10:08:23\"")
		e.Contains(body, "\"type\":\"1\"")
		e.Contains(body, "\"initvolume\":\"0.1\"")
		e.Contains(body, "\"curvolume\":\"0.1\"")
		e.Contains(body, "\"priceopen\":\"0.1\"")
		e.Contains(body, "\"sl\":\"0.1\"")
		e.Contains(body, "\"tp\":\"0.1\"")
		e.Contains(body, "\"swap\":\"0\"")
		e.Contains(body, "\"pricesl\":\"0\"")
		e.Contains(body, "\"profit\":\"1\"")
	}
}

func (e *engineTestSuite) TestWebSocket() {
	println("TestWebSocket started")

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(e.wsHandler))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	u, err := url.Parse(s.URL)
	e.Nil(err)
	u.Scheme = "ws"

	// Connect to the server - fails NO ACCOUNT
	_, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if e.Error(err) {
		e.Equal(404, resp.StatusCode)
	} else {
		return
	}

	// Create account
	_, err = e.Push(&metatrader.Message{
		Page:       "test",
		UpdateFreq: "second",
	})
	if !e.NoError(err) {
		return
	}

	// Connect to the server - SUCCESS
	ws, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if e.NoError(err) {
		e.Equal(101, resp.StatusCode) // Switching protocols
		defer ws.Close()

		_, p, err := ws.ReadMessage()
		if e.NoError(err) {
			str := string(p)
			println("Read WS message from server:", str)
			e.Contains(str, "\"page\":\"test\"")
			e.Contains(str, "\"updatefreq\":\"second\"")
			e.NotContains(str, "balance")
		}
	} else {
		return
	}

	// Update account
	mtresp, err := e.Push(&metatrader.Message{
		Balance: "100",
	})
	e.Empty(mtresp)
	if e.NoError(err) {
		_, p, err := ws.ReadMessage()
		if e.NoError(err) {
			str := string(p)
			println("Read WS update from server:", str)
			e.Contains(str, "\"balance\":\"100\"")
		}
	}
}

func (e *engineTestSuite) wsHandler(w http.ResponseWriter, r *http.Request) {
	c := e.testEcho.NewContext(r, w)
	c.SetPath("/api/rest/test")
	c.SetParamNames("page")
	c.SetParamValues("test")
	e.mt.WssAPIHandler(c)
}

// TODO:
// - Account.ToJSON() !!!
// - Validate account success
// - Validate orders
// - Performance check
// - Messages frequency
// - Read timeout
// - Max message size
// - Broadcast for low-speed or stalled clients

func (e *engineTestSuite) GetStats() (code int, body string, err error) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.testEcho.NewContext(req, rec)
	c.SetPath("/api/stats")
	err = e.mt.StatsAPIHandler(c)
	return rec.Code, rec.Body.String(), err
}

func (e *engineTestSuite) GetRest(page string) (code int, body string, err error) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.testEcho.NewContext(req, rec)
	c.SetPath("/api/rest/:page")
	c.SetParamNames("page")
	c.SetParamValues(page)

	err = e.mt.RestAPIHandler(c)
	return rec.Code, rec.Body.String(), err
}

func (e *engineTestSuite) PushToNewInstance(msg *metatrader.Message) (*metatrader.ResponseMsg, error) {
	errChan := make(chan error)
	respChan := make(chan *metatrader.ResponseMsg)
	go func() {
		server, client := net.Pipe()
		enc := gob.NewEncoder(server)
		dec := gob.NewDecoder(server)
		go e.mt.ProcessMessages(enc, dec, "testconn")

		enc2 := gob.NewEncoder(client)
		dec2 := gob.NewDecoder(client)

		if err := enc2.Encode(msg); err != nil {
			errChan <- err
			return
		}

		// Wait for confirmation!
		resp := new(metatrader.ResponseMsg)
		if err := dec2.Decode(&resp); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
		respChan <- resp
	}()

	err := <-errChan
	if err != nil {
		return nil, err
	}
	return <-respChan, err
}

func (e *engineTestSuite) Push(msg *metatrader.Message) (*metatrader.ResponseMsg, error) {
	if err := e.Send(msg); err != nil {
		return nil, err
	}
	return e.Recv()
}

func (e *engineTestSuite) Send(msg *metatrader.Message) error {
	errChan := make(chan error)
	go func() {
		err := e.enc.Encode(*msg)
		if err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	timer := time.NewTimer(TestTimeoutSeconds)
	select {
	case <-timer.C:
		return errors.New("Write failed: timeout")
	case err := <-errChan:
		if err != nil {
			return err
		}
		return nil
	}
}

func (e *engineTestSuite) Recv() (*metatrader.ResponseMsg, error) {
	errChan := make(chan error)
	respChan := make(chan *metatrader.ResponseMsg)
	go func() {
		msg := new(metatrader.ResponseMsg)
		err := e.dec.Decode(msg)
		if err != nil {
			errChan <- err
			return
		}
		respChan <- msg
	}()

	timer := time.NewTimer(TestTimeoutSeconds)
	select {
	case <-timer.C:
		return nil, errors.New("Read failed: timeout")
	case err := <-errChan:
		return nil, err
	case msg := <-respChan:
		return msg, nil
	}
}

func (e *engineTestSuite) printTestLogs() {
	println("*** Recorded logs ***")
	for _, log := range e.zapRecorder.All() {
		println(log.Message)
	}
	println("***")
	println()
}
func (e *engineTestSuite) waitNewSecond() {
	// Wait for the new second
	sec := time.Now().Unix()
	for {
		if time.Now().Unix() != sec {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
func (e *engineTestSuite) rand32Page() string {
	var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyz_-")
	b := make([]rune, 32)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
func (e *engineTestSuite) rand32String() string {
	var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-()., ")
	b := make([]rune, 32)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
func (e *engineTestSuite) randFreq() string {
	var freqs = []string{"second", "minute"}
	return freqs[rand.Intn(2)]
}
func (e *engineTestSuite) randFloat() string {
	return fmt.Sprintf("%.2f", rand.Float64()*1000.0)
}
func (e *engineTestSuite) randVersion() string {
	s := strconv.Itoa(rand.Intn(10))
	return s + ".0"
}
func (e *engineTestSuite) randInt() string {
	return strconv.Itoa(rand.Int())
}
func (e *engineTestSuite) randType() string {
	return strconv.Itoa(rand.Intn(8))
}
func (e *engineTestSuite) randTime() string {
	v := time.Duration(rand.Intn(1440)) * time.Second
	return time.Now().Add(-1 * v).Format("2006-01-02 15:04:05")
}
