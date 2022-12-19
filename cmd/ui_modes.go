package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/rivo/tview"
)

func (l *Lanes) ListValidModes(activeMode string) ([]string, int, error) {
	activeModeIndex := 0
	modes := make([]string, 0)
	modes = append(modes, "main")
	dirEntries, err := os.ReadDir(l.todoDirModes)
	if err != nil {
		return modes, 0, err
	}

	idx := 0
	for _, di := range dirEntries {
		if di.Name() == "main" {
			continue
		}
		if di.IsDir() && !strings.HasPrefix(di.Name(), ".") {
			modes = append(modes, di.Name())
			idx++
			if di.Name() == activeMode {
				activeModeIndex = idx
			}
		}
	}

	return modes, activeModeIndex, nil
}

func (l *Lanes) CmdSelectModeDialog() {
	lastIndex := l.saveActive()
	modes, activeModeIndex, err := l.ListValidModes(l.mode)
	if err != nil {
		log.Fatal(err)
	}
	modePage := tview.NewModal().
		SetTitle(" Mode Selection ").
		SetText("Select mode, add a new, or remove the current mode (tasks of all lanes are moved to the lanes of another mode or archived)").
		AddButtons(modes).
		AddButtons([]string{"Add", "Merge/Remove", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			l.cmdSelectModeDialogAction(buttonIndex, buttonLabel, lastIndex)
		})

	modePage.SetFocus(activeModeIndex)
	l.setActive()
	l.pages.RemovePage("mode")
	l.pages.AddPage("mode", modePage, false, false)
	l.pages.ShowPage("mode")
}

func (l *Lanes) cmdSelectModeDialogAction(buttonIndex int, buttonLabel string, lastIndex int) {
	if buttonIndex >= 0 {
		switch buttonLabel {
		case "Add":
			l.pages.HidePage("mode")
			l.pages.ShowPage("addMode")
			return
		case "Merge/Remove":
			// empty
		case "Cancel":
			// empty
		default:
			l.nextMode = buttonLabel
			l.app.Stop()
		}
	}
	l.pages.HidePage("mode")
	l.setActiveIndex(lastIndex)
}

func (l *Lanes) CmdAddMode() {
	l.saveActive()
	l.pages.ShowPage("addMode")
}
