package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type MainBox struct {
	application fyne.App
	window fyne.Window
	command binding.String
}

func NewMainBox() (*MainBox, error) {
	a := app.New()
	w := a.NewWindow("sshclient")

	c := binding.NewString()
	w.SetContent(container.NewVBox(
		widget.NewLabel("ssh client"),
		widget.NewEntryWithData(c),
	))
	
	return &MainBox{
		application: a,
		window: w,
		command: c,
	}, nil
}

func (self *MainBox) Run() {
	self.window.ShowAndRun()
}