package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ColorInput is a form item that combines an input field with a color drop-down.
type ColorInput struct {
	*tview.Flex
	input     *tview.InputField
	dropdown  *tview.DropDown
	colors    []string
	laneColor string
}

// NewColorInput creates a new ColorInput. colors is a slice of color names. The
// first entry should be "default".
func NewColorInput(label, laneColor string, colors []string) *ColorInput {
	input := tview.NewInputField().SetLabel(label)
	dd := tview.NewDropDown()
	dd.SetLabel("")
	disp := make([]string, len(colors))
	maxLen := 0
	for i, c := range colors {
		name := c
		if i == 0 || c == "" {
			name = "default"
		}
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}
	for i, c := range colors {
		name := c
		fg := c
		bg := laneColor
		if i == 0 || c == "" {
			name = "default"
			fg = "black"
		}
		style := "[" + fg
		if bg != "" {
			style += ":" + bg
		}
		style += "]"
		disp[i] = style + name + strings.Repeat(" ", maxLen-len(name))
	}
	dd.SetOptions(disp, nil)
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.AddItem(input, 0, 4, true)
	flex.AddItem(dd, 0, 1, false)
	return &ColorInput{flex, input, dd, colors, laneColor}
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
