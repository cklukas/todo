package model

import (
	"sort"
	"time"
)

const (
	SortNone     = ""
	SortColor    = "color"
	SortDue      = "due"
	SortCreated  = "created"
	SortModified = "modified"
	SortPriority = "priority"
)

func PriorityMark(p int) string {
	switch p {
	case 1:
		return "↑"
	case 3:
		return "↓"
	case 4:
		return "⌛"
	default:
		return ""
	}
}

func sortItems(items []Item, mode string) {
	switch mode {
	case SortColor:
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Color < items[j].Color
		})
	case SortDue:
		sort.SliceStable(items, func(i, j int) bool {
			ti := parseDue(items[i].Due)
			tj := parseDue(items[j].Due)
			return ti.Before(tj)
		})
	case SortCreated:
		sort.SliceStable(items, func(i, j int) bool {
			ti := parseTime(items[i].Created)
			tj := parseTime(items[j].Created)
			return ti.Before(tj)
		})
	case SortModified:
		sort.SliceStable(items, func(i, j int) bool {
			ti := parseTime(items[i].LastUpdate)
			if ti.IsZero() {
				ti = parseTime(items[i].Created)
			}
			tj := parseTime(items[j].LastUpdate)
			if tj.IsZero() {
				tj = parseTime(items[j].Created)
			}
			return ti.Before(tj)
		})
	case SortPriority:
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Priority < items[j].Priority
		})
	}
}

func parseDue(d string) time.Time {
	if t, err := time.Parse("2006-01-02", d); err == nil {
		return t
	}
	return time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)
}

func parseTime(tstr string) time.Time {
	if t, err := time.Parse(time.RFC3339, tstr); err == nil {
		return t
	}
	return time.Time{}
}
