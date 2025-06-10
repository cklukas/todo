package ui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// weekDateLayout returns a date layout without year including weekday.
func weekDateLayout() string {
	if localeUS() {
		return "Mon Jan 2"
	}
	return "Mon 2 Jan"
}

// ClockText returns the display string for the clock based on the screen width.
// If width is too small, an empty string is returned. If width is large enough,
// the date is included before the time.
func ClockText(now time.Time, width, base int) string {
	if width >= base+20 {
		return now.Format(weekDateLayout() + " 15:04:05")
	}
	if width >= base+9 {
		return now.Format("15:04:05")
	}
	return ""
}

// StartClock updates the given TextView every second with time and date based
// on the available screen width.
func StartClock(app *tview.Application, view *tview.TextView, lanes *Lanes, mode string) {
	ticker := time.NewTicker(time.Second)
	var width int
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		width, _ = screen.Size()
		return false
	})
	go func() {
		for now := range ticker.C {
			helpLen := lanes.MoveHelpLength()
			base := ClockBaseWidth(mode, helpLen)
			app.QueueUpdateDraw(func() {
				view.SetText(ClockText(now, width, base))
			})
		}
	}()
}

func ClockBaseWidth(mode string, moveHelpLen int) int {
	return 10 + 13 + 9 + 9 + 13 + 10 + 9 + 10 + 2 + len(mode) + moveHelpLen
}
