package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ModalInput is based on Modal from tview, but has an input field instead
type ModalInput struct {
	*tview.Form
	DialogHeight int
	frame        *tview.Frame
	main         string
	secondary    string
	due          string
	priority     int
	showPriority bool
	showDue      bool
	createdBy    string
	created      string
	updatedBy    string
	updated      string
	done         func(string, string, bool)
}

func NewModalInput(title string) *ModalInput {
	form := tview.NewForm()

	m := &ModalInput{form, 9, tview.NewFrame(form), "", "", "", 2, false, false, "", "", "", "", nil}

	form.SetCancelFunc(func() {
		if m.done != nil {
			m.done("", "", false)
		}
	})

	form.AddInputField("Task:", "", 50, nil, func(text string) {
		m.main = text
	})

	form.AddInputField("Created / due:", "", 50, nil, func(text string) {
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

func NewModalInputMode(title, modeDirectory string) *ModalInput {
	form := tview.NewForm()
	m := &ModalInput{form, 8, tview.NewFrame(form), "", "", "", 2, false, false, "", "", "", "", nil}

	form.AddInputField("Mode:", "", 50, nil, func(text string) {
		m.main = text
	})

	m.frame.AddText("New modes are saved in '"+modeDirectory+"'", false, 0, tcell.ColorDarkGray)

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
		SetBorderPadding(1, 0, 1, 1)

	return m
}

func NewModalInputLane(title, laneDescription string, dialogHeight int, initialInput1 string) *ModalInput {
	form := tview.NewForm()
	m := &ModalInput{form, dialogHeight, tview.NewFrame(form), "", "", "", 2, false, false, "", "", "", "", nil}
	m.main = initialInput1

	form.AddInputField("Lane:", initialInput1, 50, nil, func(text string) {
		m.main = text
	})

	m.frame.AddText(laneDescription, false, 0, tcell.ColorDarkGray)

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
		SetBorderPadding(1, 0, 1, 1)

	return m
}

// SetValue sets the current value in the item
func (m *ModalInput) SetValue(text string, secondary string, due string) {
	m.main = text
	m.secondary = secondary
	m.due = due
	m.Clear(false)
	height := 9
	m.AddInputField("", text, 50, nil, func(text string) {
		if len(text) == 0 {
			text = "(empty)"
		}
		m.main = text
	})
	m.AddInputField("", secondary, 50, nil, func(text string) {
		m.secondary = text
	})
	if m.showDue {
		dateField := tview.NewInputField().SetLabel("Due:").SetFieldWidth(10)
		updating := false
		dateField.SetAcceptanceFunc(func(text string, ch rune) bool {
			if ch == 0 {
				return true
			}
			if ch < '0' || ch > '9' {
				return false
			}
			digits := strings.ReplaceAll(text, ".", "")
			return len(digits) < 8
		})
		dateField.SetChangedFunc(func(text string) {
			if updating {
				return
			}
			digits := strings.ReplaceAll(text, ".", "")
			if len(digits) > 8 {
				digits = digits[:8]
			}
			formatted := digits
			if len(digits) > 4 {
				formatted = digits[:2] + "." + digits[2:4] + "." + digits[4:]
			} else if len(digits) > 2 {
				formatted = digits[:2] + "." + digits[2:]
			}
			if formatted != text {
				updating = true
				dateField.SetText(formatted)
				updating = false
			}
			m.due = formatted
		})
		dateField.SetText(due)
		m.AddFormItem(dateField)
		height++
	}
	if m.showPriority {
		options := []string{"1 (high)", "2 (normal)", "3 (low)", "4 (idle)"}
		m.AddDropDown("Priority:", options, m.priority-1, func(option string, index int) {
			m.priority = index + 1
		})
		height++
	}
	if m.createdBy != "" && m.created != "" {
		txt := fmt.Sprintf("created by: %s (%s)", m.createdBy, m.created)
		m.AddTextView("", txt, 50, 1, false, false)
		height++
	}
	if m.updatedBy != "" && m.updated != "" {
		txt := fmt.Sprintf("modified by: %s (%s)", m.updatedBy, m.updated)
		m.AddTextView("", txt, 50, 1, false, false)
		height++
	}
	m.DialogHeight = height
}

// SetPriority enables a priority dropdown with the given value (1-4).
func (m *ModalInput) SetPriority(value int) {
	if value < 1 || value > 4 {
		value = 2
	}
	m.priority = value
	m.showPriority = true
}

// SetDue enables a due date input field with the given value (dd.mm.yyyy).
func (m *ModalInput) SetDue(date string) {
	m.due = date
	m.showDue = true
}

// GetDueISO returns the due date in ISO format or empty if not set.
func (m *ModalInput) GetDueISO() string {
	if !m.showDue || m.due == "" {
		return ""
	}
	t, err := time.Parse("02.01.2006", m.due)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// GetPriority returns the currently selected priority value.
func (m *ModalInput) GetPriority() int {
	if !m.showPriority {
		return 0
	}
	return m.priority
}

// SetInfo sets the creation and modification information to be shown.
func (m *ModalInput) SetInfo(createdBy, created, updatedBy, updated string) {
	m.createdBy = createdBy
	m.created = created
	m.updatedBy = updatedBy
	m.updated = updated
}

// ClearExtras disables priority dropdown and info lines.
func (m *ModalInput) ClearExtras() {
	m.showPriority = false
	m.showDue = false
	m.createdBy = ""
	m.created = ""
	m.updatedBy = ""
	m.updated = ""
	m.due = ""
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
	height := m.DialogHeight
	width += 4
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)

	// Draw the frame.
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}
