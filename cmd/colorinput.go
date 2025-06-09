package cmd

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ColorInput is a form item that combines an input field with a color drop-down.
type ColorInput struct {
	*tview.Flex
	input    *tview.InputField
	dropdown *tview.DropDown
	colors   []string
}

// NewColorInput creates a new ColorInput. colors is a slice of color names. The
// first entry should be "default".
func NewColorInput(label string, colors []string) *ColorInput {
	input := tview.NewInputField().SetLabel(label)
	dd := tview.NewDropDown()
	dd.SetLabel("")
	disp := make([]string, len(colors))
	for i, c := range colors {
		if i == 0 || c == "" {
			disp[i] = "default"
		} else {
			disp[i] = "[" + c + "]" + c
		}
	}
	dd.SetOptions(disp, nil)
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.AddItem(input, 0, 4, true)
	flex.AddItem(dd, 0, 1, false)
	return &ColorInput{flex, input, dd, colors}
}

// GetLabel returns the item's label text.
func (c *ColorInput) GetLabel() string {
	return c.input.GetLabel()
}

// SetFormAttributes sets form attributes for both widgets.
func (c *ColorInput) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	c.input.SetFormAttributes(labelWidth, labelColor, bgColor, fieldTextColor, fieldBgColor)
	c.dropdown.SetFormAttributes(0, labelColor, bgColor, fieldTextColor, fieldBgColor)
	return c
}

// GetFieldWidth returns 0 to use available width.
func (c *ColorInput) GetFieldWidth() int {
	return 0
}

// GetFieldHeight returns the field height.
func (c *ColorInput) GetFieldHeight() int {
	return 1
}

// SetFinishedFunc installs a handler for the input field.
func (c *ColorInput) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	c.input.SetFinishedFunc(handler)
	c.dropdown.SetDoneFunc(handler)
	return c
}

// InputHandler handles tab navigation between input field and dropdown.
func (c *ColorInput) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return c.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if c.input.HasFocus() && event.Key() == tcell.KeyTab {
			setFocus(c.dropdown)
			return
		}
		if c.dropdown.HasFocus() && event.Key() == tcell.KeyBacktab {
			setFocus(c.input)
			return
		}
		if handler := c.Flex.InputHandler(); handler != nil {
			handler(event, setFocus)
		}
	})
}

// removePrefix returns text without a leading color prefix.
func removePrefix(text string) string {
	if strings.HasPrefix(text, "[") {
		if idx := strings.Index(text, "]"); idx > 0 {
			return text[idx+1:]
		}
	}
	return text
}

// applyPrefix adds or replaces the color prefix.
func applyPrefix(text, color string) string {
	base := removePrefix(text)
	if color == "" || color == "default" {
		return base
	}
	return "[" + color + "]" + base
}

// parsePrefix splits a title into color prefix and the remaining text.
func parsePrefix(text string) (string, string) {
	if strings.HasPrefix(text, "[") {
		if idx := strings.Index(text, "]"); idx > 0 {
			return strings.ToLower(text[1:idx]), text[idx+1:]
		}
	}
	return "", text
}
