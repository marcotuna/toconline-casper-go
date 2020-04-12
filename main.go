// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"toconline-casper-go/model"
	wsClient "toconline-casper-go/pkg/ws/client"
	"toconline-casper-go/utils"

	"github.com/joho/godotenv"
)

const (
	serverAddr = "app3.toconline.pt"
)

const regexExpression = `^([0-9]+):([a-zA-Z]+):([0-9]+|\{.*?\}+)?:?(\{.*\}+)`

func main() {
	log.SetFlags(0)
	godotenv.Load()

	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)

	// Sign in
	reqClient, err := utils.HTTPClientReq(
		"https://app3.toconline.pt/login/sign-in",
		nil,
		[]*utils.HTTPClientHeader{
			{
				Key:   "Authorization",
				Value: fmt.Sprintf("Basic %s", utils.SetBasicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))),
			},
		},
		nil,
	)

	if err != nil {
		log.Fatal("Login:", err)
	}

	userSession := &model.Session{
		AccessToken: "",
	}

	if reqClient.StatusCode >= 200 && reqClient.StatusCode < 300 {
		userSession = model.SessionFromJSON(bytes.NewReader(reqClient.Body))
	} else {
		log.Fatal("Could not authenticate with the provided credentials")
	}

	if userSession == nil || len(userSession.AccessToken) == 0 {
		log.Fatal("Session is not valid")
	}

	u := url.URL{Scheme: "wss", Host: serverAddr, Path: "/epaper"}
	log.Printf("Connecting to %s", u.String())

	httpHeader := http.Header{}
	httpHeader.Add("Sec-WebSocket-Protocol", "casper-epaper")

	wsClientHandler, err := wsClient.New(u.String(), httpHeader)

	if err != nil {
		log.Fatal("Websocket:", err)
	}

	// Captures Received Websocket Messages
	go wsClientHandler.OnMessage()

	setSession := model.NewMessage(
		model.VerbSet,
		model.MessageTarget{
			Target: model.TargetSession,
		},
		userSession,
	)

	wsClientHandler.Send(setSession.ToCasper())
	_ = model.WaitReceive(wsClientHandler.Payload, setSession.GetInvokeID())

	wsClientHandler.Send([]byte(`1:GET:{"target":"http","url":"https://cdb-master.toconline.pt/cdb/entity/list","headers":{"content-type":"application/json","accept":"application/json"}}`))
	_ = model.WaitReceive(wsClientHandler.Payload, setSession.GetInvokeID())

	// for {
	// 	select {
	// 	case msg := <-wsClientHandler.Payload:
	// 		msgReponse := model.MessageFromCasper(*msg)
	// 		fmt.Println(msgReponse)
	// 	case <-interrupt:
	// 		log.Println("interrupt")
	// 		err := wsClientHandler.Client.WriteMessage(
	// 			websocket.CloseMessage,
	// 			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	// 		)
	// 		if err != nil {
	// 			log.Println("write close:", err)
	// 			return
	// 		}

	// 		return
	// 	}
	// }
}
