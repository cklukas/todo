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
	color        string
	showColor    bool
	colors       []string
	createdBy    string
	created      string
	updatedBy    string
	updated      string
	titleField   *tview.InputField
	okButton     *tview.Button
	done         func(string, string, bool)
}

func (m *ModalInput) GetFrame() *tview.Frame {
	return m.frame
}

func (m *ModalInput) updateOKButton() {
	if m.okButton == nil {
		return
	}
	if len(m.main) == 0 {
		m.okButton.SetLabelColor(tcell.ColorDarkGray)
		m.okButton.SetSelectedFunc(func() {})
	} else {
		m.okButton.SetLabelColor(tview.Styles.PrimaryTextColor)
		m.okButton.SetSelectedFunc(func() {
			if m.done != nil {
				m.done(m.main, m.secondary, true)
			}
		})
	}
}

func NewModalInput(title string) *ModalInput {
	form := tview.NewForm()

	m := &ModalInput{Form: form, DialogHeight: 9, frame: tview.NewFrame(form), main: "", secondary: "", due: "", priority: 2, showPriority: false, showDue: false, color: "", showColor: false, colors: []string{"default", "blue", "green", "red", "yellow"}, createdBy: "", created: "", updatedBy: "", updated: "", titleField: nil, okButton: nil, done: nil}

	form.SetCancelFunc(func() {
		if m.done != nil {
			m.done("", "", false)
		}
	})

	var titleField *tview.InputField
	form, titleField = form.AddInputField("Title:", "", 50, nil, func(text string) {
		m.main = text
		m.updateOKButton()
	})
	m.titleField = titleField

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
	m.okButton = m.GetButton(0)
	m.updateOKButton()
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
	m := &ModalInput{Form: form, DialogHeight: 8, frame: tview.NewFrame(form), main: "", secondary: "", due: "", priority: 2, showPriority: false, showDue: false, createdBy: "", created: "", updatedBy: "", updated: "", titleField: nil, okButton: nil, done: nil}

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
	m := &ModalInput{Form: form, DialogHeight: dialogHeight, frame: tview.NewFrame(form), main: "", secondary: "", due: "", priority: 2, showPriority: false, showDue: false, createdBy: "", created: "", updatedBy: "", updated: "", titleField: nil, okButton: nil, done: nil}
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

	var titleField *tview.InputField
	if m.showColor {
		ci := NewColorInput("Title:", m.colors)
		titleField = ci.input
		titleField.SetText(text)
		titleField.SetChangedFunc(func(text string) {
			if len(text) == 0 {
				text = "(empty)"
			}
			m.main = text
			m.updateOKButton()
		})
		idx := 0
		for i, c := range m.colors {
			if c == m.color {
				idx = i
				break
			}
		}
		ci.dropdown.SetCurrentOption(idx)
		ci.dropdown.SetSelectedFunc(func(option string, index int) {
			m.color = m.colors[index]
		})
		m.AddFormItem(ci)
	} else {
		m.Form, titleField = m.Form.AddInputField("Title:", text, 50, nil, func(text string) {
			if len(text) == 0 {
				text = "(empty)"
			}
			m.main = text
			m.updateOKButton()
		})
	}
	m.titleField = titleField
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
		dateField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				newText := removeLastDueDigit(m.due)
				updating = true
				dateField.SetText(newText)
				updating = false
				m.due = newText
				return nil
			}
			return event
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
		txt := fmt.Sprintf("%s (%s)", m.created, m.createdBy)
		tv := tview.NewTextView().SetLabel("Created:").SetSize(1, 50).SetText(txt).SetScrollable(false)
		tv.SetTextColor(tcell.ColorDarkGray)
		m.AddFormItem(tv)
	}
	if m.updatedBy != "" && m.updated != "" {
		txt := fmt.Sprintf("%s (%s)", m.updated, m.updatedBy)
		tv := tview.NewTextView().SetLabel("Modified:").SetSize(1, 50).SetText(txt).SetScrollable(false)
		tv.SetTextColor(tcell.ColorDarkGray)
		m.AddFormItem(tv)
	}

	itemCount := m.GetFormItemCount()
	m.DialogHeight = 2*itemCount + 5
	m.updateOKButton()
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

// removeLastDueDigit removes the last numeric digit from a formatted due date.
func removeLastDueDigit(text string) string {
	digits := strings.ReplaceAll(strings.ReplaceAll(text, ".", ""), "/", "")
	if len(digits) == 0 {
		return ""
	}
	digits = digits[:len(digits)-1]
	return formatDueInput(digits)
}

// SetPriority enables a priority dropdown with the given value (1-4).
func (m *ModalInput) SetPriority(value int) {
	if value < 1 || value > 4 {
		value = 2
	}
	m.priority = value
	m.showPriority = true
}

// SetColor enables a color dropdown with the given value.
func (m *ModalInput) SetColor(color string) {
	m.color = color
	m.showColor = true
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

// GetColor returns the currently selected color value.
func (m *ModalInput) GetColor() string {
	if !m.showColor {
		return ""
	}
	return m.color
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
	m.showColor = false
	m.color = ""
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

// Focus delegates focus to the internal form.
func (m *ModalInput) Focus(delegate func(p tview.Primitive)) {
	delegate(m.Form)
}

// HasFocus returns whether the form has focus.
func (m *ModalInput) HasFocus() bool {
	return m.Form.HasFocus()
}

// SetFocus passes the focus index to the embedded form.
func (m *ModalInput) SetFocus(index int) *ModalInput {
	m.Form.SetFocus(index)
	return m
}

// MouseHandler forwards mouse events to the form and captures clicks inside the dialog.
func (m *ModalInput) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
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
func (m *ModalInput) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if m.frame.HasFocus() {
			if handler := m.frame.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
		}
	})
}
