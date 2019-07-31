package ui

import (
	"net"
	"xmpp"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var Username = ""
var Password = ""
var userToAdd = ""
var messageToSend = ""
var selectedUser = ""

var contactList = tview.NewList()
var chatRooms = tview.NewList()
var messageList = tview.NewList()
var pages = tview.NewPages()

// SetUI starts the UI
func SetUI(conn *net.Conn, ui *tview.Application, host string) {
	loginGrid := tview.NewGrid().
		SetRows(0, 0, 0).
		SetColumns(0, 0, 0)
	form := tview.NewForm().
		AddInputField("First name", "", 30, nil, func(text string) {
			Username = text
		}).
		AddPasswordField("Password", "", 30, '*', func(text string) {
			Password = text
		}).
		AddButton("Register", func() {
			xmpp.Register(conn, Username, Password)
		}).
		AddButton("Login", func() {
			xmpp.Authenticate(conn, Username, Password)
		})
	form.SetBorder(true).SetTitle("XMPPClient").SetTitleAlign(tview.AlignCenter)
	form.SetButtonsAlign(tview.AlignCenter)
	loginGrid.AddItem(form, 1, 1, 1, 1, 0, 0, true)

	pages.AddPage("login", loginGrid, true, true)

	mainGrid := tview.NewGrid().
		SetColumns(0, 0, 0, 0).
		SetRows(0, 0, 0, 0)

	newUserForm := tview.NewForm().
		AddInputField("Username", "", 20, nil, func(text string) {
			userToAdd = text
		}).
		AddButton("Add", func() {
			xmpp.AddUser(conn, userToAdd)
		})
	newUserForm.SetBorder(true).SetTitle("Add Contact")

	newMessageForm := tview.NewForm().
		AddInputField("Message", "", 80, nil, func(text string) {
			messageToSend = text
		}).
		AddButton("Send", func() {
			xmpp.SendMessage(conn, selectedUser, messageToSend)
		})
	newMessageForm.SetBorder(true).SetTitle("Write Message").SetTitleAlign(tview.AlignLeft)
	newMessageForm.SetButtonsAlign(tview.AlignRight)

	contactList.SetBorder(true).SetTitle("Contacts").SetTitleAlign(tview.AlignLeft)
	chatRooms.SetBorder(true).SetTitle("ChatRooms").SetTitleAlign(tview.AlignLeft)
	messageList.SetBorder(true).SetTitle("Messages")

	mainGrid.AddItem(newUserForm, 0, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(contactList, 1, 0, 2, 1, 0, 0, true)
	mainGrid.AddItem(chatRooms, 3, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(messageList, 0, 1, 3, 3, 0, 0, false)
	mainGrid.AddItem(newMessageForm, 3, 1, 1, 2, 0, 0, true)

	pages.AddPage("main", mainGrid, true, false)

	ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlQ {
			ui.Stop()
		} else if event.Key() == tcell.KeyF1 {
			ui.SetFocus(newUserForm)
		} else if event.Key() == tcell.KeyF2 {
			ui.SetFocus(contactList)
		} else if event.Key() == tcell.KeyF3 {
			ui.SetFocus(chatRooms)
		} else if event.Key() == tcell.KeyF4 {
			ui.SetFocus(newMessageForm)
		} else if event.Key() == tcell.KeyF5 {

		}
		return event
	})

	if err := ui.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

// GoToMainView changes the view after login
func GoToMainView() {
	pages.SwitchToPage("main")
}

func AddContact(app *tview.Application, user string) {
	rune := int32(contactList.GetItemCount())
	contactList.AddItem(user, "", rune, func() {
		selectedUser = user
	})
	app.Draw()
}

func AddMessage(app *tview.Application, message string) {
	rune := int32(messageList.GetItemCount())
	messageList.AddItem(message, "", rune, nil)
	app.Draw()
}
