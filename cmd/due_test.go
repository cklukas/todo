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
