package ui

import (
	"os"
	"testing"
	"time"
)

func TestDueSuffix(t *testing.T) {
	now := time.Date(2025, 6, 10, 10, 0, 0, 0, time.UTC)
	if s := dueSuffix("2025-06-10", now); s != "[due!]" {
		t.Fatalf("expected [due!] got %s", s)
	}
	if s := dueSuffix("2025-06-11", now); s != "[tomorrow]" {
		t.Fatalf("expected [tomorrow] got %s", s)
	}
	if s := dueSuffix("2025-06-12", now); s != "" {
		t.Fatalf("expected empty suffix got %s", s)
	}
}

func TestFormatDueInput(t *testing.T) {
	os.Unsetenv("LC_TIME")
	cases := map[string]string{
		"1":          "1",
		"12":         "12.",
		"120":        "12.0",
		"1203":       "12.03.",
		"120320":     "12.03.20",
		"12032023":   "12.03.2023",
		"12.03.2023": "12.03.2023",
	}
	for in, expect := range cases {
		if out := formatDueInput(in); out != expect {
			t.Fatalf("formatDueInput(%q) = %q, want %q", in, out, expect)
		}
	}
}

func TestFormatDueInputUS(t *testing.T) {
	os.Setenv("LC_TIME", "en_US")
	defer os.Unsetenv("LC_TIME")

	cases := map[string]string{
		"1":          "1",
		"12":         "12/",
		"120":        "12/0",
		"1203":       "12/03/",
		"120320":     "12/03/20",
		"12032023":   "12/03/2023",
		"12/03/2023": "12/03/2023",
	}
	for in, expect := range cases {
		if out := formatDueInput(in); out != expect {
			t.Fatalf("formatDueInput(%q) = %q, want %q", in, out, expect)
		}
	}
}

func TestRemoveLastDueDigit(t *testing.T) {
	os.Unsetenv("LC_TIME")
	cases := map[string]string{
		"":           "",
		"1":          "",
		"12.":        "1",
		"12.03.2023": "12.03.202",
	}
	for in, expect := range cases {
		if out := removeLastDueDigit(in); out != expect {
			t.Fatalf("removeLastDueDigit(%q) = %q, want %q", in, out, expect)
		}
	}

	os.Setenv("LC_TIME", "en_US")
	defer os.Unsetenv("LC_TIME")
	if out := removeLastDueDigit("12/03/2023"); out != "12/03/202" {
		t.Fatalf("removeLastDueDigit US failed: %s", out)
	}
}
