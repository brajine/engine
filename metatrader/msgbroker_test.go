package metatrader

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type brokerTestSuite struct {
	suite.Suite
	br          *BrokerFactory
	ws          *websocket.Conn
	testServer  *httptest.Server
	zapRecorder *observer.ObservedLogs
	zapObserver *zap.Logger
}

func TestBroker(t *testing.T) {
	suite.Run(t, new(brokerTestSuite))
}

func (b *brokerTestSuite) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	b.ws = ws
}

func (b *brokerTestSuite) SetupTest() {
	core, recorder := observer.New(zapcore.DebugLevel)
	b.zapRecorder = recorder
	b.zapObserver = zap.New(core)
	b.br = NewBroker(b.zapObserver.Sugar())
	b.testServer = httptest.NewServer(b)
}

func (b *brokerTestSuite) TearDownTest() {
	b.testServer.Close()

	println("*** Recorded logs ***")
	for _, log := range b.zapRecorder.All() {
		println(log.Message)
	}
	println("***")
}

func (b *brokerTestSuite) TestAddRemoveViewer() {
	println("TestAddViewer started")

	srvWs, _, err := b.makeClient()
	if b.NoError(err) {
		b.br.AddViewer(srvWs)
		b.NoError(b.waitForLogMessage("new viewer"))
		b.Equal(1, b.br.ViewersNumber())

		b.br.RemoveViewer(srvWs)
		b.NoError(b.waitForLogMessage("Viewer disconnected"))
		b.Equal(0, b.br.ViewersNumber())
	}
}

func (b *brokerTestSuite) TestSameViewer() {
	println("TestSameViewer started")

	srvWs, _, err := b.makeClient()
	if b.NoError(err) {
		b.NoError(b.br.AddViewer(srvWs))
		b.NoError(b.waitForLogMessage("new viewer"))
		b.Equal(1, b.br.ViewersNumber())

		b.Error(b.br.AddViewer(srvWs))
		b.Equal(1, b.br.ViewersNumber())
	}
}

func (b *brokerTestSuite) TestViewersNumber() {
	println("Test1KViewers started")

	cnt := 100
	srv := make([]*websocket.Conn, cnt)
	for i := 0; i < cnt; i++ {
		srvconn, _, err := b.makeClient()
		if b.NoError(err) {
			b.NoError(b.br.AddViewer(srvconn))
			srv[i] = srvconn
		}
	}
	// Wait Broker to complete
	b.NoError(b.waitForCountLogMessages("new viewer", cnt))
	b.Equal(cnt, b.br.ViewersNumber())

	// Remove all
	for i := 0; i < cnt; i++ {
		b.br.RemoveViewer(srv[i])
	}
	b.NoError(b.waitForCountLogMessages("Viewer disconnected", cnt))
	b.Equal(0, b.br.ViewersNumber())

	println("Intermedate messages for this test was switched off")
	b.zapRecorder.TakeAll()
}

func (b *brokerTestSuite) TestDirectMessage() {
	println("TestBroadcast started")

	cnt := 30
	srv := make([]*websocket.Conn, cnt)
	cli := make([]*websocket.Conn, cnt)
	for i := 0; i < cnt; i++ {
		srvconn, cliconn, err := b.makeClient()
		if b.NoError(err) {
			b.NoError(b.br.AddViewer(srvconn))
			srv[i] = srvconn
			cli[i] = cliconn
		}
	}
	// Wait Broker to complete
	b.NoError(b.waitForCountLogMessages("new viewer", cnt))
	b.Equal(cnt, b.br.ViewersNumber())

	readerNum := 22
	b.br.SendMessageToViewer(srv[readerNum], []byte("test"))
	b.NoError(b.waitForLogMessage("particular message"))

	for i := 0; i < cnt; i++ {
		cli[i].SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		_, data, err := cli[i].ReadMessage()
		if i == readerNum {
			if b.NoError(err) {
				b.Equal("test", string(data))
			}
		} else {
			b.Error(err)
		}
	}

	println("Intermedate messages for this test was switched off")
	b.zapRecorder.TakeAll()
}

func (b *brokerTestSuite) TestBroadcast() {
	println("TestBroadcast started")

	cnt := 100
	srv := make([]*websocket.Conn, cnt)
	cli := make([]*websocket.Conn, cnt)
	for i := 0; i < cnt; i++ {
		srvconn, cliconn, err := b.makeClient()
		if b.NoError(err) {
			b.NoError(b.br.AddViewer(srvconn))
			srv[i] = srvconn
			cli[i] = cliconn
		}
	}
	// Wait Broker to complete
	b.NoError(b.waitForCountLogMessages("new viewer", cnt))
	b.Equal(cnt, b.br.ViewersNumber())

	// Broadcast message
	b.br.SendMessage([]byte("test"))

	for i := 0; i < cnt; i++ {
		_, data, err := cli[i].ReadMessage()
		if b.NoError(err) {
			b.Equal("test", string(data))
		} else {
			return
		}
	}

	println("Intermedate messages for this test was switched off")
	b.zapRecorder.TakeAll()
}

func (b *brokerTestSuite) TestStop() {
	println("TestStop started")
	b.br.Stop()
	b.NoError(b.waitForLogMessage("Broker just closed all the Viewers"))
	b.NoError(b.waitForLogMessage("Broker is closed"))
}

func (b *brokerTestSuite) makeClient() (server, client *websocket.Conn, err error) {
	wsURL := strings.Replace(b.testServer.URL, "http://", "ws://", 1)
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if b.NoError(err) {
		return b.ws, ws, nil
	}
	return nil, nil, err
}

func (b *brokerTestSuite) waitForCountLogMessages(msg string, count int) error {
	st := time.Now()
	for {
		if b.countLogMessages(msg) == count {
			return nil
		}
		if time.Since(st) >= time.Second {
			return errors.New("timeout")
		}
		time.Sleep(time.Millisecond)
	}
}

func (b *brokerTestSuite) countLogMessages(msg string) int {
	cnt := 0
	for _, log := range b.zapRecorder.All() {
		if strings.Contains(log.Message, msg) {
			cnt++
		}
	}
	return cnt
}

func (b *brokerTestSuite) waitForLogMessage(msg string) error {
	st := time.Now()
	for {
		for _, log := range b.zapRecorder.All() {
			if strings.Contains(log.Message, msg) {
				return nil
			}
		}
		if time.Since(st) >= 100*time.Millisecond {
			return errors.New("timeout")
		}
		time.Sleep(time.Millisecond)
	}
}
