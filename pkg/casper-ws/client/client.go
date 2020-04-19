package client

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Websocket ...
type Websocket struct {
	Client  *websocket.Conn
	Payload chan *[]byte
}

// SetDefaultHeaders ...
func SetDefaultHeaders() *http.Header {
	httpHeader := http.Header{}
	httpHeader.Add("Sec-WebSocket-Protocol", "casper-epaper")

	return &httpHeader
}

// New instanciates connection with the websocket server
func New(url string, httpHeader http.Header) (*Websocket, error) {

	// Set default headers when no headers were passed
	if httpHeader == nil {
		httpHeader = *SetDefaultHeaders()
	}

	// Try to establish a websocket connection
	c, _, err := websocket.DefaultDialer.Dial(url, httpHeader)
	if err != nil {
		return nil, err
	}

	return &Websocket{
		Client:  c,
		Payload: make(chan *[]byte),
	}, nil
}

// Close stops connection with the websocket server
func (ws *Websocket) Close() error {
	return ws.Client.Close()
}

// Send emits a message to the websocket server
func (ws *Websocket) Send(message []byte) error {
	err := ws.Client.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	return nil
}

// Receive waits for incoming message and returns its payload
func (ws *Websocket) Receive() []byte {
	for {
		select {
		case msg := <-ws.Payload:
			return *msg
		}
	}
}

// OnMessage handles message when receives it from the websocket server
func (ws *Websocket) OnMessage() {
	for {
		_, message, err := ws.Client.ReadMessage()
		if err != nil {
			return
		}

		ws.Payload <- &message
	}

}
