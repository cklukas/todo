package cmd

import "testing"
import "time"

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
	cases := map[string]string{
		"1":          "1",
		"12":         "12.",
		"120":        "12.0",
		"1203":       "12.03",
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
