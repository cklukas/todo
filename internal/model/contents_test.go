package model

import (
	"testing"

	"github.com/cklukas/todo/internal/util"
)

func TestInsertNewLane(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	idx := c.InsertNewLane(false, "Review", 0)
	if idx != 1 {
		t.Fatalf("expected index 1 got %d", idx)
	}
	if len(c.Titles) != 4 {
		t.Fatalf("expected 4 lanes got %d", len(c.Titles))
	}
	if c.Titles[1] != "Review" {
		t.Fatalf("lane title not inserted: %v", c.Titles)
	}
}

func TestMoveItem(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "task", "", 2, "", "")
	c.MoveItem(0, 0, 1, 0)
	if len(c.Items[0]) != 0 {
		t.Fatalf("item not removed from source")
	}
	if len(c.Items[1]) != 1 || c.Items[1][0].Title != "task" {
		t.Fatalf("item not moved correctly")
	}
}

func TestRemoveLane(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	idx := c.InsertNewLane(false, "X", 1)
	c.RemoveLane(idx)
	if len(c.Titles) != 3 {
		t.Fatalf("expected 3 lanes got %d", len(c.Titles))
	}
	if c.Titles[1] != "Doing" {
		t.Fatalf("expected lane 'Doing' at index 1 got %s", c.Titles[1])
	}
}

func TestGetLaneTitle(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "task", "", 2, "", "")
	title := c.GetLaneTitle(0)
	if title != " To Do (1) " {
		t.Fatalf("unexpected title: %s", title)
	}
}

func TestAddItemPriority(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "p task", "", 3, "", "")
	if c.Items[0][0].Priority != 3 {
		t.Fatalf("expected priority 3 got %d", c.Items[0][0].Priority)
	}
}

func TestAddItemDue(t *testing.T) {
	c := &ToDoContent{}
	c.InitializeNew()
	c.AddItem(0, 0, "due task", "", 2, "2025-06-10", "")
	if c.Items[0][0].Due != "2025-06-10" {
		t.Fatalf("expected due 2025-06-10 got %s", c.Items[0][0].Due)
	}
}

func TestNormPos(t *testing.T) {
	tests := []struct{ pos, length, exp int }{
		{3, 5, 3},
		{5, 5, 0},
		{-1, 5, 4},
		{6, 5, 1},
	}
	for _, tc := range tests {
		if got := util.NormPos(tc.pos, tc.length); got != tc.exp {
			t.Errorf("normPos(%d,%d)=%d expected %d", tc.pos, tc.length, got, tc.exp)
		}
	}
}

func TestSplit(t *testing.T) {
	words, err := util.Split("one 'two three' four")
	if err != nil {
		t.Fatalf("Split returned error: %v", err)
	}
	exp := []string{"one", "two three", "four"}
	if len(words) != len(exp) {
		t.Fatalf("expected %d words got %d", len(exp), len(words))
	}
	for i, w := range exp {
		if words[i] != w {
			t.Fatalf("expected word %q got %q", w, words[i])
		}
	}
}

func TestNormalizeSortModes(t *testing.T) {
	c := &ToDoContent{Titles: []string{"A", "B"}, Items: make([][]Item, 2)}
	c.normalize()
	if len(c.SortModes) != 2 {
		t.Fatalf("expected 2 sort modes got %d", len(c.SortModes))
	}
	for i, m := range c.SortModes {
		if m != "" {
			t.Fatalf("sort mode %d should be empty", i)
		}
	}
}

func TestNormalizeLaneColors(t *testing.T) {
	c := &ToDoContent{Titles: []string{"A", "B"}, Items: make([][]Item, 2)}
	c.normalize()
	if len(c.LaneColors) != 2 {
		t.Fatalf("expected 2 lane colors got %d", len(c.LaneColors))
	}
	for i, col := range c.LaneColors {
		if col != "" {
			t.Fatalf("lane color %d should be empty", i)
		}
	}
}
