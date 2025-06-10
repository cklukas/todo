package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ColorModal is a simple modal with a dropdown for lane color
type ColorModal struct {
	*tview.Form
	DialogHeight int
	frame        *tview.Frame
	optionIndex  int
	options      []string
	done         func(string, bool)
}

func (m *ColorModal) GetFrame() *tview.Frame {
	return m.frame
}

func NewColorModal(title, current string) *ColorModal {
	form := tview.NewForm()
	m := &ColorModal{Form: form, DialogHeight: 7, frame: tview.NewFrame(form), optionIndex: 0,
		options: []string{"", "blue", "green", "red", "yellow", "white", "darkcyan", "black", "darkmagenta"}, done: nil}

	form.SetCancelFunc(func() {
		if m.done != nil {
			m.done("", false)
		}
	})

	labels := []string{"default", "blue", "green", "red", "yellow", "white", "darkcyan", "black", "darkmagenta"}
	idx := 0
	for i, v := range m.options {
		if v == current {
			idx = i
			break
		}
	}
	m.optionIndex = idx
	form.AddDropDown("Color:", labels, idx, func(option string, index int) {
		m.optionIndex = index
	})
	m.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(0, 0, 0, 0)

	m.AddButton("OK", func() {
		if m.done != nil {
			m.done(m.options[m.optionIndex], true)
		}
	})
	m.AddButton("Cancel", func() {
		if m.done != nil {
			m.done("", false)
		}
	})

	m.frame.SetTitle(" " + title + " ")
	m.frame.SetBorders(0, 0, 1, 0, 0, 0).
		SetBorder(true).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(1, 1, 1, 1)

	return m
}

func (m *ColorModal) SetDoneFunc(handler func(string, bool)) {
	m.done = handler
}

// Draw draws this modal with a surrounding frame.
func (m *ColorModal) Draw(screen tcell.Screen) {
	buttonsWidth := 30
	screenWidth, screenHeight := screen.Size()
	width := screenWidth / 3
	if width < buttonsWidth {
		width = buttonsWidth
	}
	height := m.DialogHeight
	width += 4
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}

func (m *ColorModal) Focus(delegate func(p tview.Primitive)) {
	delegate(m.Form)
}

func (m *ColorModal) HasFocus() bool {
	return m.Form.HasFocus()
}

func (m *ColorModal) SetFocus(index int) *ColorModal {
	m.Form.SetFocus(index)
	return m
}

func (m *ColorModal) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		consumed, capture := m.Form.MouseHandler()(action, event, setFocus)
		if !consumed && action == tview.MouseLeftDown && m.InRect(event.Position()) {
			setFocus(m)
			consumed = true
		}
		return consumed, capture
	})
}

func (m *ColorModal) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if m.frame.HasFocus() {
			if handler := m.frame.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
		}
	})
}
