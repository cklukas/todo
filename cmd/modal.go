package cmd

import "github.com/rivo/tview"

// modal returns a primitive which places the provided primitive in the center
// of the screen with the given width and height. This follows the approach
// outlined in the tview Modal documentation and ensures mouse events outside
// the dialog are captured by the surrounding Flex layout.
func modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
