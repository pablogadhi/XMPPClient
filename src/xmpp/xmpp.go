package xmpp

import (
	"encoding/base64"
	"fmt"
	"net"
	"stanza"
)

// Register user
func Register(conn *net.Conn, username string, password string) {
	fmt.Fprintf(*conn, stanza.GetRegistrationFields)
	fmt.Fprintf(*conn, stanza.Register, username, password)
}

// Authenticate user
func Authenticate(conn *net.Conn, username string, password string) {
	raw := "\x00" + username + "\x00" + password
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	base64.StdEncoding.Encode(enc, []byte(raw))
	fmt.Fprintf(*conn, stanza.Auth, "urn:ietf:params:xml:ns:xmpp-sasl", enc)
}

// DeleteUser from server
func DeleteUser(conn *net.Conn, username string, host string) {
	mail := username + "@" + host
	fmt.Fprintf(*conn, stanza.Unregister, mail)
}

//SendSimplePresence sends starting presence
func SendSimplePresence(conn *net.Conn) {
	fmt.Fprintf(*conn, stanza.SimplePresence)
}

// GetRoster for user
func GetRoster(conn *net.Conn) {
	fmt.Fprintf(*conn, stanza.GetRoster)
}

// AddUser adds user to roster
func AddUser(conn *net.Conn, to string) {
	fmt.Fprintf(*conn, stanza.AddUser, to)
}

// AcceptUser sends presence confirmation
func AcceptUser(conn *net.Conn, to string) {
	fmt.Fprintf(*conn, stanza.AcceptUser, to)
}

// SendMessage does what it says
func SendMessage(conn *net.Conn, to string, message string) {
	fmt.Fprintf(*conn, stanza.SendMessage, to, "chat", "slakdfjn324jsakldfj324", message)
}
