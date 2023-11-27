package gui

import (
	"io"

	"github.com/brylemp/chatroom/pkg/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GUI struct {
	app        *tview.Application
	mainView   *tview.Grid
	chatBox    *tview.TextView
	inputField *tview.InputField
}

func New() *GUI {
	app := tview.NewApplication()

	chatBox := tview.NewTextView()
	chatBox.SetBorder(true)
	chatBox.SetChangedFunc(func() {
		app.Draw()
		chatBox.ScrollToEnd()
	})

	placeHolderStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).Dim(true)

	inputField := tview.NewInputField().
		SetPlaceholder("Enter text here...").
		SetPlaceholderStyle(placeHolderStyle).
		SetFieldBackgroundColor(tcell.ColorBlack)

	mainView := tview.NewGrid().
		AddItem(chatBox, 0, 0, 20, 1, 1, 0, false).
		AddItem(inputField, 20, 0, 1, 1, 1, 1, true)

	app.SetRoot(mainView, true).
		EnableMouse(true)

	return &GUI{
		app:        app,
		mainView:   mainView,
		chatBox:    chatBox,
		inputField: inputField,
	}
}

func (gui *GUI) Run(title string) error {
	gui.chatBox.SetTitle(title)

	return gui.app.Run()
}

func (gui *GUI) SetInputReader(rw io.ReadWriter) {
	gui.inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() != tcell.KeyEnter {
			return event
		}
		err := util.SendMessage(rw, gui.inputField.GetText())
		if err == nil {
			gui.inputField.SetText("")
			return event
		}

		err = util.SendMessage(gui.chatBox, "Disconnected from the server")
		if err != nil {
			return nil
		}
		gui.inputField.SetText("")
		gui.inputField.SetDisabled(true)
		gui.inputField.SetPlaceholder("")
		gui.inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { return nil })

		return event
	})
}

func (gui *GUI) GetOutputWriter() io.Writer {
	return gui.chatBox
}
