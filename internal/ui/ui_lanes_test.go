package ui

import (
	"testing"

	"github.com/cklukas/todo/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestNoSelectionChangeDuringEdit(t *testing.T) {
	c := &model.ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "task 1", "", 2, "", "")
	c.AddItem(0, 1, "task 2", "", 2, "", "")
	app := tview.NewApplication()
	l := NewLanes(c, app, "", t.TempDir(), "")

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
	c := &model.ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "task 1", "", 2, "", "")
	c.AddItem(0, 1, "task 2", "", 2, "", "")
	app := tview.NewApplication()
	l := NewLanes(c, app, "", t.TempDir(), "")

	l.CmdAddTask()

	event := tcell.NewEventMouse(50, 50, tcell.Button1, 0)
	l.pages.MouseHandler()(tview.MouseLeftClick, event, func(p tview.Primitive) {})

	if l.lanes[0].GetCurrentItem() != 0 {
		t.Fatalf("selection changed while adding")
	}
}

func TestAddTaskFocus(t *testing.T) {
	c := &model.ToDoContent{}
	c.InitializeNew()
	app := tview.NewApplication()
	l := NewLanes(c, app, "", t.TempDir(), "")

	l.CmdAddTask()

	if app.GetFocus() != l.add.titleField {
		t.Fatalf("focus not on add task dialog")
	}
}

func TestEditTaskFocus(t *testing.T) {
	c := &model.ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "task", "", 2, "", "")
	app := tview.NewApplication()
	l := NewLanes(c, app, "", t.TempDir(), "")

	l.CmdEditTask()

	if app.GetFocus() != l.edit.titleField {
		t.Fatalf("focus not on edit task dialog")
	}
}

func TestMouseCaptureDuringAdd(t *testing.T) {
	c := &model.ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "task 1", "", 2, "", "")
	app := tview.NewApplication()
	called := false
	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		called = true
		return event, action
	})
	l := NewLanes(c, app, "", t.TempDir(), "")

	l.CmdAddTask()
	event := tcell.NewEventMouse(50, 50, tcell.Button1, 0)
	ev, _ := app.GetMouseCapture()(event, tview.MouseLeftClick)
	if ev != nil {
		t.Fatalf("mouse event not swallowed")
	}
	if called {
		t.Fatalf("original capture called during dialog")
	}

	l.add.done("", "", false)
	ev, _ = app.GetMouseCapture()(event, tview.MouseLeftClick)
	if ev == nil {
		t.Fatalf("mouse event swallowed after dialog")
	}
}

func TestInputCaptureDuringAdd(t *testing.T) {
	c := &model.ToDoContent{}
	c.InitializeNew()
	app := tview.NewApplication()
	called := false
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		called = true
		return nil
	})
	l := NewLanes(c, app, "", t.TempDir(), "")

	l.CmdAddTask()
	key := tcell.NewEventKey(tcell.KeyF1, rune(0), tcell.ModNone)
	ret := app.GetInputCapture()(key)
	if ret != key {
		t.Fatalf("input capture modified event")
	}
	if called {
		t.Fatalf("original input capture called during dialog")
	}

	l.add.done("", "", false)
	called = false
	ret = app.GetInputCapture()(key)
	if ret != nil {
		t.Fatalf("expected nil from original capture")
	}
	if !called {
		t.Fatalf("original input capture not called after dialog")
	}
}

func TestInputCaptureDuringLaneDialog(t *testing.T) {
	c := &model.ToDoContent{}
	c.InitializeNew()
	app := tview.NewApplication()
	called := false
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		called = true
		return nil
	})
	l := NewLanes(c, app, "", t.TempDir(), "")

	l.CmdLanesCmds()
	key := tcell.NewEventKey(tcell.KeyF1, rune(0), tcell.ModNone)
	ret := app.GetInputCapture()(key)
	if ret != key {
		t.Fatalf("input capture modified event")
	}
	if called {
		t.Fatalf("original input capture called during dialog")
	}

	l.hideDialog("laneDialog")
	called = false
	ret = app.GetInputCapture()(key)
	if ret != nil {
		t.Fatalf("expected nil from original capture")
	}
	if !called {
		t.Fatalf("original input capture not called after dialog")
	}
}
