package cmd

import (
	"fmt"
	"strings"

	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Lanes struct {
	nextMode        string
	todoDirModes    string
	mode            string
	content         *Content
	lanes           []*tview.List
	active          int
	lastActive      int
	lastActiveSaved bool
	pages           *tview.Pages
	app             *tview.Application
	inselect        bool
	add             *ModalInput
	edit            *ModalInput
	addMode         *ModalInput

	bMoveHelp *tview.Button
}

func (l *Lanes) CmdAbout() {
	l.saveActive()
	l.pages.ShowPage("help")
}

func (l *Lanes) CmdExit() {
	l.setActive()
	l.pages.ShowPage("quit")
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

func (l *Lanes) CmdSelectMode() {
	modes, activeModeIndex, err := l.ListValidModes(l.mode)
	if err != nil {
		log.Fatal(err)
	}
	modePage := tview.NewModal().
		SetText("Select mode").
		AddButtons(modes).
		AddButtons([]string{"Add", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Add":
				l.pages.HidePage("mode")
				l.pages.ShowPage("addMode")
				return
			case "Cancel":
				// empty
			default:
				l.nextMode = buttonLabel
				l.app.Stop()
			}
			l.pages.HidePage("mode")
		})

	modePage.SetFocus(activeModeIndex)
	l.setActive()
	l.pages.RemovePage("mode")
	l.pages.AddPage("mode", modePage, false, false)
	l.pages.ShowPage("mode")
}

func (l *Lanes) CmdAddTask() {
	l.saveActive()
	now := time.Now()
	l.add.SetValue("", fmt.Sprintf("created: %v", now.Format("2006-01-02")))
	l.pages.ShowPage("add")
}

func (l *Lanes) CmdEditTask() {
	l.saveActive()
	if item := l.currentItem(); item != nil {
		l.edit.SetValue(item.Title, item.Secondary)
		l.pages.ShowPage("edit")
	}
}

func (l *Lanes) CmdAddMode() {
	l.saveActive()
	l.pages.ShowPage("addMode")
}

func (l *Lanes) CmdEditNote() {
	l.saveActive()
	if runtime.GOOS == "windows" {
		l.pages.ShowPage("wait")
		l.app.ForceDraw()
		l.editNote()
		l.pages.HidePage("wait")
	} else {
		l.app.Suspend(l.editNote)
	}
}

func (l *Lanes) CmdArchiveNote() {
	l.saveActive()
	l.pages.ShowPage("archive")
}

func (l *Lanes) CmdSelectNote() {
	l.setActive()
	l.selected()
}

func (l *Lanes) Focus() {
	l.setActive()
	l.selected()
}

func NewLanes(content *Content, app *tview.Application, mode, todoDirModes string) *Lanes {
	l := &Lanes{"", todoDirModes, mode, content, make([]*tview.List, content.GetNumLanes()), 0, 0, false, tview.NewPages(), app, false, NewModalInput("Add Task"), NewModalInput("Edit Task"), NewModalInputMode("Add Mode", todoDirModes), nil}
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
		l.lanes[i].SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyF10:
				l.CmdExit()
				return nil
			case tcell.KeyF6:
				l.CmdSelectNote()
				return nil
			case tcell.KeyF5:
				l.CmdArchiveNote()
				return nil
			case tcell.KeyF4:
				l.CmdEditNote()
				return nil
			case tcell.KeyF3:
				l.CmdEditTask()
				return nil
			case tcell.KeyF2:
				fallthrough
			case tcell.KeyInsert:
				l.CmdAddTask()
				return nil
			case tcell.KeyDelete:
				l.pages.ShowPage("delete")
				return nil
			case tcell.KeyTab:
				l.incActive()
				return nil
			case tcell.KeyBacktab:
				l.decActive()
				return nil
			case tcell.KeyUp:
				if l.inselect {
					l.up()
					return nil
				}
			case tcell.KeyDown:
				if l.inselect {
					l.down()
					return nil
				}
			case tcell.KeyLeft:
				if l.inselect {
					l.moveSelectionLeft()
				} else {
					l.decActive()
				}
				return nil
			case tcell.KeyRight:
				if l.inselect {
					l.moveSelectionRight()
				} else {
					l.incActive()
				}
				return nil
			case tcell.KeyF1:
				l.CmdAbout()
				return nil
			}
			switch event.Rune() {
			case 'q':
				l.CmdExit()
			case 'h':
				fallthrough
			case '?':
				l.CmdAbout()
			case 'd':
				l.pages.ShowPage("delete")
			case '+':
				l.CmdAddTask()
				return nil
			case 'a':
				l.CmdArchiveNote()
				return nil
			case 'e':
				l.CmdEditTask()
			case 'n':
				l.CmdEditNote()
			case 'm':
				l.CmdSelectMode()
			}
			return event
		})
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
				l.redraw(l.active, item)
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
				l.redraw(l.active, item)
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
			l.content.AddItem(l.active, item, text, secondary)
			l.redraw(l.active, item)
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
			l.redraw(l.active, item)
		}
		l.pages.HidePage("edit")
		l.setActive()
	})
	l.pages.AddPage("edit", l.edit, false, false)

	return l
}

func (l *Lanes) selected() {
	if l.inselect {
		l.lanes[l.active].SetSelectedBackgroundColor(tcell.ColorWhite)
		l.lanes[l.active].SetSelectedTextColor(tcell.ColorBlack)
	} else {
		l.lanes[l.active].SetSelectedBackgroundColor(tcell.ColorNavy)
		l.lanes[l.active].SetSelectedTextColor(tcell.ColorWhite)
	}
	l.inselect = !l.inselect
	if l.inselect {
		l.bMoveHelp.SetLabel("[blue::-]↔ ↕ [black::b]Use Arrow Keys to Move Task")
	} else {
		l.bMoveHelp.SetLabel("")
	}
}

func (l *Lanes) redraw(lane, active int) {
	l.lanes[lane].Clear()
	for _, item := range l.content.GetLaneItems(lane) {
		l.lanes[lane].AddItem(item.Title, item.Secondary, 0, nil)
	}
	num := l.lanes[lane].GetItemCount()
	if num > 0 {
		l.lanes[lane].SetCurrentItem(normPos(active, num))
	}
	l.lanes[lane].SetTitle(l.content.GetLaneTitle(lane))
}

func (l *Lanes) RedrawAllLanes() {
	for idx := range l.lanes {
		currentPos := l.lanes[idx].GetCurrentItem()
		newPos := normPos(currentPos, l.lanes[idx].GetItemCount())
		l.redraw(idx, newPos)
	}
}

func (l *Lanes) up() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newPos := normPos(currentPos-1, l.lanes[l.active].GetItemCount())
	l.content.MoveItem(l.active, currentPos, l.active, newPos)
	l.redraw(l.active, newPos)
}

func (l *Lanes) down() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newPos := normPos(currentPos+1, l.lanes[l.active].GetItemCount())
	l.content.MoveItem(l.active, currentPos, l.active, newPos)
	l.redraw(l.active, newPos)
}

func (l *Lanes) moveSelectionLeft() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newLane := normPos(l.active-1, len(l.lanes))
	newPos := l.lanes[newLane].GetCurrentItem()
	l.content.MoveItem(l.active, currentPos, newLane, newPos)
	l.redraw(l.active, currentPos)
	l.redraw(newLane, newPos)
	l.selected()
	l.decActive()
	l.selected()
}

func (l *Lanes) moveSelectionRight() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newLane := normPos(l.active+1, len(l.lanes))
	newPos := l.lanes[newLane].GetCurrentItem()
	l.content.MoveItem(l.active, currentPos, newLane, newPos)
	l.redraw(l.active, currentPos)
	l.redraw(newLane, newPos)
	l.selected()
	l.incActive()
	l.selected()

}

func (l *Lanes) decActive() {
	l.lanes[l.active].SetSelectedStyle(tcell.StyleDefault)
	l.active--
	l.setActive()
	l.lanes[l.active].SetSelectedBackgroundColor(tcell.ColorLightBlue)
	l.lanes[l.active].SetSelectedTextColor(tcell.ColorBlack)
}

func (l *Lanes) incActive() {
	l.lanes[l.active].SetSelectedStyle(tcell.StyleDefault)
	l.active++
	l.setActive()
	l.lanes[l.active].SetSelectedBackgroundColor(tcell.ColorLightBlue)
	l.lanes[l.active].SetSelectedTextColor(tcell.ColorBlack)
}

func normPos(pos, length int) int {
	for pos < 0 {
		pos += length
	}
	if length > 0 {
		pos %= length
	}
	return pos
}

func (l *Lanes) setActive() {
	l.active = normPos(l.active, len(l.lanes))
	l.app.SetFocus(l.lanes[l.active])
}

func (l *Lanes) saveActive() {
	l.setActive()
	l.lastActive = l.active
	l.lastActiveSaved = true
}

func (l *Lanes) currentItem() *Item {
	pos := l.lanes[l.active].GetCurrentItem()
	content := l.content.GetLaneItems(l.active)
	if pos < 0 || pos >= len(content) {
		return nil
	}
	return &content[pos]
}

func (l *Lanes) editNote() {
	item := l.currentItem()
	if item != nil {
		tmp, err := os.CreateTemp("", "todo_temp_note_")
		if err == nil {
			name := tmp.Name()
			defer os.Remove(name)
			tmp.Write([]byte(item.Note))
			tmp.Close()
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("notepad", name)
				err = cmd.Start()
				if err != nil {
					log.Fatal(err)
				}
				err = cmd.Wait()
				if err != nil {
					log.Fatal(err)
				}
				if err == nil {
					note_raw, err := os.ReadFile(name)
					if err == nil {
						item.Note = string(note_raw)
					}
				}
			} else {
				cmd = exec.Command("vim", name)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err == nil {
					note_raw, err := os.ReadFile(name)
					if err == nil {
						item.Note = string(note_raw)
					}
				}
			}
		}
	}
}

func (l *Lanes) GetUi() *tview.Pages {
	return l.pages
}
