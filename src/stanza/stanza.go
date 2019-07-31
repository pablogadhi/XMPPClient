package stanza

const Start = "<?xml version='1.0'?>\n" +
	"<stream:stream to='%s' xmlns='jabber:client'\n" +
	" xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"

const StartTLS = "<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>"

const Stream = "<stream:stream to='%s' xmlns:stream='http://etherx.jabber.org/streams' xmlns='jabber:client' xml:lang='en' version='1.0'>"

const GetRegistrationFields = "<iq id='reg1' type='get'><query xmlns='jabber:iq:register' /></iq>"

const Register = "<iq id='reg2' type='set'><query xmlns='jabber:iq:register'><username>%s</username><password>%s</password></query></iq>"

const Auth = "<auth xmlns='%s' mechanism='PLAIN'>%s</auth>\n"
