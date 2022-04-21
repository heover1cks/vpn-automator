package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/heover1cks/vpn-automator/config"
	"github.com/heover1cks/vpn-automator/pkg/client"
	"github.com/heover1cks/vpn-automator/pkg/utils"
	"github.com/muesli/termenv"
)

var (
	term           = termenv.EnvColorProfile()
	selected       = ""
	selectedStatus = ""
)

type Model struct {
	conf   config.Config
	vpns   map[string]string
	cursor int
}

func InitialModel(conf config.Config) Model {
	vpns := client.GetConnectionStatus(conf)
	return Model{
		conf: client.GetConnectionStatus(conf),
		vpns: vpns.RetrieveAccountAliases(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.vpns)-1 {
				m.cursor++
			}
		case "help", "h":
			printDescription()
		case "enter", " ":
			m.manageConnection()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) manageConnection() {
	m.changeSelected()
	client.ConnectVPN(m.conf, selected, selectedStatus == client.Connected)
}

func (m Model) changeSelected() {

}

func (m Model) View() string {
	s := "⚡️Select VPN to connect\n\n"
	idx := 0
	vpns := utils.SortedStringKeys(m.vpns)
	for _, vpn := range vpns {
		cursor := " " // no cursor
		if m.cursor == idx {
			cursor = ">" // cursor!
			selected = vpn
			selectedStatus = m.vpns[vpn]
		}
		s += fmt.Sprintf("[%s]%s %s\n", statusColorFg(m.vpns[vpn]), colorFg(cursor, "159"), vpn)
		idx += 1
	}
	return s //+ printDescription()
}

func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func statusColorFg(val string) string {
	color := statusColorSelector(val)
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func statusColorSelector(status string) string {
	if status == client.Connected {
		return "10"
	} else if status == client.ConnectedUnknown {
		return "11"
	} else if status == client.Disconnected {
		return "9"
	} else {
		return "13"
	}
}

func printDescription() string {
	ret := "\n"
	ret += fmt.Sprintf(" *[%s]: Process is alive, but not sure which network is connected\n", statusColorFg(client.ConnectedUnknown))
	ret += fmt.Sprintf(" *[%s]: Connected\n", statusColorFg(client.Connected))
	ret += fmt.Sprintf(" *[%s]: Disconnected or process is dead\n", statusColorFg(client.Disconnected))
	return ret
}
