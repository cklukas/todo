package model

import "testing"

func TestPriorityMark(t *testing.T) {
	if m := PriorityMark(1); m == "" {
		t.Fatalf("high priority mark empty")
	}
	if m := PriorityMark(2); m != "" {
		t.Fatalf("default priority should be empty")
	}
}

func TestSortByDue(t *testing.T) {
	items := []Item{
		{Title: "t1", Due: "2025-06-10"},
		{Title: "t2", Due: ""},
		{Title: "t3", Due: "2025-01-01"},
	}
	sortItems(items, SortDue)
	if items[0].Title != "t3" || items[1].Title != "t1" || items[2].Title != "t2" {
		t.Fatalf("due sort failed: %#v", items)
	}
}
