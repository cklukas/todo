package cmd

import "github.com/rivo/tview"

// modal wraps the provided primitive in an overlay centered on the screen.
// The surrounding boxes make sure that mouse events outside the modal are
// consumed so that underlying pages cannot receive them.
func modal(p tview.Primitive) tview.Primitive {
	return tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(p, 0, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), 0, 1, true).
		AddItem(tview.NewBox(), 0, 1, false)
}
