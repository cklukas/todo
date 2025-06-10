package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"testing"
)

func TestNewSortModalManualSelected(t *testing.T) {
	dlg := NewSortModal("Sort", "lane", "")
	item := dlg.GetFormItem(0).(*tview.DropDown)
	_, text := item.GetCurrentOption()
	if text != "manual" {
		t.Fatalf("expected initial option 'manual', got '%s'", text)
	}
	if dlg.GetButtonIndex("Cancel") == -1 {
		t.Fatalf("cancel button missing")
	}
}

func TestSortModalDrawSetsFrame(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("screen init failed: %v", err)
	}
	dlg := NewSortModal("Sort", "lane", "")
	dlg.Draw(screen)
	_, _, w, h := dlg.frame.GetRect()
	if w == 0 || h == 0 {
		t.Fatalf("frame rect not set by Draw")
	}
}
