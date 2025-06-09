package cmd

import "testing"

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
	c.AddItem(0, 0, "task", "")
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
	c.AddItem(0, 0, "task", "")
	title := c.GetLaneTitle(0)
	if title != " To Do (1) " {
		t.Fatalf("unexpected title: %s", title)
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
		if got := normPos(tc.pos, tc.length); got != tc.exp {
			t.Errorf("normPos(%d,%d)=%d expected %d", tc.pos, tc.length, got, tc.exp)
		}
	}
}

func TestSplit(t *testing.T) {
	words, err := Split("one 'two three' four")
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
