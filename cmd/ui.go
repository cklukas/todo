package cmd

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type dialogWithFrame interface {
	tview.Primitive
	GetFrame() *tview.Frame
}

type Lanes struct {
	nextMode        string
	nextLaneFocus   int
	todoDirModes    string
	mode            string
	content         *ToDoContent
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

	dialogActive     bool
	activeDialog     dialogWithFrame
	origInputCapture func(event *tcell.EventKey) *tcell.EventKey
	origMouseCapture func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction)
}

func (l *Lanes) CmdAbout() {
	l.saveActive()
	l.pages.ShowPage("help")
}

func (l *Lanes) CmdExit() {
	l.setActive()
	l.pages.ShowPage("quit")
}

func (l *Lanes) Focus() {
	l.setActive()
	l.selected()
}

func (l *Lanes) redrawLane(laneIndex, active int) error {
	if laneIndex >= len(l.lanes) {
		return fmt.Errorf("invalid index '%v', visible lines count is only '%v'", laneIndex, len(l.lanes))
	}

	prev := ""
	if laneIndex == l.active {
		if it := l.currentItem(); it != nil {
			prev = it.Guid
		}
	}
	l.content.SortLane(laneIndex)
	l.lanes[laneIndex].Clear()
	now := time.Now()
	for _, item := range l.content.GetLaneItems(laneIndex) {
		title := item.Title
		if item.Color != "" {
			title = "[" + item.Color + "]" + title
		}
		if suffix := dueSuffix(item.Due, now); suffix != "" {
			title += " " + suffix
		}
		secondary := item.Secondary
		if mark := priorityMark(item.Priority); mark != "" {
			if len(secondary) > 0 {
				secondary += " "
			}
			secondary += mark
		}
		l.lanes[laneIndex].AddItem(title, secondary, 0, nil)
	}

	num := l.lanes[laneIndex].GetItemCount()
	if num > 0 {
		if prev != "" {
			for i, item := range l.content.GetLaneItems(laneIndex) {
				if item.Guid == prev {
					active = i
					break
				}
			}
		}
		l.lanes[laneIndex].SetCurrentItem(normPos(active, num))
	}

	l.lanes[laneIndex].SetTitle(l.content.GetLaneTitle(laneIndex))
	return nil
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

func dueSuffix(due string, now time.Time) string {
	if len(due) == 0 {
		return ""
	}
	d, err := time.Parse("2006-01-02", due)
	if err != nil {
		return ""
	}
	d = d.In(now.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	days := int(d.Sub(today).Hours() / 24)
	switch days {
	case 0:
		return "[due!]"
	case 1:
		return "[tomorrow]"
	default:
		return ""
	}
}

func (l *Lanes) setActive() {
	l.active = normPos(l.active, len(l.lanes))
	l.app.SetFocus(l.lanes[l.active])
}

func (l *Lanes) setActiveIndex(index int) {
	l.active = normPos(index, len(l.lanes))
	l.app.SetFocus(l.lanes[l.active])
}

func (l *Lanes) saveActive() int {
	l.setActive()
	l.lastActive = l.active
	l.lastActiveSaved = true
	return l.lastActive
}

func (l *Lanes) currentItem() *Item {
	pos := l.lanes[l.active].GetCurrentItem()
	content := l.content.GetLaneItems(l.active)
	if pos < 0 || pos >= len(content) {
		return nil
	}
	return &content[pos]
}

func (l *Lanes) GetUi() *tview.Pages {
	return l.pages
}

func (l *Lanes) appInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if l.dialogActive {
		return event
	}
	if l.origInputCapture != nil {
		return l.origInputCapture(event)
	}
	return event
}

func (l *Lanes) appMouseCapture(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	if l.dialogActive && l.activeDialog != nil {
		x, y := event.Position()
		if !l.activeDialog.GetFrame().InRect(x, y) {
			return nil, action
		}
	}
	if l.origMouseCapture != nil {
		return l.origMouseCapture(event, action)
	}
	return event, action
}

func (l *Lanes) showDialog(name string, modal dialogWithFrame) {
	l.dialogActive = true
	l.activeDialog = modal
	l.pages.ShowPage(name)
	l.app.SetFocus(modal)
}

func (l *Lanes) hideDialog(name string) {
	l.dialogActive = false
	l.activeDialog = nil
	l.pages.HidePage(name)
	l.setActive()
}
