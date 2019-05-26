/*
Copyright (C) 2019 by Martin Langlotz aka stackshadow

This file is part of gopilot, an rewrite of the copilot-project in go

gopilot is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, version 3 of this License

gopilot is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with gopilot.  If not, see <http://www.gnu.org/licenses/>.
*/

package gwebsocket

import (
	"flag"
	"fmt"
	"gopilot/clog"
	"gopilot/gbus"
	"gopilot/nodeName"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Gwebsocket struct {
	logging        clog.Logger
	conn           *websocket.Conn
	bus            gbus.Socketbus
	remoteNodeName string
}

var startWebSocket bool
var webSocketAddr string
var startWebServer bool
var webServerRoot string
var webServerAddr string

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func ParseCmdLine() {
	flag.StringVar(&webSocketAddr, "websocket.addr", "127.0.0.1:3333", "Web-Socket Adress for webinterface")

	flag.BoolVar(&startWebServer, "webserver", false, "Enable Webserver")
	flag.StringVar(&webServerAddr, "webserver.addr", "127.0.0.1:9090", "Web-Server Adress for webinterface")
	flag.StringVar(&webServerRoot, "webserver.root", "/app/www/gopilot", "Root directory of webfiles")

}

// Init the current websocket
func (curCWs *Gwebsocket) Init() {
	curCWs.logging = clog.New("WS")
	curCWs.bus.Init()
}

// Serve [GOROUTINE] Will start subroutines for websocket/webserver
func (curCWs *Gwebsocket) Serve() {

	curCWs.bus.Subscribe("ws-all", "", "", curCWs.OnMessage)

	go curCWs.serveWebsocket()

	curCWs.remoteNodeName, _ = curCWs.bus.Connect(
		gbus.SocketFileName,
		gbus.Msg{NodeSource: "", GroupSource: ""},
	)

	if startWebServer == true {
		go curCWs.serveWebserver()
	}
}

func (curCWs *Gwebsocket) serveWebsocket() {
	curCWs.logging.Info(fmt.Sprintf("Start websocker-server on %s", webSocketAddr))
	http.HandleFunc("/echo-protocol", curCWs.onWebsocketMessage)
	http.ListenAndServe(webSocketAddr, nil)
}

func (curCWs *Gwebsocket) serveWebserver() {
	curCWs.logging.Info(fmt.Sprintf("Start webserver on %s", webServerAddr))
	fs := http.FileServer(http.Dir(webServerRoot))
	http.Handle("/", fs)
	http.ListenAndServe(webServerAddr, nil)
}

func (curCWs *Gwebsocket) onWebsocketMessage(w http.ResponseWriter, r *http.Request) {

	var err error
	curCWs.conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer curCWs.conn.Close()

	for {
		messageType, message, err := curCWs.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if messageType == websocket.BinaryMessage {
			log.Printf("We dont handle binary-messages yet")
			continue
		}

		if messageType == websocket.TextMessage {
			curCWs.logging.Debug(string(message[:]))
		}

		curMessage, err := gbus.FromJSONString(string(message))
		if err != nil {
			log.Println("error: ", err)
			continue
		}
		curMessage.NodeSource = mynodename.NodeName
		curMessage.GroupSource = "ws"
		curMessage.NodeTarget = curCWs.remoteNodeName

		curMessage.ContextSet("ws")
		curCWs.bus.PublishMsg(curMessage)

	}
}

func (curCWs *Gwebsocket) OnMessage(message *gbus.Msg, group, command, payload string) {
	if curCWs.conn == nil {
		return
	}

	// is the message from us ?
	if message.ContextGet() == "ws" {
		return
	}

	curCWs.logging.Info(fmt.Sprintf("%s/%s", group, command))

	jsonString, err := message.ToJSONString()
	if err != nil {
		fmt.Println("error:", err)
	}

	curCWs.logging.Debug(jsonString)
	err = curCWs.conn.WriteMessage(websocket.TextMessage, []byte(jsonString))
	if err != nil {
		log.Println("write:", err)
	}
}
