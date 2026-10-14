package main

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestEnterSends(t *testing.T) {
	m := newModel()
	m.input.SetValue("hello")
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if len(m.lines) != 1 || m.lines[0] != "alice: hello" {
		t.Fatalf("lines = %#v", m.lines)
	}
	if m.input.Value() != "" {
		t.Fatalf("input = %q", m.input.Value())
	}
}
