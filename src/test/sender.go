package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-xmpp"
)

var server = "alumchat.xyz"
var username = "test1@alumchat.xyz"
var password = "test1"
var status = flag.String("status", "xa", "status")
var statusMessage = flag.String("status-msg", "I for one welcome our new codebot overlords.", "status message")
var notls = flag.Bool("notls", true, "No TLS")
var debug = flag.Bool("debug", true, "debug output")
var session = flag.Bool("session", false, "use server session")

func serverName(host string) string {
	return strings.Split(host, ":")[0]
}

func main() {
	if !*notls {
		xmpp.DefaultConfig = tls.Config{
			ServerName:         server,
			InsecureSkipVerify: false,
		}
	} else {
		xmpp.DefaultConfig = tls.Config{
			ServerName:         server,
			InsecureSkipVerify: true,
		}
	}

	var talk *xmpp.Client
	var err error
	options := xmpp.Options{Host: server,
		User:                         username,
		Password:                     password,
		NoTLS:                        *notls,
		Debug:                        *debug,
		Session:                      *session,
		Status:                       *status,
		StatusMessage:                *statusMessage,
		InsecureAllowUnencryptedAuth: true,
		StartTLS:                     true,
	}

	// talk, err = xmpp.NewClient(server, username, password, *debug)
	talk, err = options.NewClient()

	if err != nil {
		fmt.Println(err.Error())
	}

	go func() {
		for {
			chat, err := talk.Recv()
			if err != nil {
				log.Fatal(err)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				fmt.Println(v.Remote, v.Text)
			case xmpp.Presence:
				fmt.Println(v.From, v.Show)
			}
		}
	}()
	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.TrimRight(line, "\n")

		tokens := strings.SplitN(line, " ", 2)
		if len(tokens) == 2 {
			talk.Send(xmpp.Chat{Remote: tokens[0], Type: "chat", Text: tokens[1]})
		}
	}
}
