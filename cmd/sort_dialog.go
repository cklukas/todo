package cmd

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SortModal is a simple modal with a dropdown for sort mode
type SortModal struct {
	*tview.Form
	DialogHeight int
	frame        *tview.Frame
	optionIndex  int
	options      []string
	done         func(string, bool)
}

func (m *SortModal) GetFrame() *tview.Frame {
	return m.frame
}

func NewSortModal(title, lane string, current string) *SortModal {
	form := tview.NewForm()
	m := &SortModal{Form: form, DialogHeight: 7, frame: tview.NewFrame(form), optionIndex: 0,
		options: []string{"", SortColor, SortDue, SortCreated, SortModified, SortPriority}, done: nil}

	form.SetCancelFunc(func() {
		if m.done != nil {
			m.done("", false)
		}
	})

	labels := []string{"manual", "color", "due", "created", "modified", "priority"}
	idx := 0
	for i, v := range m.options {
		if v == current {
			idx = i
			break
		}
	}
	m.optionIndex = idx
	form.AddDropDown("Sort:", labels, idx, func(option string, index int) {
		m.optionIndex = index
	})
	m.frame.AddText(fmt.Sprintf("Sort tasks in lane '%s'", lane), false, 0, tcell.ColorDarkGray)
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

	m.frame.SetTitle(fmt.Sprintf(" %v ", title))
	m.frame.SetBorders(0, 0, 1, 0, 0, 0).
		SetBorder(true).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(1, 1, 1, 1)

	return m
}

func (m *SortModal) SetDoneFunc(handler func(string, bool)) {
	m.done = handler
}

// Draw draws this modal with a surrounding frame.
func (m *SortModal) Draw(screen tcell.Screen) {
	// Determine width similar to ModalInput.
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

	// Draw the frame.
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}

// Focus delegates focus to the embedded form.
func (m *SortModal) Focus(delegate func(p tview.Primitive)) {
	delegate(m.Form)
}

// HasFocus returns whether the form has focus.
func (m *SortModal) HasFocus() bool {
	return m.Form.HasFocus()
}

// SetFocus passes the focus index to the embedded form.
func (m *SortModal) SetFocus(index int) *SortModal {
	m.Form.SetFocus(index)
	return m
}

// MouseHandler forwards mouse events to the form and captures clicks inside the dialog.
func (m *SortModal) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		consumed, capture = m.Form.MouseHandler()(action, event, setFocus)
		if !consumed && action == tview.MouseLeftDown && m.InRect(event.Position()) {
			setFocus(m)
			consumed = true
		}
		return
	})
}

// InputHandler returns the handler for this primitive.
func (m *SortModal) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if m.frame.HasFocus() {
			if handler := m.frame.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
		}
	})
}
