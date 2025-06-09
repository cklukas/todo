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

	form.AddInputField("Details:", "", 50, nil, func(text string) {
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

	m.AddInputField("Task:", text, 50, nil, func(text string) {
		if len(text) == 0 {
			text = "(empty)"
		}
		m.main = text
	})
	m.AddInputField("Details:", secondary, 50, nil, func(text string) {
		m.secondary = text
	})
	if m.showDue {
		dateField := tview.NewInputField().SetLabel("Due:").SetFieldWidth(10).SetPlaceholder(duePlaceholder())
		updating := false
		dateField.SetAcceptanceFunc(func(text string, ch rune) bool {
			if ch == 0 {
				return true
			}
			if ch < '0' || ch > '9' {
				return false
			}
			digits := strings.ReplaceAll(strings.ReplaceAll(text, ".", ""), "/", "")
			return len(digits) <= 8
		})
		dateField.SetChangedFunc(func(text string) {
			if updating {
				return
			}
			formatted := formatDueInput(text)
			if formatted != text {
				updating = true
				dateField.SetText(formatted)
				updating = false
			}
			m.due = formatted
		})
		dateField.SetText(due)
		m.AddFormItem(dateField)
	}
	if m.showPriority {
		options := []string{"1 (high)", "2 (normal)", "3 (low)", "4 (idle)"}
		m.AddDropDown("Priority:", options, m.priority-1, func(option string, index int) {
			m.priority = index + 1
		})
	}
	if m.createdBy != "" && m.created != "" {
		tv := tview.NewTextView().SetLabel("Created by:").SetSize(1, 50).SetText(m.createdBy).SetScrollable(false)
		tv.SetTextColor(tcell.ColorDarkGray)
		m.AddFormItem(tv)
		tv2 := tview.NewTextView().SetLabel("Created at:").SetSize(1, 50).SetText(m.created).SetScrollable(false)
		tv2.SetTextColor(tcell.ColorDarkGray)
		m.AddFormItem(tv2)
	}
	if m.updatedBy != "" && m.updated != "" {
		tv := tview.NewTextView().SetLabel("Modified by:").SetSize(1, 50).SetText(m.updatedBy).SetScrollable(false)
		tv.SetTextColor(tcell.ColorDarkGray)
		m.AddFormItem(tv)
		tv2 := tview.NewTextView().SetLabel("Modified at:").SetSize(1, 50).SetText(m.updated).SetScrollable(false)
		tv2.SetTextColor(tcell.ColorDarkGray)
		m.AddFormItem(tv2)
	}

	itemCount := m.GetFormItemCount()
	m.DialogHeight = 2*itemCount + 5
}

// formatDueInput formats a date input so that dots are inserted after day and
// month. Input may already contain dots and may be partially entered.
// At most 8 digits are considered.
func formatDueInput(text string) string {
	digits := strings.ReplaceAll(strings.ReplaceAll(text, ".", ""), "/", "")
	if len(digits) > 8 {
		digits = digits[:8]
	}
	if localeUS() {
		switch {
		case len(digits) > 4:
			return digits[:2] + "/" + digits[2:4] + "/" + digits[4:]
		case len(digits) == 4:
			return digits[:2] + "/" + digits[2:4] + "/"
		case len(digits) >= 2:
			return digits[:2] + "/" + digits[2:]
		default:
			return digits
		}
	}
	switch {
	case len(digits) > 4:
		return digits[:2] + "." + digits[2:4] + "." + digits[4:]
	case len(digits) == 4:
		return digits[:2] + "." + digits[2:4] + "."
	case len(digits) >= 2:
		return digits[:2] + "." + digits[2:]
	default:
		return digits
	}
}

// SetPriority enables a priority dropdown with the given value (1-4).
func (m *ModalInput) SetPriority(value int) {
	if value < 1 || value > 4 {
		value = 2
	}
	m.priority = value
	m.showPriority = true
}

// SetDue enables a due date input field with the given value using locale formatting.
func (m *ModalInput) SetDue(date string) {
	m.due = date
	m.showDue = true
}

// GetDueISO returns the due date in ISO format or empty if not set.
func (m *ModalInput) GetDueISO() string {
	if !m.showDue || m.due == "" {
		return ""
	}
	t, err := time.Parse(dueLayout(), m.due)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// DueValid reports whether the current due value can be parsed.
func (m *ModalInput) DueValid() bool {
	if !m.showDue || m.due == "" {
		return true
	}
	_, err := time.Parse(dueLayout(), m.due)
	return err == nil
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
