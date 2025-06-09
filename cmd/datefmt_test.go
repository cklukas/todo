package cmd

import "testing"

func TestLocaleUSApple(t *testing.T) {
	cases := map[string]bool{
		"en_US":           true,
		"en_US@rg=uszzzz": true,
		"en_US@rg=dezzzz": false,
		"de_DE":           false,
		"en_GB":           false,
	}
	for in, want := range cases {
		if got := localeUSApple(in); got != want {
			t.Fatalf("localeUSApple(%q)=%v want %v", in, got, want)
		}
	}
}
