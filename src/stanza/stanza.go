package stanza

const Start = "<?xml version='1.0'?>\n" +
	"<stream:stream to='%s' xmlns='jabber:client'\n" +
	" xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"

const StartTLS = "<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>"

const Stream = "<stream:stream to='%s' xmlns:stream='http://etherx.jabber.org/streams' xmlns='jabber:client' xml:lang='en' version='1.0'>"

const GetRegistrationFields = "<iq id='reg1' type='get'><query xmlns='jabber:iq:register' /></iq>"

const Register = "<iq id='reg2' type='set'><query xmlns='jabber:iq:register'><username>%s</username><password>%s</password></query></iq>"

const Auth = "<auth xmlns='%s' mechanism='PLAIN'>%s</auth>\n"

const Unregister = "<iq type='set' from='%s' id='unreg1'>\n" +
	"<query xmlns='jabber:iq:register'>\n" +
	"<remove/>\n" +
	"</query>\n" +
	"</iq>"

const SessionBind = "<iq id='c109ad79-c372-4afd-b1bd-ac5df9d07a93-3' type='set'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind' /></iq>"

const SimplePresence = "<presence/>"

const GetRoster = "<iq id='roster1' type='get'>\n" +
	"<query xmlns='jabber:iq:roster'/>\n" +
	"</iq>"

const AddUser = "<presence to='%s' type='subscribe'/>"

const AcceptUser = "<presence to='%s' type='subscribed'/>"

const SendMessage = "<message to='%s' type='%s' id='%s' xml:lang='en'><body>%s</body></message>"
