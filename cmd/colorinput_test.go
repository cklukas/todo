package cmd

import "testing"

func TestParsePrefix(t *testing.T) {
	color, text := parsePrefix("[blue]hello")
	if color != "blue" || text != "hello" {
		t.Fatalf("expected blue/hello got %s/%s", color, text)
	}
	color, text = parsePrefix("hello")
	if color != "" || text != "hello" {
		t.Fatalf("expected empty/hello got %s/%s", color, text)
	}
}

func TestApplyPrefix(t *testing.T) {
	if out := applyPrefix("hello", "blue"); out != "[blue]hello" {
		t.Fatalf("applyPrefix failed: %s", out)
	}
	if out := applyPrefix("[red]hello", ""); out != "hello" {
		t.Fatalf("applyPrefix remove failed: %s", out)
	}
}
