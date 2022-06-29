package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type FyneBox struct {
	application fyne.App
	window      fyne.Window
	command     binding.String
}

func NewFyneBox() (*FyneBox, error) {
	a := app.New()
	w := a.NewWindow("sshclient")
	w.Resize(fyne.Size{Width: 800, Height: 600})

	inputText := binding.NewString()
	outputText := binding.NewString()
	output := widget.NewEntryWithData(outputText)
	output.MultiLine = true
	input := widget.NewEntryWithData(inputText)
	input.OnChanged = func(text string) {
		if strings.Index(text, "\n") < 0 {
			return
		}
		if ot, err := outputText.Get(); err == nil {
			outputText.Set(ot + "\n")
		}
	}
	w.SetContent(container.NewVBox(
		widget.NewLabel("ssh client"),
		output,
		input,
	))

	return &FyneBox{
		application: a,
		window:      w,
		command:     inputText,
	}, nil
}

func (self *FyneBox) Run() {
	self.window.ShowAndRun()
}
