package cmd

import (
	"fmt"
	"log"
	"os/user"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (l *Lanes) GetActiveLaneName() string {
	return l.content.Titles[l.active]
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (l *Lanes) RedrawAllLanes() {
	l.content.readWriteMutex.Lock()
	defer l.content.readWriteMutex.Unlock()

	for laneIdx := 0; laneIdx < min(len(l.lanes), len(l.content.Items)); laneIdx++ {
		if len(l.content.Items[laneIdx]) == 0 {
			continue
		}
		currentIndexInLine := l.lanes[laneIdx].GetCurrentItem()
		validIndexInLine := normPos(currentIndexInLine, len(l.content.Items[laneIdx]))
		l.redrawLane(laneIdx, validIndexInLine)
	}
}

func NewLanes(content *ToDoContent, app *tview.Application, mode, todoDirModes string) *Lanes {
	l := &Lanes{"", 0, todoDirModes, mode, content, make([]*tview.List, content.GetNumLanes()), 0, 0, false, tview.NewPages(), app, false,
		NewModalInput("Add Task"),
		NewModalInput("Edit Task"),
		NewModalInputMode("Add Mode", todoDirModes), nil}
	flex := tview.NewFlex()
	for i := 0; i < l.content.GetNumLanes(); i++ {
		l.lanes[i] = tview.NewList()
		xi := i
		l.lanes[i].SetFocusFunc(func() {
			l.lanes[xi].SetSelectedStyle(tcell.StyleDefault)
			l.active = xi
			l.lanes[xi].SetSelectedBackgroundColor(tcell.ColorLightBlue)
			l.lanes[xi].SetSelectedTextColor(tcell.ColorBlack)
			if l.lastActiveSaved {
				l.lastActiveSaved = false
				if l.lastActive > 0 {
					for i := 0; i < l.lastActive; i++ {
						l.incActive()
					}
				}
			}
		})
		l.lanes[i].ShowSecondaryText(true).SetBorder(true)
		l.lanes[i].SetTitle(l.content.GetLaneTitle(i))
		l.lanes[i].SetInputCapture(l.HotKeyHandler)
		l.lanes[i].SetSelectedFunc(func(w int, x string, y string, z rune) {
			if l.inselect {
				l.selected()
				content.Save()
			} else {
				l.selected()
			}
		})
		l.lanes[i].SetDoneFunc(func() {
			// Cancel select on Done (escape)
			if l.inselect {
				l.selected()
				content.Save()
			}
		})
		for _, item := range l.content.GetLaneItems(i) {
			l.lanes[i].AddItem(item.Title, item.Secondary, 0, nil)
		}
		flex.AddItem(l.lanes[i], 0, 1, i == 0)
	}
	l.pages.AddPage("lanes", flex, true, true)

	quit := tview.NewModal().
		SetText("Do you want to quit the application?").
		SetTitle(" Exit ").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			}
			l.pages.HidePage("quit")
		})
	l.pages.AddPage("quit", quit, false, false)

	// help := tview.NewModal().
	help := tview.NewModal()
	help = help.
		SetText("- developed by C. Klukas -\n\n- adapted from toukan (https://github.com/witchard/toukan) -\n\nUsage/Keys:\nEnter/space - mark task, cursor keys - move marked task, +/Insert - add, e - edit, Del/d - delete task, n - note, a - archive, Tab - switch lane, m - select mode, q - quit").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			l.pages.HidePage("help")
			l.setActive()
		})

	help.SetTitle(" About TODO ")

	l.pages.AddPage("help", help, false, false)

	delete := tview.NewModal().
		SetText("About to delete selected task. Continue?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				item := l.lanes[l.active].GetCurrentItem()
				l.content.DelItem(l.active, item)
				l.redrawLane(l.active, item)
				content.Save()
			}
			l.pages.HidePage("delete")
			l.setActive()
		}).SetTitle(" Delete Task ")

	l.pages.AddPage("delete", delete, false, false)

	archive := tview.NewModal().
		SetText("About to archive selected task. Continue?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				item := l.lanes[l.active].GetCurrentItem()
				err := l.content.ArchiveItem(l.active, item)
				if err != nil {
					app.Stop()
					log.Fatal(err)
				}
				l.redrawLane(l.active, item)
				content.Save()
			}
			l.pages.HidePage("archive")
			l.setActive()
		}).SetTitle(" Archive Task ")

	l.pages.AddPage("archive", archive, false, false)

	waitPage := tview.NewModal().
		SetText("When finished editing the note, save the changes and close Notepad. The item note text will be updated and you can continue to use the ToDo app.").
		SetTitle(" Editing Note ")

	l.pages.AddPage("wait", waitPage, false, false)

	l.add.SetDoneFunc(func(text string, secondary string, success bool) {
		if success {
			item := l.lanes[l.active].GetCurrentItem()
			if len(text) == 0 {
				text = "(empty)"
			}
			prio := l.add.GetPriority()
			due := l.add.GetDueISO()
			l.content.AddItem(l.active, item, text, secondary, prio, due)
			l.redrawLane(l.active, item)
			content.Save()
		}
		l.pages.HidePage("add")
		l.setActive()
	})
	l.pages.AddPage("add", l.add, false, false)

	l.addMode.SetDoneFunc(func(text string, secondary string, success bool) {
		l.pages.HidePage("addMode")
		l.setActive()
		if success {
			l.nextMode = text
			l.app.Stop()
		}
	})
	l.pages.AddPage("addMode", l.addMode, false, false)

	l.edit.SetDoneFunc(func(text string, secondary string, success bool) {
		if success {
			item := l.lanes[l.active].GetCurrentItem()
			itemVal := l.currentItem()
			itemVal.Title = text
			itemVal.Secondary = secondary
			itemVal.Priority = l.edit.GetPriority()
			itemVal.Due = l.edit.GetDueISO()
			itemVal.LastUpdate = time.Now().UTC().Format(time.RFC3339)
			if usr, err := user.Current(); err == nil {
				itemVal.UpdatedByName = usr.Username
			}
			l.redrawLane(l.active, item)
			l.content.Save()
		}
		l.pages.HidePage("edit")
		l.setActive()
	})
	l.pages.AddPage("edit", l.edit, false, false)

	return l
}

func (l *Lanes) CmdLanesCmds() {
	initActiveLane := l.saveActive()
	addToLeft := false
	lanePage := tview.NewModal().
		SetTitle(" Lane Commands ").
		SetText(fmt.Sprintf("Rename lane '%v', add a new lane, or remove it (tasks of current lane are moved to another lane or archived):", l.GetActiveLaneName())).
		AddButtons([]string{"Rename", "Add to left", "Add to right", "Merge/Remove", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Rename":
				l.renameLaneCommand(initActiveLane)
				return
			case "Add to left":
				addToLeft = true
				fallthrough
			case "Add to right":
				l.addLaneLeftRightCommand(addToLeft, initActiveLane)
				return
			case "Merge/Remove":
				l.removeMergeLaneDialog(initActiveLane)
				return
			case "Cancel":
				// empty
			default:
				// l.nextMode = buttonLabel
				// l.app.Stop()
			}
			l.pages.HidePage("laneDialog")
			l.setActiveIndex(initActiveLane)
		})

	lanePage.SetFocus(0)
	l.pages.AddPage("laneDialog", lanePage, false, true)
}

func (l *Lanes) renameLaneCommand(initActiveLane int) {
	addLaneDialog := NewModalInputLane("Rename Lane", "", 7, l.GetActiveLaneName())

	addLaneDialog.SetDoneFunc(func(lane, _ string, success bool) {
		l.pages.HidePage("addLane")
		l.setActiveIndex(initActiveLane)
		l.content.SetLaneTitle(initActiveLane, lane)
		l.content.Save()
		l.RedrawAllLanes()
	})

	l.pages.HidePage("laneDialog")
	l.pages.AddPage("addLane", addLaneDialog, false, true)
}

func (l *Lanes) addLaneLeftRightCommand(addToLeft bool, initActiveLane int) {
	leftRight := "right"
	if addToLeft {
		leftRight = "left"
	}

	addLaneDialog := NewModalInputLane(
		"Add Lane",
		fmt.Sprintf("New lane will be created %v of lane '%v'.", leftRight, l.GetActiveLaneName()), 8, "")

	addLaneDialog.SetDoneFunc(func(lane, _ string, success bool) {
		if success && len(lane) > 0 {
			laneIndex := l.content.InsertNewLane(addToLeft, lane, initActiveLane)
			l.content.Save()
			l.nextMode = l.mode
			l.nextLaneFocus = laneIndex
			l.app.Stop()
		}
		l.pages.HidePage("addLane")
		l.setActiveIndex(initActiveLane)
	})

	l.pages.HidePage("laneDialog")
	l.pages.AddPage("addLane", addLaneDialog, false, true)
}

func (l *Lanes) removeMergeLaneDialog(initActiveLane int) {
	activeLane := l.GetActiveLaneName()
	targetLanes := make([]string, 0)
	nItemsInActiveLane := len(l.content.GetLaneItems(initActiveLane))
	for _, title := range l.content.Titles {
		if title != activeLane {
			targetLanes = append(targetLanes, title)
		}
	}

	endS := "s"
	themIt := "them"
	if nItemsInActiveLane == 1 {
		endS = ""
		themIt = "it"
	}

	removeLaneDialog := tview.NewModal().
		SetTitle(" Remove Lane ").
		SetText(
			fmt.Sprintf("CAUTION: About to [red]delete[white] lane '%v'.\n\nSelect a target lane, where the %v task%v of the current lane will be moved to, or select 'Archive' to remove %v.",
				activeLane, nItemsInActiveLane, endS, themIt)).
		AddButtons(append(targetLanes, "Archive", "Cancel"))

	removeLaneDialog.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		l.removeMergeLaneCommand(buttonIndex, buttonLabel, targetLanes, initActiveLane)
	})

	l.pages.HidePage("laneDialog")
	l.pages.AddPage("removeLane", removeLaneDialog, false, true)
}

func (l *Lanes) removeMergeLaneCommand(buttonIndex int, buttonLabel string, targetLanes []string, initActiveLane int) {
	var removeLaneOK bool
	if buttonIndex < len(targetLanes) {
		// move tasks into target lane
		if buttonIndex >= initActiveLane {
			buttonIndex++
		}
		for {
			if len(l.content.Items[initActiveLane]) == 0 {
				break
			}
			l.content.MoveItem(initActiveLane, len(l.content.Items[initActiveLane])-1, buttonIndex, 0)
		}
		removeLaneOK = true
	} else {
		if buttonIndex == len(targetLanes) {
			// archive all tasks in line
			for {
				if len(l.content.Items[initActiveLane]) == 0 {
					break
				}
				err := l.content.ArchiveItem(initActiveLane, 0)
				if err != nil {
					l.app.Stop()
					log.Fatal(err)
				}
			}
			removeLaneOK = true
		} else {
			// Cancel
			removeLaneOK = false
		}
	}

	if removeLaneOK {
		l.content.RemoveLane(initActiveLane)
		l.content.Save()
		l.nextMode = l.mode
		l.nextLaneFocus = initActiveLane
		l.app.Stop()
		return
	}

	l.pages.HidePage("removeLane")
	l.setActiveIndex(initActiveLane)
}
