package promptrm

import (
	"fmt"

	"github.com/IPA-CyberLab/h132/cmd/h132/common"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	warningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF0000")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF4500")).
			Padding(1, 2)

	nameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF00"))
)

type model struct {
	// Name of the object to confirm deletion of - we ask the user to type this
	name string

	// Type of the object to confirm deletion of
	typ string

	ti   textinput.Model
	done bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update function to handle interactions
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			v := m.ti.Value()
			if v == "" {
				// Empty input -> exit the prompt
				return m, tea.Quit
			} else if v == m.name {
				// If the name matches the input, exit the prompt
				m.done = true
				return m, tea.Quit
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	// Pass the msg to the text input component
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

// View function to render the UI
func (m model) View() string {
	if m.done {
		return ""
	}

	warningMessage := warningStyle.Render(fmt.Sprintf(
		"WARNING: You are about to delete the %s '%s'.\nThis operation is irreversible!",
		m.typ, m.name,
	))

	return fmt.Sprintf(
		"%s\n\nTo confirm deletion, please type the %s name '%s':\n\n%s\n",
		warningMessage,
		m.typ,
		nameStyle.Render(m.name),
		m.ti.View(),
	)
}

func Prompt(name, typ string) error {
	ti := textinput.New()
	ti.Placeholder = fmt.Sprintf("Type the name of the %s to be deleted", typ)
	ti.Focus()

	p := tea.NewProgram(model{
		name: name,
		typ:  typ,
		ti:   ti,
		done: false,
	})
	m, err := p.Run()
	if err != nil {
		return err
	}
	m2 := m.(model)
	if !m2.done {
		return common.ErrAbort{}
	}
	return nil
}
