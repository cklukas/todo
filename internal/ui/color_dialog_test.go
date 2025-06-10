package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"testing"
)

func TestNewColorModalDefaultSelected(t *testing.T) {
	dlg := NewColorModal("Color", "")
	item := dlg.GetFormItem(0).(*tview.DropDown)
	_, text := item.GetCurrentOption()
	if text != "default" {
		t.Fatalf("expected initial option 'default', got '%s'", text)
	}
	if dlg.GetButtonIndex("Cancel") == -1 {
		t.Fatalf("cancel button missing")
	}
}

func TestColorModalDrawSetsFrame(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("screen init failed: %v", err)
	}
	dlg := NewColorModal("Color", "")
	dlg.Draw(screen)
	_, _, w, h := dlg.GetFrame().GetRect()
	if w == 0 || h == 0 {
		t.Fatalf("frame rect not set by Draw")
	}
}
