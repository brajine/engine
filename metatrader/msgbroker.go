package metatrader

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type viewUpdater struct {
	ws         *websocket.Conn
	dataChan   chan []byte
	closeChan  chan bool
	signalChan chan *websocket.Conn
	log        *zap.SugaredLogger
}

func (v *viewUpdater) run() {
	go func() {
		defer func() {
			v.ws.Close()
			v.log.Info("Viewer disconnected ", v.ws.RemoteAddr())
		}()
		for {
			select {
			case <-v.closeChan:
				return
			case data := <-v.dataChan:
				v.log.Debug("Viewer sent a message to websocket ", v.ws.RemoteAddr())
				v.ws.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
				err := v.ws.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					// Let the Broker know we're finished
					v.signalChan <- v.ws
					return
				}
			}
		}
	}()
}

// Send data to Updater
func (v *viewUpdater) send(data []byte) {
	v.dataChan <- data
}

func (v *viewUpdater) close() {
	v.closeChan <- true
}

// BrokerFactory manage broadcasting of messages
type BrokerFactory struct {
	dataChan   chan []byte
	customChan chan *customMessage
	closeChan  chan bool
	addChan    chan *websocket.Conn
	removeChan chan *websocket.Conn
	updaters   map[*websocket.Conn]*viewUpdater
	signalChan chan *websocket.Conn
	log        *zap.SugaredLogger
}

// NewBroker ...
func NewBroker(log *zap.SugaredLogger) *BrokerFactory {
	br := BrokerFactory{
		updaters:   make(map[*websocket.Conn]*viewUpdater),
		dataChan:   make(chan []byte, 5),
		customChan: make(chan *customMessage, 5),
		closeChan:  make(chan bool, 1),
		signalChan: make(chan *websocket.Conn, 5),
		addChan:    make(chan *websocket.Conn, 5),
		removeChan: make(chan *websocket.Conn, 5),
		log:        log,
	}

	br.run()
	return &br
}

// run broker Manager
func (b *BrokerFactory) run() {
	go func() {
		defer func() {
			b.log.Debug("Broker is closed")
		}()
		for {
			select {
			case <-b.closeChan: // Close the Broker and all viewers
				for _, upd := range b.updaters {
					upd.close()
				}
				b.log.Debug("Broker just closed all the Viewers")
				return
			case c := <-b.customChan: // Send a message to one particular viewer
				b.updaters[c.ws].send(c.data)
				b.log.Debug("Broker sent particular message to viewer ", c.ws.RemoteAddr())
			case data := <-b.dataChan: // Broadcast message to all viewers
				for _, upd := range b.updaters {
					upd.send(data)
				}
				b.log.Debug("Broker broadcasted a message")
			case ws := <-b.addChan: // Add new Viewer
				b.updaters[ws] = &viewUpdater{
					ws:         ws,
					log:        b.log,
					dataChan:   make(chan []byte, 5),
					closeChan:  make(chan bool),
					signalChan: b.signalChan,
				}
				b.updaters[ws].run()
				b.log.Debug("Broker added new viewer to pool ", ws.RemoteAddr)
			case ws := <-b.removeChan: // Remove Viewer
				if _, ok := b.updaters[ws]; ok {
					b.updaters[ws].close()
					delete(b.updaters, ws)
					b.log.Debug("Broker removed viewer from pool ", ws.RemoteAddr)
				}
			case closedUpdater := <-b.signalChan: // Viewer got a Send error and sould be removed
				delete(b.updaters, closedUpdater)
				b.log.Debug("Broker got Closed signal from viewer, and removed it from pool ", closedUpdater)
			}
		}
	}()
}

// AddViewer to viewers pool, also create processing goroutine
func (b *BrokerFactory) AddViewer(viewer *websocket.Conn) error {
	if _, ok := b.updaters[viewer]; !ok {
		b.addChan <- viewer
		return nil
	}

	return errors.New("Failed to add a Viewer: already exists")
}

// RemoveViewer from viewers pool
func (b *BrokerFactory) RemoveViewer(viewer *websocket.Conn) {
	b.removeChan <- viewer
}

// SendMessage to Broker Manager (and further for all connected Viewers)
func (b *BrokerFactory) SendMessage(data []byte) {
	b.dataChan <- data
}

// ViewersNumber for testing purposes
func (b *BrokerFactory) ViewersNumber() int {
	return len(b.updaters)
}

// Send message to one particular Viewer
type customMessage struct {
	ws   *websocket.Conn
	data []byte
}

// SendMessageToViewer sends a direct message to one particular viewer
func (b *BrokerFactory) SendMessageToViewer(ws *websocket.Conn, data []byte) {
	b.customChan <- &customMessage{
		ws:   ws,
		data: data,
	}
}

// Stop the broker
func (b *BrokerFactory) Stop() {
	b.closeChan <- true
}
