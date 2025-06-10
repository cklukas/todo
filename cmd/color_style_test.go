package cmd

import (
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/cklukas/todo/internal/ui"
	"github.com/rivo/tview"
)

func getDropDownTexts(dd *tview.DropDown) []string {
	v := reflect.ValueOf(dd).Elem().FieldByName("options")
	res := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		res[i] = v.Index(i).Elem().FieldByName("Text").String()
	}
	return res
}

func getDropDownTextsFromCI(ci *ui.ColorInput) []string {
	v := reflect.ValueOf(ci).Elem().FieldByName("dropdown")
	dd := (*tview.DropDown)(unsafe.Pointer(v.Pointer()))
	return getDropDownTexts(dd)
}

func TestColorModalLabelsHaveBackground(t *testing.T) {
	dlg := ui.NewColorModal("Color", "")
	dd := dlg.GetFormItem(0).(*tview.DropDown)
	texts := getDropDownTexts(dd)
	trimmed := strings.TrimRight(texts[1], " ")
	if trimmed != "[black:black]black" || len(texts[1]) == len(trimmed) {
		t.Fatalf("color modal text for black wrong: %q", texts[1])
	}
}

func TestColorInputLabelsUseLaneBackground(t *testing.T) {
	ci := ui.NewColorInput("Title:", "yellow", []string{"default", "red"})
	texts := getDropDownTextsFromCI(ci)
	trimmed := strings.TrimRight(texts[1], " ")
	if trimmed != "[red:yellow]red" || len(texts[1]) == len(trimmed) {
		t.Fatalf("color input text for red wrong: %q", texts[1])
	}
}
