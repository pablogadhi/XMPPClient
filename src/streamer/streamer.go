package streamer

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"stanza"
	"ui"
	"xmpp"

	"github.com/rivo/tview"
)

type decoderIO struct {
	r io.Reader
	w io.Writer
}

type message struct {
	XMLName xml.Name `xml:"message"`
	From    string   `xml:"from,attr"`
	Body    string   `xml:"body"`
}

var connected = false
var connectionId = ""

// Listen starts listening to connection
func Listen(conn *net.Conn, host string, debugMode bool, app *tview.Application) {
	// First XMPP call
	_, err := fmt.Fprintf(*conn, stanza.Start, host)

	if err != nil {
		log.Println(err)
		return
	}

	// Set XML Decoder as connection listener
	var decoder *xml.Decoder
	if debugMode {
		decoder = xml.NewDecoder(decoderIO{*conn, os.Stderr})
	} else {
		decoder = xml.NewDecoder(*conn)
	}

	for {
		t, err := decoder.Token()
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			// After feature response (second response)
			case "features":
				if !connected {
					_, err = fmt.Fprintf(*conn, stanza.StartTLS)

					if err != nil {
						log.Println(err)
						return
					}
				}
			// Proceed to do TLS handshake
			case "proceed":
				tlsConf := &tls.Config{
					ServerName:         host,
					InsecureSkipVerify: true,
				}

				tlsConn := tls.Client(*conn, tlsConf)
				err := tlsConn.Handshake()
				if err != nil {
					log.Println(err.Error())
					return
				}

				//Set decoder with new connection and restart stream
				if debugMode {
					decoder = xml.NewDecoder(decoderIO{tlsConn, os.Stderr})
				} else {
					decoder = xml.NewDecoder(tlsConn)
				}
				_, err = fmt.Fprintf(tlsConn, stanza.Stream, host)

				if err != nil {
					log.Println(err)
					return
				}

				*conn = tlsConn
				connected = true

				// xmpp.Authenticate(conn, ui.Username, ui.Password)
				xmpp.Authenticate(conn, "yair", "test1")
			case "success":
				ui.GoToMainView()
				_, err = fmt.Fprintf(*conn, stanza.Stream, host)
				_, err = fmt.Fprintf(*conn, stanza.SessionBind)
				xmpp.GetRoster(conn)
				xmpp.SendSimplePresence(conn)

			case "stream":
				if connected {
					fmt.Println(t.Name.Local)
					connectionId = t.Attr[3].Value
				}
			case "presence":
				// Accept Subscriptions
				if len(t.Attr) == 3 && t.Attr[2].Value == "subscribe" {
					xmpp.AcceptUser(conn, t.Attr[0].Value)
				}
			case "item":
				if len(t.Attr) >= 1 && t.Name.Space == "jabber:iq:roster" {
					ui.AddContact(app, t.Attr[0].Value)
				}
			case "message":
				var msg message
				err := decoder.DecodeElement(&msg, &t)
				if err != nil {
					log.Println(err)
				}
				ui.AddMessage(app, msg.From+": "+msg.Body)
				xmpp.SendMessage(conn, "test1@alumchat.xyz", "Ora ora")
				// fmt.Printf(stanza.SendMessage, "test1@alumchat.xyz", "yair@alumchat.xyz/"+connectionId, "Ora ora")
			}
			// fmt.Println(t.Name.Space + " " + t.Name.Local)
		}

		if err != nil || t == nil {
			// fmt.Println("TagError!")
		}
	}
}

func (dio decoderIO) Read(p []byte) (n int, err error) {
	n, err = dio.r.Read(p)
	if n > 0 {
		dio.w.Write(p[0:n])
		dio.w.Write([]byte("\n"))
	}
	return
}
