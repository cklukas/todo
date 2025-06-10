package ui

import (
	"os"
	"testing"
	"time"
)

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

func TestIsoTimeToLocal(t *testing.T) {
	iso := "2025-03-23T14:11:51Z"

	os.Unsetenv("LC_TIME")
	os.Unsetenv("AppleLocale")
	t1, _ := time.Parse(time.RFC3339, iso)
	want := t1.Local().Format(dueLayout() + " 15:04")
	if got := isoTimeToLocal(iso); got != want {
		t.Fatalf("isoTimeToLocal non-US got %q want %q", got, want)
	}

	os.Setenv("LC_TIME", "en_US")
	os.Setenv("AppleLocale", "en_US")
	t2, _ := time.Parse(time.RFC3339, iso)
	want = t2.Local().Format(dueLayout() + " 15:04")
	if got := isoTimeToLocal(iso); got != want {
		t.Fatalf("isoTimeToLocal US got %q want %q", got, want)
	}
	os.Unsetenv("LC_TIME")
	os.Unsetenv("AppleLocale")
}
