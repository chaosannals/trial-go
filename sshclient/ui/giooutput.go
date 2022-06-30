package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type GioOutput struct {
	widget.Editor
	old string
}

func NewGioOutput() *GioOutput {
	return &GioOutput{}
}

func (ed *GioOutput) Changed() bool {
	newText := ed.Editor.Text()
	changed := newText != ed.old
	ed.old = newText
	return changed
}

func (ed *GioOutput) SetText(s string) {
	ed.old = s
	ed.Editor.SetText(s)
}

func (ed *GioOutput) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	// Determine colors based on the state of the editor.
	borderWidth := float32(0.5)
	borderColor := color.NRGBA{A: 107}
	switch {
	case ed.Editor.Focused():
		borderColor = th.Palette.ContrastBg
		borderWidth = 2
		// case ed.Invalid:
		// 	borderColor = color.NRGBA{R: 200, A: 0xFF}
	}

	// draw an editor with a border.
	return widget.Border{
		Color:        borderColor,
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(borderWidth),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(4)).Layout(gtx,
			material.Editor(th, &ed.Editor, "").Layout)
	})
}
