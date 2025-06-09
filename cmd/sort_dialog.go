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

	labels := []string{"default", "color", "due", "created", "modified", "priority"}
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
