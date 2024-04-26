package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func loginUI(m Model) string {
	var sb strings.Builder
	sb.WriteString(headlineStyle.Render("Login with your master password:\n"))
	// a fix to make text input align again. Otherwise textinput align itself
	// in the middle when I add the headline style.
	sb.WriteString("\n")
	sb.WriteString(m.masterPassword.View())
	sb.WriteString("\n\n")
	if m.masterPasswordFocusIndex == 1 {
		sb.WriteString(activeBtnStyle.Render("Submit"))

	} else {
		sb.WriteString(btnStyle.Render("Submit"))

	}

	sb.WriteString("\n")
	sb.WriteString(faintTextStyle.Render(m.err.Error()))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%s %s", boldFaintTextStyle.Render("ctrl+c/esc"), faintTextStyle.Render("exit")))
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
				m.err = err
				return m, nil
			} else {
				m.secretKey = generateSecretKey([]byte(m.masterPassword.Value()), m.salt)
				m.pageIndex = allPasswordsPage
			}

		}
	case saltFound:
		m.salt = []byte(msg)
	case error:
		if msg == os.ErrNotExist {
			m.err = errors.New("salt not found")
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
