package app

import (
	"errors"
	"fmt"
	"time"
	"toconline-casper-go/model"
	"toconline-casper-go/pkg/casper-ws/client"
)

// App ...
type App struct {
	UserSession *model.Authentication
	WsClient    *client.Websocket
	Request     *model.CasperRequest
}

// NewApp creates new app
func NewApp(userSession *model.Authentication, wsClient *client.Websocket) *App {

	if wsClient != nil {
		// Captures Received Websocket Messages
		go wsClient.OnMessage()
	}

	return &App{
		UserSession: userSession,
		WsClient:    wsClient,
		Request:     model.NewCasperRequest(),
	}
}

// WaitReceive ...
func WaitReceive(payloadMessage chan *[]byte, messageID *int) (*model.CasperResponse, error) {
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-payloadMessage:
			casperMessage := model.MessageFromCasper(*msg)

			if casperMessage == nil {
				return nil, errors.New("Received message is not valid")
			}

			if *casperMessage.InvokeID == *messageID {
				return casperMessage, nil
			}
		case <-ticker.C:
			return nil, errors.New("Message not received. Timeout")
		}
	}

}

// CasperRequest sends payload to ws server and waits for response
// This is done by sending an ID on the request and waiting for it on the response
func (a *App) CasperRequest(reqSend *model.CasperRequest) (*model.CasperResponse, error) {
	// Send request
	reqMessage := reqSend.ToCasper()
	fmt.Println(string(reqMessage))
	a.WsClient.Send(reqMessage)

	// Wait for the response
	return WaitReceive(a.WsClient.Payload, reqSend.InvokeID)
}

// CasperRequestRaw sends payload to ws server and waits for response
// This is done by sending an ID on the request and waiting for it on the response
func (a *App) CasperRequestRaw(reqMessage []byte) (*model.CasperResponse, error) {
	fmt.Println(string(reqMessage))
	// Send request
	a.WsClient.Send(reqMessage)

	// Wait for the response
	return WaitReceive(a.WsClient.Payload, model.NewInt(0))
}
