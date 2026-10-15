package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"tincan/internal/transport"
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
	conn     net.Conn
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
	if m.conn != nil {
		conn := m.conn
		return func() tea.Msg {
			buf, err := transport.ReadFrame(conn)
			if err != nil {
				return nil
			}
			return lineMsg(string(buf))
		}
	}

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
			if m.conn != nil {
				m.conn.Close()
			}
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

	listen := flag.String("listen", "", "show received frames in the screen")
	dial := flag.String("dial", "", "send stdin lines to this localhost address")
	flag.Parse()

	if *debug && *config != "" {
		fmt.Fprintf(os.Stderr, "config: %s\n", *config)
	}

	if *listen != "" && *dial != "" {
		fmt.Fprintln(os.Stderr, "use --listen or --dial, not both")
		os.Exit(1)
	}
	if *listen != "" {
		conn, err := acceptOne(*listen)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer conn.Close()
		m := newModel()
		m.conn = conn
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if *dial != "" {
		if err := oneDial(*dial); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func requireLocalhost(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if host != "127.0.0.1" {
		return fmt.Errorf("localhost only, got %s", host)
	}
	return nil
}

func acceptOne(addr string) (net.Conn, error) {
	if err := requireLocalhost(addr); err != nil {
		return nil, err
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer ln.Close()
	return ln.Accept()
}

func oneDial(addr string) error {
	if err := requireLocalhost(addr); err != nil {
		return err
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if err := transport.WriteFrame(conn, []byte(line)); err != nil {
			return err
		}
	}
	return sc.Err()
}
