package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	focusIndex        int
	focusShift        bool
	inputs            []textinput.Model
	err               error
	data              *EmailSetup
	testingConnection bool
	spinner           spinner.Model
	testInitial       bool
}

func (m model) Host() string {
	return m.inputs[0].Value()
}
func (m model) User() string {
	return m.inputs[2].Value()
}
func (m model) Password() string {
	return m.inputs[3].Value()
}

func (m model) Port() (int32, error) {
	port, err := strconv.Atoi(m.inputs[1].Value())
	if err != nil {
		return 0, fmt.Errorf("Invalid port number: %v", m.inputs[1].Value())
	}
	return int32(port), nil
}

func (m model) UpdateSettings() error {
	port, err := m.Port()
	if err != nil {
		return err
	}

	if m.User() == "" {
		return fmt.Errorf("User field is required")
	}

	m.data.Host = m.Host()
	m.data.Port = port
	m.data.User = m.User()
	m.data.Password = m.Password()

	return nil
}

type Button struct {
	focused string
	blurred string
}

func newBtn(label string) Button {
	return Button{
		focused: focusedStyle.Copy().Render(
			fmt.Sprintf("[ %s ]", label),
		),
		blurred: fmt.Sprintf("[ %s ]",
			blurredStyle.Render(label),
		),
	}
}

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("70"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("124"))

	submitBtn = newBtn("Submit")
	skipBtn   = newBtn("Skip")
)

func initialModel(s *EmailSetup) *model {
	var host = textinput.New()
	host.Prompt = "Host:\t\t"
	host.PromptStyle = focusedStyle
	host.TextStyle = focusedStyle
	host.SetValue(s.Host)
	host.SetSuggestions([]string{"sandbox.smtp.mailtrap.io"})

	var port = textinput.New()
	port.Prompt = "Port:\t\t"
	port.SetValue(fmt.Sprintf("%d", s.Port))

	var user = textinput.New()
	user.Prompt = "User:\t\t"
	user.SetValue(s.User)

	var password = textinput.New()
	password.Prompt = "Password:\t"
	password.SetValue(s.Password)

	m := model{
		inputs:      []textinput.Model{host, port, user, password},
		data:        s,
		spinner:     spinner.New(spinner.WithSpinner(spinner.Points)),
		testInitial: !s.noAuto && s != nil && s.User != "" && s.Password != "" && s.Host != "" && s.Port != 0,
		focusIndex:  0,
		focusShift:  false,
	}

	return &m
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.spinner.Tick,
		tea.Printf("[ SMTP configuration ]"),
		m.Focus(0),
	)
}

func (m *model) Focus(index int) tea.Cmd {
	var cmd tea.Cmd
	for i := range m.inputs {
		if i == m.focusIndex {
			// Set focused state
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
			cmd = m.inputs[i].Focus()
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = noStyle
		m.inputs[i].TextStyle = noStyle
	}
	return cmd
}

func (m *model) TestConnection() {
	m.data.connectionOK = false
	m.testingConnection = true
	m.testInitial = false
	err := m.UpdateSettings()
	if err != nil {
		m.err = fmt.Errorf("Failed to get settings: %v", err)
		m.testingConnection = false
		return
	}
	_, err = m.data.Dialer().Dial()
	if err != nil {
		m.err = fmt.Errorf(
			"🔴 Connection to the SMTP server failed using this configuration.",
		)
		m.testingConnection = false
		return
	}
	m.testingConnection = false
	m.data.connectionOK = true
	return
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.data.connectionOK {
		return m, tea.Quit
	}
	if m.testInitial {
		go m.TestConnection()
	}
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEsc:
			m.focusIndex = len(m.inputs)
			m.focusShift = true

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter, tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
			if msg.Type == tea.KeyEnter && m.focusIndex == len(m.inputs) {
				if m.focusShift {
					m.data.skip = true
					return m, tea.Quit
				} else {
					m.testingConnection = true
					go m.TestConnection()
					return m, m.updateInputs(msg)
				}
			}

			// Cycle indexes
			switch msg.Type {
			case tea.KeyUp:
				m.focusIndex--
			case tea.KeyDown:
				m.focusIndex++
			case tea.KeyLeft:
				m.focusShift = false || m.focusIndex < len(m.inputs)
			case tea.KeyRight:
				m.focusShift = true && m.focusIndex == len(m.inputs)
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			return m, tea.Batch(m.Focus(m.focusIndex))
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		if i != 1 {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		} else {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if strings.Contains("1234567890", msg.String()) || msg.Type == tea.KeyBackspace {
					m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
				}
			}
		}
	}
	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	if m.testingConnection {
		b.WriteString(focusedStyle.Render("\n\n", m.spinner.View(), "\tTesting connection\n"))
	} else if m.data.connectionOK {
		b.WriteString(successStyle.Render("\n\n🟢 Connection succeeded"))
	} else {
		submit := &submitBtn.blurred
		if m.focusIndex == len(m.inputs) && !m.focusShift {
			submit = &submitBtn.focused
		}
		skip := &skipBtn.blurred
		if m.focusIndex == len(m.inputs) && m.focusShift {
			skip = &skipBtn.focused
		}
		fmt.Fprintf(&b, "\n\n%s\t%s\n\n", *submit, *skip)
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render(m.err.Error()))
	}

	return b.String()
}
