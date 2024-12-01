package promptnewemergency

import (
	"fmt"
	"strings"

	"github.com/IPA-CyberLab/h132/cmd/h132/common"
	"github.com/IPA-CyberLab/h132/keys/emergency"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type state int

const (
	stateQuerying state = iota
	stateAccepted
	stateAborted
)

type newEmergencyKeyFunc func() *emergency.EmergencyBackupKey

type model struct {
	newEmergencyKeyFunc
	key *emergency.EmergencyBackupKey
	state
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.state != stateQuerying {
		return m, tea.Quit // shouldn't reach here in the first place?
	}

	switch msg.(type) {
	case tea.KeyMsg:
		switch msg.(tea.KeyMsg).String() {
		case "y", "Y":
			m.state = stateAccepted
			return m, tea.Quit

		case "n", "N":
			// Generate a new key
			m.key = m.newEmergencyKeyFunc()
			m.state = stateQuerying
			return m, nil

		case "ctrl+c", "q", "esc":
			m.state = stateAborted
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	mneu := m.key.Mneumonic()
	words := strings.Split(mneu, " ")

	mneuLines := ""
	for l := 0; l < len(words); l += 6 {
		mneuLines += strings.Join(words[l:l+6], " ") + "\n"
	}

	switch m.state {
	case stateQuerying:
		return fmt.Sprintf(
			"Generated Mnemonic:\n%s\nDo you like it? (y/n)\n",
			mneuLines)

	case stateAccepted:
		return fmt.Sprintf("Mnemonic accepted:\n%s\n", mneuLines)

	case stateAborted:
		return "Session Aborted\n"
	}
	zap.S().Panic("unexpected state")
	return ""
}

func Prompt(newEmergencyKeyFunc newEmergencyKeyFunc) (*emergency.EmergencyBackupKey, error) {
	p := tea.NewProgram(model{
		newEmergencyKeyFunc: newEmergencyKeyFunc,
		key:                 newEmergencyKeyFunc(),
		state:               stateQuerying,
	})

	m, err := p.Run()
	if err != nil {
		return nil, err
	}
	m2 := m.(model)
	switch m2.state {
	case stateAccepted:
		return m2.key, nil
	case stateAborted:
		return nil, common.ErrAbort{}
	default:
		zap.S().Fatalf("unexpected state: %v", m2.state)
	}

	return m.(model).key, nil
}
