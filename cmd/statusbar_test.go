package cmd

import (
	"testing"

	"github.com/cklukas/todo/internal/model"
	"github.com/cklukas/todo/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestStatusBarUsesGrayBackground(t *testing.T) {
	c := &model.ToDoContent{}
	c.InitializeNew()
	app := tview.NewApplication()
	lanes := ui.NewLanes(c, app, "main", t.TempDir(), "")

	status := getStatusBar(lanes, "main")

	if status.GetBackgroundColor() != tcell.ColorLightGray {
		t.Fatalf("status bar background color %v, want %v", status.GetBackgroundColor(), tcell.ColorLightGray)
	}
}

func TestStatusBarContainsClock(t *testing.T) {
	c := &model.ToDoContent{}
	c.InitializeNew()
	app := tview.NewApplication()
	lanes := ui.NewLanes(c, app, "main", t.TempDir(), "")

	status := getStatusBar(lanes, "main")

	if status.GetItemCount() != 11 {
		t.Fatalf("status bar item count %d, want %d", status.GetItemCount(), 11)
	}
	if _, ok := status.GetItem(10).(*tview.TextView); !ok {
		t.Fatalf("last status bar item should be TextView")
	}
}
