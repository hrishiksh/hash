package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	btnStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1).Bold(true)
	activeBtnStyle  = btnStyle.Copy().Background(lipgloss.Color("#5539a8"))
	faintText       = lipgloss.NewStyle().Faint(true)
	displayMsgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef233c")).Bold(true)
)

type model struct {
	salt                     []byte
	masterPasswordFocusIndex int
	masterPassword           textinput.Model
	authenticated            bool
	focusIndex               int
	txtInputs                []textinput.Model
	err                      error
	displayMsg               string
}

type saltFound []byte

func initialModel() model {

	m := model{
		authenticated: false,
		focusIndex:    0,
		txtInputs:     make([]textinput.Model, 3),
		err:           nil,
	}

	// setting master password text filed
	masterPasswordTxtField := textinput.New()
	masterPasswordTxtField.CharLimit = 200
	masterPasswordTxtField.Width = 50
	masterPasswordTxtField.Placeholder = "Master password"
	masterPasswordTxtField.EchoMode = textinput.EchoPassword
	masterPasswordTxtField.EchoCharacter = '*'
	m.masterPassword = masterPasswordTxtField

	for i := range m.txtInputs {

		t := textinput.New()
		t.CharLimit = 156
		t.Width = 30

		switch i {
		case 0:
			t.Placeholder = "Name"
			t.Focus()

		case 1:
			t.Placeholder = "Email"

		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
		}

		m.txtInputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return readSaltFileTUI
}

func (m model) View() string {

	var sb strings.Builder

	if m.authenticated {
		sb.WriteString("Add new password\n\n")
		for _, v := range m.txtInputs {
			sb.WriteString(v.View())
			sb.WriteString("\n")
		}

		sb.WriteString("\n")
		if m.focusIndex == len(m.txtInputs) {
			sb.WriteString(activeBtnStyle.Render("Save"))
		} else {
			sb.WriteString(btnStyle.Render("Save"))

		}

		sb.WriteString(faintText.Render("\n(press ctrl+c or esc to exit)\n"))

	} else {
		sb.WriteString("Login with your master password\n\n")
		sb.WriteString(m.masterPassword.View())
		sb.WriteString("\n\n")
		if m.masterPasswordFocusIndex == 1 {
			sb.WriteString(activeBtnStyle.Render("Submit"))

		} else {
			sb.WriteString(btnStyle.Render("Submit"))

		}

		sb.WriteString("\n")
		sb.WriteString(displayMsgStyle.Render(m.displayMsg))
		sb.WriteString(faintText.Render("\n(press ctrl+c or esc to exit)\n"))

	}

	return sb.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyDown:
			if m.authenticated {
				if m.focusIndex >= len(m.txtInputs) {
					m.focusIndex = 0
				} else {
					m.focusIndex++
				}

			} else {
				if m.masterPasswordFocusIndex >= 1 {
					m.masterPasswordFocusIndex = 0
				} else {
					m.masterPasswordFocusIndex++
				}
			}

		case tea.KeyUp:
			if m.authenticated {
				if m.focusIndex == 0 {
					m.focusIndex = len(m.txtInputs)
				} else {
					m.focusIndex--
				}

			} else {
				if m.masterPasswordFocusIndex == 0 {
					m.masterPasswordFocusIndex = 1
				} else {
					m.masterPasswordFocusIndex--
				}
			}

		case tea.KeyEnter:
			if err := varifySaltAndPassword(m.salt, []byte(m.masterPassword.Value())); err != nil {
				m.displayMsg = "Password and salt doesn't pair off"
			} else {
				m.authenticated = true
			}

		}
	case saltFound:
		m.salt = []byte(msg)
	case error:
		if msg == os.ErrNotExist {
			m.displayMsg = "Salt Not found"
		} else {
			m.err = msg
		}
		return m, nil
	}

	if m.authenticated {
		updatedCmds := []tea.Cmd{}
		for i, v := range m.txtInputs {
			if i == m.focusIndex {
				v.Focus()
			} else {
				v.Blur()
			}

			var cmd tea.Cmd
			v, cmd = v.Update(msg)
			m.txtInputs[i] = v
			updatedCmds = append(updatedCmds, cmd)
		}

		return m, tea.Batch(updatedCmds...)
	} else {
		if m.masterPasswordFocusIndex == 0 {
			m.masterPassword.Focus()
		} else {
			m.masterPassword.Blur()
		}
		tm, tc := m.masterPassword.Update(msg)
		m.masterPassword = tm
		return m, tc
	}

}
