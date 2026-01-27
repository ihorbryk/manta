package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	work     = 25 * 60
	rest     = 5 * 60
	WORKTIME = "work"
	RESTTIME = "rest"
)

type Mapping map[string]int

var mapping = Mapping{
	WORKTIME: work,
	RESTTIME: rest,
}

var choices = []string{WORKTIME, RESTTIME}

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type model struct {
	progress progress.Model
	timeLeft int
	timeType string
	cursor   int
	choice   string
	pause    bool
}

func NewModel() model {
	return model{
		progress: progress.New(progress.WithDefaultGradient()),
		timeLeft: 0,
		timeType: WORKTIME,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			switch choices[m.cursor] {
			case WORKTIME:
				m.timeLeft = work
				m.timeType = WORKTIME
			case RESTTIME:
				m.timeLeft = rest
				m.timeType = RESTTIME
			}

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case " ":
			m.pause = !m.pause

		case "esc":
			m.timeLeft = 0
			m.pause = false

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.pause {
			return m, tickCmd()
		}

		if m.progress.Percent() == 1.0 && m.timeLeft == 0 {
			_ = notify(fmt.Sprintf("Time to %s is left", m.timeType), "")
		}

		m.timeLeft -= 1

		percent := 0.0

		if m.timeType == WORKTIME {
			percent = 1.0 - float64(m.timeLeft)/float64(work)
		}

		if m.timeType == RESTTIME {
			percent = 1.0 - float64(m.timeLeft)/float64(rest)
		}

		cmd := m.progress.SetPercent(float64(percent))

		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.timeLeft <= 0 {
		s := strings.Builder{}
		s.WriteString("Choose time type:\n")

		for i := 0; i < len(choices); i++ {
			if m.cursor == i {
				s.WriteString("[•] ")
			} else {
				s.WriteString("[ ] ")
			}
			s.WriteString(choices[i])
			totalTime := mapping[choices[i]]
			minutes := (totalTime % 3600) / 60
			s.WriteString(fmt.Sprintf(" (%02dm)", minutes))
			s.WriteString("\n")
		}
		s.WriteString("\n(press q to quit)\n")

		return s.String()
	}

	pad := strings.Repeat(" ", padding)

	minutes := (m.timeLeft % 3600) / 60
	seconds := m.timeLeft - minutes*60

	pause := "▶️"
	if m.pause {
		pause = "⏸️"
	}

	return "\n" +
		pad + fmt.Sprintf("%02dm/%02dm", (work%3600)/60, (rest%3600)/60) + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + fmt.Sprintf("%s -> %02dm%02ds %v", m.timeType, minutes, seconds, pause) +
		pad + helpStyle("Press 'q' key to quit")
}
