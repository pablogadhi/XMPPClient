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
)

const host = "alumchat.xyz"
const port = "5222"

type decoderIO struct {
	r io.Reader
	w io.Writer
}

var connected = false

// Listen starts listening to connection
func Listen() {
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Println(err)
		return
	}
	for _, ip := range ips {
		fmt.Println(ip.String())
	}

	conn, err := net.Dial("tcp", ips[0].String()+":"+port)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	// First XMPP call
	_, err = fmt.Fprintf(conn, stanza.Start, host)

	if err != nil {
		log.Println(err)
		return
	}

	// Set XML Decoder as connection listener
	decoder := xml.NewDecoder(decoderIO{conn, os.Stderr})
	if decoder != nil {
		fmt.Println("Succes!")
	}

	for {
		t, err := decoder.Token()
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			// After feature response (second response)
			case "features":
				if !connected {
					_, err = fmt.Fprintf(conn, stanza.StartTLS)

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

				tlsConn := tls.Client(conn, tlsConf)
				err := tlsConn.Handshake()
				if err != nil {
					log.Println(err.Error())
					return
				}

				//Set decoder with new connection and restart stream
				decoder = xml.NewDecoder(decoderIO{tlsConn, os.Stderr})
				_, err = fmt.Fprintf(tlsConn, stanza.Stream, host)

				if err != nil {
					log.Println(err)
					return
				}

				connected = true

				fmt.Fprintf(tlsConn, stanza.GetRegistrationFields)
				fmt.Fprintf(tlsConn, stanza.Register, "yair", "test1")

				// Authentication
				// user := "test1"
				// password := "test1"
				// raw := "\x00" + user + "\x00" + password
				// enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
				// base64.StdEncoding.Encode(enc, []byte(raw))
				// fmt.Fprintf(tlsConn, stanza.Auth, "urn:ietf:params:xml:ns:xmpp-sasl", enc)

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
