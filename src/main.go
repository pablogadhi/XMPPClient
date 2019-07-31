package main

import (
	"log"
	"net"
	"os"
	"streamer"
	"ui"

	"github.com/rivo/tview"
)

const host = "alumchat.xyz"
const port = "5222"

var debugMode = false

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "--debug" {
		debugMode = true
	}

	// DNS Lookout
	ips, err := net.LookupIP(host)
	if err != nil {
		if debugMode {
			log.Println(err)
		}
		return
	}

	// Connect to server
	conn, err := net.Dial("tcp", ips[0].String()+":"+port)
	if err != nil {
		if debugMode {
			log.Println(err)
		}
		return
	}

	defer conn.Close()

	app := tview.NewApplication()

	go streamer.Listen(&conn, host, debugMode, app)

	if !debugMode {
		ui.SetUI(&conn, app, host)
	} else {
		for {

		}
	}
}
