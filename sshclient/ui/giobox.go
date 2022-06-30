package ui

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

type GioBox struct {
	window *app.Window
	theme  *material.Theme
	title  *material.LabelStyle
	input  *GioInput
	output *GioOutput
}

func NewGioBox() (*GioBox, error) {
	w := app.NewWindow()
	t := material.NewTheme(gofont.Collection())

	title := material.H1(t, "Gio")
	maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	title.Color = maroon
	title.Alignment = text.Middle

	input := NewGioInput()
	output := NewGioOutput()
	input.SingleLine = true

	return &GioBox{
		window: w,
		theme:  t,
		title:  &title,
		input:  input,
		output: output,
	}, nil
}

func (self *GioBox) runLoop() error {
	var ops op.Ops

	for {
		e := <-self.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			self.title.Layout(gtx)

			layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return self.output.Layout(self.theme, gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return self.input.Layout(self.theme, gtx)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func (self *GioBox) Run() {
	go func() {
		err := self.runLoop()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
