package promptcode

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textInput    textinput.Model
	message      string
	errorMessage string
	quitting     bool
	codeValid    bool
	validate     func(string) error
}

var (
	fixedErrorSize = 2 // Define a fixed height for the error message space

	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Height(fixedErrorSize)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
)

type tickMsg struct{}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tickCmd())
}

// tickCmd creates a command that sends a `tickMsg` every 300ms.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*300, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Update handles the tea.Msgs (inputs or system messages) and updates the model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if err := m.validate(m.textInput.Value()); err != nil {
			m.errorMessage = err.Error()
		} else {
			m.errorMessage = ""
			if len(m.textInput.Value()) > 0 && err == nil {
				m.message = successStyle.Render("Code accepted! Session ended.")
				m.codeValid = true
				return m, tea.Quit
			}
		}
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.textInput.Value() != "" {
				m.textInput.SetValue("")
				return m, tickCmd()
			} else {
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the Bubble Tea application's UI.
func (m model) View() string {
	if m.quitting {
		return ""
	}

	return fmt.Sprintf(
		"%s\n%s",
		errorStyle.Render(m.errorMessage),
		m.textInput.View(),
	)
}

func Prompt(url string, validate func(string) error) error {
	fmt.Printf(
		"Please navigate to the following URL in your browser: %s\n"+
			"Then, paste the base64 code below. Press ctrl+c to quit.",
		url,
	)

	ti := textinput.New()
	ti.Placeholder = "Paste your base64 code here"
	ti.Focus()

	p := tea.NewProgram(model{
		textInput: ti,
		validate:  validate,
	})
	_, err := p.Run()
	return err
}
