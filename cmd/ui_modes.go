package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
)

func (l *Lanes) ListValidModesRemoveProvided(activeMode string) ([]string, int, error) {
	res := make([]string, 0)
	modes, activeModeIndex, _ := l.ListValidModes(activeMode)
	for _, m := range modes {
		if m != activeMode {
			res = append(res, m)
		}
	}
	return res, activeModeIndex, nil
}

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
		if _, err := os.Stat(filepath.Join(l.todoDirModes, di.Name(), "todo.json")); errors.Is(err, os.ErrNotExist) {
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
		l.app.Stop()
		log.Fatal(err)
	}
	modePage := tview.NewModal().
		SetTitle(" Mode Selection ").
		// SetText("Modes allow separation of ToDo lists into categories.\n\nSelect an existing mode, 'Add' a new, or 'Merge/Remove' the current mode (tasks of all lanes are moved to the lanes of another mode or archived)").
		SetText("Modes allow separation of ToDo lists into categories. Select an existing mode from the list:").
		// AddButtons(append(modes, "Add", "Merge/Remove", "Cancel")).
		AddButtons(append(modes, "Cancel")).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			l.cmdSelectModeDialogAction(buttonIndex, buttonLabel, lastIndex)
		})

	modePage.SetFocus(activeModeIndex)
	l.setActive()
	l.pages.RemovePage("mode")
	l.pages.AddPage("mode", modePage, false, false)
	l.pages.ShowPage("mode")
}

func (l *Lanes) CmdRemoveModeDialog() {
	lastIndex := l.saveActive()
	modes, _, err := l.ListValidModesRemoveProvided(l.mode)
	if err != nil {
		l.app.Stop()
		log.Fatal(err)
	}
	modePage := tview.NewModal().
		SetTitle(" Remove Mode ").
		SetText(
			fmt.Sprintf("CAUTION: About to [red]delete[white] mode '%v'.\n\nTasks of the current mode will be moved into the respective lanes of the selected target mode, or can be removed by selecting 'Archive'.\n\nOther instances of this application showing this mode will close themself.", l.mode)).
		AddButtons(append(modes, "Archive", "Cancel")).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			l.cmdRemoveModeAction(buttonIndex, buttonLabel, lastIndex)
		})

	modePage.SetFocus(len(modes) + 1)
	l.setActive()
	l.pages.RemovePage("removeMode")
	l.pages.AddPage("removeMode", modePage, false, false)
	l.pages.ShowPage("removeMode")
}

func (l *Lanes) cmdRemoveModeAction(buttonIndex int, buttonLabel string, lastIndex int) {
	l.pages.RemovePage("removeMode")
}

func (l *Lanes) cmdSelectModeDialogAction(buttonIndex int, buttonLabel string, lastIndex int) {
	if buttonIndex >= 0 {
		switch buttonLabel {
		case "Add":
			l.pages.HidePage("mode")
			l.pages.ShowPage("addMode")
			return
		case "Merge/Remove":
			l.pages.HidePage("mode")
			l.CmdRemoveModeDialog()
			return
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
