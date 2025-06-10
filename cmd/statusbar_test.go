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
