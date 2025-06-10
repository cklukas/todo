package ui

import "github.com/gdamore/tcell/v2"

func (l *Lanes) HotKeyHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyF10:
		l.CmdExit()
		return nil
		// case tcell.KeyF9:
		// 	l.CmdSelectMode()
		// 	return nil
	case tcell.KeyF7:
		l.CmdLanesCmds()
		return nil
	case tcell.KeyF8:
		l.CmdSortDialog()
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
		l.CmdSelectModeDialog()
	}
	return event
}
