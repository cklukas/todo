package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gdamore/tcell/v2"
)

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

func (l *Lanes) CmdEditNote() {
	if runtime.GOOS == "windows" {
		l.pages.ShowPage("wait")
		l.app.ForceDraw()
		l.editNote()
		l.pages.HidePage("wait")
	} else {
		if !l.app.Suspend(l.editNote) {
			l.app.Stop()
			log.Fatal("internal suspend error")
		}
	}
	l.content.Save()
}

func (l *Lanes) CmdArchiveNote() {
	if len(l.lanes) > 0 {
		ll := (*l.lanes[l.active]).GetItemCount()
		if ll > 0 {
			l.saveActive()
			l.pages.ShowPage("archive")
		}
	}
}

func (l *Lanes) CmdSelectNote() {
	l.setActive()
	l.selected()
}

func (l *Lanes) up() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newPos := normPos(currentPos-1, l.lanes[l.active].GetItemCount())
	l.content.MoveItem(l.active, currentPos, l.active, newPos)
	l.redrawLane(l.active, newPos)
}

func (l *Lanes) down() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newPos := normPos(currentPos+1, l.lanes[l.active].GetItemCount())
	l.content.MoveItem(l.active, currentPos, l.active, newPos)
	l.redrawLane(l.active, newPos)
}

func (l *Lanes) moveSelectionLeft() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newLane := normPos(l.active-1, len(l.lanes))
	newPos := l.lanes[newLane].GetCurrentItem()
	l.content.MoveItem(l.active, currentPos, newLane, newPos)
	l.redrawLane(l.active, currentPos)
	l.redrawLane(newLane, newPos)
	l.selected()
	l.decActive()
	l.selected()
}

func (l *Lanes) moveSelectionRight() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newLane := normPos(l.active+1, len(l.lanes))
	newPos := l.lanes[newLane].GetCurrentItem()
	l.content.MoveItem(l.active, currentPos, newLane, newPos)
	l.redrawLane(l.active, currentPos)
	l.redrawLane(newLane, newPos)
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

func (l *Lanes) selected() {
	selectDisable := true
	if len(l.lanes) > 0 {
		ll := (*l.lanes[l.active]).GetItemCount()
		if ll > 0 {
			selectDisable = false
		}
	}
	if l.inselect || selectDisable {
		l.lanes[l.active].SetSelectedBackgroundColor(tcell.ColorWhite)
		l.lanes[l.active].SetSelectedTextColor(tcell.ColorBlack)
	} else {
		l.lanes[l.active].SetSelectedBackgroundColor(tcell.ColorNavy)
		l.lanes[l.active].SetSelectedTextColor(tcell.ColorWhite)
	}
	l.inselect = !l.inselect
	if selectDisable {
		l.inselect = false
	}
	if l.inselect {
		l.bMoveHelp.SetLabel("[blue::-]↔ ↕ [black::b]Use Arrow Keys to Move Task")
	} else {
		l.bMoveHelp.SetLabel("")
	}
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
					l.app.Stop()
					log.Fatal(err)
				}
				err = cmd.Wait()
				if err != nil {
					l.app.Stop()
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
					if err != nil {
						l.app.Stop()
						log.Fatal(err)
					}
					item.Note = string(note_raw)
				}
			}
		}
	}
}
