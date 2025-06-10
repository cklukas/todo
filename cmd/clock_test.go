package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/cklukas/todo/internal/ui"
)

func TestClockText(t *testing.T) {
	base := ui.ClockBaseWidth("work", 0)
	ts := time.Date(2025, time.June, 10, 9, 8, 7, 0, time.Local)

	os.Setenv("LC_TIME", "en_US")
	if got := ui.ClockText(ts, base+5, base); got != "" {
		t.Fatalf("got %q want empty", got)
	}
	if got := ui.ClockText(ts, base+12, base); got != "09:08:07" {
		t.Fatalf("got %q want %q", got, "09:08:07")
	}
	if got := ui.ClockText(ts, base+25, base); got != "Tue Jun 10 09:08:07" {
		t.Fatalf("got %q", got)
	}

	os.Setenv("LC_TIME", "de_DE")
	if got := ui.ClockText(ts, base+25, base); got != "Tue 10 Jun 09:08:07" {
		t.Fatalf("got %q", got)
	}
	os.Unsetenv("LC_TIME")
}

func TestClockBaseWidth(t *testing.T) {
	noHelp := ui.ClockBaseWidth("work", 0)
	help := ui.ClockBaseWidth("work", len("Use Arrow Keys to Move Task"))
	if help <= noHelp {
		t.Fatalf("base width with help %d <= without %d", help, noHelp)
	}
}
