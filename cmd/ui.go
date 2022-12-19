package cmd

import (
	"fmt"

	"github.com/rivo/tview"
)

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

	l.lanes[laneIndex].Clear()
	for _, item := range l.content.GetLaneItems(laneIndex) {
		l.lanes[laneIndex].AddItem(item.Title, item.Secondary, 0, nil)
	}

	num := l.lanes[laneIndex].GetItemCount()
	if num > 0 {
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
