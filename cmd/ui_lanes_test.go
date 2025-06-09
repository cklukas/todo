package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"testing"
)

func TestNoSelectionChangeDuringEdit(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "task 1", "", 2, "")
	c.AddItem(0, 1, "task 2", "", 2, "")
	app := tview.NewApplication()
	l := NewLanes(c, app, "", t.TempDir())

	l.CmdEditTask()

	x, y, _, _ := l.lanes[0].GetInnerRect()
	// position of second item
	event := tcell.NewEventMouse(x, y+1, tcell.Button1, 0)
	l.pages.MouseHandler()(tview.MouseLeftClick, event, func(p tview.Primitive) {})

	if l.lanes[0].GetCurrentItem() != 0 {
		t.Fatalf("selection changed while editing")
	}
}

func TestNoSelectionChangeDuringAdd(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "task 1", "", 2, "")
	c.AddItem(0, 1, "task 2", "", 2, "")
	app := tview.NewApplication()
	l := NewLanes(c, app, "", t.TempDir())

	l.CmdAddTask()

	x, y, _, _ := l.lanes[0].GetInnerRect()
	event := tcell.NewEventMouse(x, y+1, tcell.Button1, 0)
	l.pages.MouseHandler()(tview.MouseLeftClick, event, func(p tview.Primitive) {})

	if l.lanes[0].GetCurrentItem() != 0 {
		t.Fatalf("selection changed while adding")
	}
}

func TestAddTaskFocus(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	app := tview.NewApplication()
	l := NewLanes(c, app, "", t.TempDir())

	l.CmdAddTask()

	if app.GetFocus() != l.add.titleField {
		t.Fatalf("focus not on add task dialog")
	}
}

func TestEditTaskFocus(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "task", "", 2, "")
	app := tview.NewApplication()
	l := NewLanes(c, app, "", t.TempDir())

	l.CmdEditTask()

	if app.GetFocus() != l.edit.titleField {
		t.Fatalf("focus not on edit task dialog")
	}
}
