package cmd

import (
	"reflect"
	"testing"

	"github.com/cklukas/todo/internal/ui"
)

func TestDOSColorLists(t *testing.T) {
	wantInput := []string{"default", "black", "blue", "green", "aqua", "red", "purple", "brown", "silver", "gray", "darkblue", "darkgreen", "darkcyan", "darkred", "fuchsia", "yellow"}
	mi := ui.NewModalInput("Add Task")
	v := reflect.ValueOf(mi).Elem().FieldByName("colors")
	gotInput := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		gotInput[i] = v.Index(i).String()
	}
	if !reflect.DeepEqual(gotInput, wantInput) {
		t.Fatalf("modal input colors %v, want %v", gotInput, wantInput)
	}

	wantModal := []string{"", "black", "blue", "green", "aqua", "red", "purple", "brown", "silver", "gray", "darkblue", "darkgreen", "darkcyan", "darkred", "fuchsia", "yellow"}
	cm := ui.NewColorModal("Color", "")
	v2 := reflect.ValueOf(cm).Elem().FieldByName("options")
	gotModal := make([]string, v2.Len())
	for i := 0; i < v2.Len(); i++ {
		gotModal[i] = v2.Index(i).String()
	}
	if !reflect.DeepEqual(gotModal, wantModal) {
		t.Fatalf("color modal options %v, want %v", gotModal, wantModal)
	}
}
