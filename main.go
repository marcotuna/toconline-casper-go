// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/url"
	"os"
	"toconline-casper-go/app"
	wsClient "toconline-casper-go/pkg/casper-ws/client"

	"github.com/joho/godotenv"
)

const (
	serverAddr = "app3.toconline.pt"
)

func main() {
	log.SetFlags(0)
	godotenv.Load()

	// Initialize Websocket
	u := url.URL{Scheme: "wss", Host: serverAddr, Path: "/epaper"}
	log.Printf("Connecting to %s", u.String())

	wsClientHandler, err := wsClient.New(u.String(), nil)

	if err != nil {
		log.Fatal("Websocket: ", err)
	}

	// Initialize App
	app := app.NewApp(nil, wsClientHandler)

	// Authenticate
	authReq, err := app.Authenticate(
		os.Getenv("USERNAME"),
		os.Getenv("PASSWORD"),
	)

	if err != nil {
		log.Fatal(err)
	}

	// Set User Session on Casper
	_, err = app.SetSession(authReq)

	if err != nil {
		log.Fatal(err)
	}

	// Get all Entities
	entityList, err := app.GetEntityList()

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range entityList {
		// fmt.Printf("%#v\n", v)
		_, err := app.EntitySwitch(v.ID)

		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("%#v\n", msgStatus)
	}

	app.GetEntityReport()

}
