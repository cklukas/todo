package ui

func (l *Lanes) CmdSortDialog() {
	l.saveActive()
	laneTitle := l.GetActiveLaneName()
	current := l.content.SortModes[l.active]
	dlg := NewSortModal(" Sort Tasks ", laneTitle, current)
	dlg.SetDoneFunc(func(mode string, ok bool) {
		l.pages.HidePage("sort")
		l.setActive()
		if ok {
			l.content.SetLaneSort(l.active, mode)
			l.redrawLane(l.active, l.lanes[l.active].GetCurrentItem())
			l.content.Save()
		}
	})
	l.pages.AddPage("sort", modal(dlg, 0, 0), false, true)
	l.showDialog("sort", dlg)
}
