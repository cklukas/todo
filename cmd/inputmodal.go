package cmd

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ModalInput is based on Modal from tview, but has an input field instead
type ModalInput struct {
	*tview.Form
	frame     *tview.Frame
	main      string
	secondary string
	done      func(string, string, bool)
}

func NewModalInput(title string) *ModalInput {
	form := tview.NewForm()
	_, taskInp := form.AddInputField("Task:", "", 50, nil, nil)
	_, secondaryInp := form.AddInputField("Created / due:", "", 50, nil, nil)

	m := &ModalInput{form, tview.NewFrame(form), "", "", nil}

	taskInp.SetChangedFunc(func(text string) {
		m.main = text
	})
	secondaryInp.SetChangedFunc(func(text string) {
		m.secondary = text
	})

	m.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(0, 0, 0, 0)

	m.AddButton("OK", func() {
		if m.done != nil {
			m.done(m.main, m.secondary, true) // Passed
		}
	})
	m.AddButton("Cancel", func() {
		if m.done != nil {
			m.done(m.main, m.secondary, false)
		}
	})
	m.frame.SetTitle(fmt.Sprintf(" %v ", title))
	m.frame.SetBorders(0, 0, 1, 0, 0, 0).
		SetBorder(true).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(1, 1, 1, 1)

	return m
}

// SetValue sets the current value in the item
func (m *ModalInput) SetValue(text string, secondary string) {
	m.main = text
	m.secondary = secondary
	m.Clear(false)
	m.AddInputField("", text, 50, nil, func(text string) {
		if len(text) == 0 {
			text = "(empty)"
		}
		m.main = text
	})
	m.AddInputField("", secondary, 50, nil, func(text string) {
		m.secondary = text
	})
}

// SetDoneFunc sets the done func for this input.
// Will be called with the text of the input and a boolean for OK or cancel button.
func (m *ModalInput) SetDoneFunc(handler func(string, string, bool)) *ModalInput {
	m.done = handler
	return m
}

// Draw draws this primitive onto the screen.
func (m *ModalInput) Draw(screen tcell.Screen) {
	// Calculate the width of this modal.
	buttonsWidth := 50
	screenWidth, screenHeight := screen.Size()
	width := screenWidth / 3
	if width < buttonsWidth {
		width = buttonsWidth
	}
	// width is now without the box border.

	// Set the modal's position and size.
	height := 9
	width += 4
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)

	// Draw the frame.
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}
