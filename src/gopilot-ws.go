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

package main

import (
	"flag"
	"gopilot/clog"
	"gopilot/nodeName"
	"gwebsocket"
	"time"
)

func main() {

	var websocket gwebsocket.Gwebsocket

	// ########################## Command line parse ##########################
	clog.ParseCmdLine()
	mynodename.ParseCmdLine()
	gwebsocket.ParseCmdLine()
	flag.Parse()

	// ########################## Init ##########################
	clog.Init()
	mynodename.Init()

	// the webclient
	websocket.Init()
	websocket.Serve()

	for {
		time.Sleep(time.Minute * 1)
	}
}
