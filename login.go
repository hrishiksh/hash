package main

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func loginUI(m Model) string {
	var sb strings.Builder
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
	sb.WriteString("\n")
	return sb.String()
}

func onLoginUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyDown:

			if m.masterPasswordFocusIndex >= 1 {
				m.masterPasswordFocusIndex = 0
			} else {
				m.masterPasswordFocusIndex++
			}

		case tea.KeyUp:

			if m.masterPasswordFocusIndex == 0 {
				m.masterPasswordFocusIndex = 1
			} else {
				m.masterPasswordFocusIndex--
			}

		case tea.KeyEnter:

			if err := varifySaltAndPassword(m.salt, []byte(m.masterPassword.Value())); err != nil {
				m.displayMsg = "Password and salt doesn't pair off"
			} else {
				m.secretKey = generateSecretKey([]byte(m.masterPassword.Value()), m.salt)
				m.pageIndex = allPasswords
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

	if m.masterPasswordFocusIndex == 0 {
		m.masterPassword.Focus()
	} else {
		m.masterPassword.Blur()
	}
	tm, tc := m.masterPassword.Update(msg)
	m.masterPassword = tm
	return m, tc

}
