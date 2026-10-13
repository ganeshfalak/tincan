package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "ctrl+c" {
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() tea.View {
	v := tea.NewView("tincan\n\n ctrl+c quits\n")
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

	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
