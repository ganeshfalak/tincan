package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
}

func newModel() *model {
	input := textinput.New()
	input.Placeholder = "write a message"
	return &model{
		input: input,
		peers: []string{"alice", "bob", "cara"},
	}
}

func (m *model) Init() tea.Cmd {
	return m.input.Focus()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c":
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
