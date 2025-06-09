package cmd

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestParsePrefix(t *testing.T) {
	color, text := parsePrefix("[blue]hello")
	if color != "blue" || text != "hello" {
		t.Fatalf("expected blue/hello got %s/%s", color, text)
	}
	color, text = parsePrefix("hello")
	if color != "" || text != "hello" {
		t.Fatalf("expected empty/hello got %s/%s", color, text)
	}
}

func TestApplyPrefix(t *testing.T) {
	if out := applyPrefix("hello", "blue"); out != "[blue]hello" {
		t.Fatalf("applyPrefix failed: %s", out)
	}
	if out := applyPrefix("[red]hello", ""); out != "hello" {
		t.Fatalf("applyPrefix remove failed: %s", out)
	}
}

func TestColorInputTabSwitch(t *testing.T) {
	ci := NewColorInput("Title:", []string{"default", "blue"})
	focused := tview.Primitive(nil)
	setFocus := func(p tview.Primitive) { focused = p; p.Focus(func(p tview.Primitive) {}) }

	ci.Focus(setFocus)
	if focused != ci.input {
		t.Fatalf("initial focus not on input")
	}
	ci.InputHandler()(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone), setFocus)
	if focused != ci.dropdown {
		t.Fatalf("tab did not move focus to dropdown")
	}
	ci.InputHandler()(tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone), setFocus)
	if focused != ci.input {
		t.Fatalf("backtab did not return focus to input")
	}
}

func TestColorInputFinishedFunc(t *testing.T) {
	ci := NewColorInput("Title:", []string{"default", "blue"})
	called := false
	ci.SetFinishedFunc(func(k tcell.Key) { called = true })

	ci.dropdown.InputHandler()(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone), func(p tview.Primitive) {})
	if !called {
		t.Fatalf("finished func not called for dropdown")
	}
}
