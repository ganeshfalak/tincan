package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// model is the screen's memory. Bubble Tea keeps this value and passes it
// back on every key.
type model struct {
	input    textinput.Model
	lines    []string
	peers    []string
	selected int
	ctx      context.Context
	cancel   context.CancelFunc
}

func newModel() *model {
	input := textinput.New()
	input.Placeholder = "write a message"
	ctx, cancel := context.WithCancel(context.Background())
	return &model{
		input:  input,
		peers:  []string{"alice", "bob", "cara"},
		ctx:    ctx,
		cancel: cancel,
	}
}

type lineMsg string

func (m *model) waitLine() tea.Cmd {
	ctx := m.ctx
	return func() tea.Msg {
		select {
		case <-time.After(2 * time.Second):
			return lineMsg("bob: ping")
		case <-ctx.Done():
			return nil
		}
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.input.Focus(), m.waitLine())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg == nil {
		return m, nil
	}
	if line, ok := msg.(lineMsg); ok {
		m.lines = append(m.lines, string(line))
		return m, m.waitLine()
	}

	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c":
			m.cancel()
			return m, tea.Quit
		case "tab":
			m.selected = (m.selected + 1) % len(m.peers)
			return m, nil
		case "shift+tab":
			m.selected = (m.selected - 1 + len(m.peers)) % len(m.peers)
			return m, nil
		case "enter":
			line := strings.TrimSpace(m.input.Value())
			if line != "" {
				m.lines = append(m.lines, m.peers[m.selected]+": "+line)
			}
			m.input.SetValue("")
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *model) View() tea.View {
	body := "tincan\n\n"
	for i, name := range m.peers {
		if i == m.selected {
			body += "> " + name + "\n"
		} else {
			body += " " + name + "\n"
		}
	}

	body += "\n"

	if len(m.lines) == 0 {
		body += "(no messages)\n"
	} else {
		for _, line := range m.lines {
			body += line + "\n"
		}
	}

	body += "\n" + m.input.View() + "\n"

	body += "tab next    shift+tab back    enter send    ctrl+c quit\n"

	v := tea.NewView(body)
	v.AltScreen = true
	return v
}

func main() {
	config := flag.String("config", "", "config directory")
	debug := flag.Bool("debug", false, "print the --config value on stderr")
	flag.Parse()

	if *debug && *config != "" {
		fmt.Fprintf(os.Stderr, "config: %s\n", *config)
	}

	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
