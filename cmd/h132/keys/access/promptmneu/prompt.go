package promptmneu

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/IPA-CyberLab/h132/cmd/h132/common"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
)

type model struct {
	ti   textinput.Model
	done bool
}

var (
	fixedErrorSize = 2 // Define a fixed height for the error message space

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Height(fixedErrorSize)

	ErrEmpty = errors.New("input is empty")
)

// Check `v` is a valid bip39 mneumonic phrase.
// This func is meant to provide user friendly error message compared to
// bare `bip39.EntropyFromMnemonic`.
func friendlyVerify(v string) error {
	if v == "" {
		return ErrEmpty
	}
	bip39.SetWordList(wordlists.Japanese)

	words := strings.Fields(v)
	for _, w := range words {
		if !slices.Contains(bip39.GetWordList(), w) {
			return fmt.Errorf("Unknown word: %q", w)
		}
	}

	if len(words) != 24 {
		return fmt.Errorf("Not enough words: %d / need 24", len(words))
	}

	if _, err := bip39.EntropyFromMnemonic(v); err != nil {
		return err
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles the tea.Msgs (inputs or system messages) and updates the model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if v := m.mneumonic(); v != "" {
				if err := friendlyVerify(v); err == nil {
					m.done = true
					return m, tea.Quit
				}
			}

		case "ctrl+c", "esc":
			if m.ti.Value() != "" {
				m.ti.SetValue("")
				return m, nil
			} else {
				m.done = true
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)

	return m, cmd
}

func (m model) mneumonic() string {
	s := m.ti.Value()
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s,
		"　", // Fullwidth space
		" ")
	return s
}

// View renders the Bubble Tea application's UI.
func (m model) View() string {
	if m.done {
		return "Mneumonic accepted!"
	}

	errorMessage := ""
	if v := m.mneumonic(); v != "" {
		if err := friendlyVerify(v); err != nil {
			errorMessage = err.Error()
		}
	}

	return fmt.Sprintf(
		"%s\n%s",
		errorStyle.Render(errorMessage),
		m.ti.View(),
	)
}

func Prompt(name, hint string) (string, error) {
	fmt.Printf("To access the emergency key %q, please enter its mneumonic.\n", name)
	fmt.Printf("Hint: %s\n", hint)

	ti := textinput.New()
	ti.Placeholder = "Please enter the mneumonic here"
	ti.Focus()

	p := tea.NewProgram(model{
		ti:   ti,
		done: false,
	})
	m, err := p.Run()
	if err != nil {
		return "", err
	}
	m2 := m.(model)

	v := m2.mneumonic()
	if v == "" {
		return "", common.ErrAbort{}
	}

	bip39.SetWordList(wordlists.Japanese)
	if _, err := bip39.EntropyFromMnemonic(v); err != nil {
		return "", fmt.Errorf("failed to decode mneumonic: %w", err)
	}

	return v, nil
}
