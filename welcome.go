package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func welcomeUI(m Model) string {
	var sb strings.Builder

	sb.WriteString(headlineStyle.Render("Welcome to Hash. A password manager for everyone"))
	sb.WriteString("\n")
	sb.WriteString("Enter a master password to generate your secret key 🔑\n")
	sb.WriteString(lBfaintTStyle.Render("A Secret key will be generated based on your master password and stored on your device.\nCombination of master password and secret key is used to decrypt your passwords.\nPlease remember your master password. There is no way to reset the master password."))

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
	sb.WriteString(fmt.Sprintf("%s %s", boldFaintTextStyle.Render("ctrl+ →"), faintTextStyle.Render("login ")))
	sb.WriteString(fmt.Sprintf("%s %s", boldFaintTextStyle.Render("・ ctrl+c/esc"), faintTextStyle.Render("exit")))
	return sb.String()
}

func onWelcomeUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {

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
			salt, err := generateSalt([]byte(m.masterPassword.Value()))
			if err != nil {
				m.err = err
				return m, nil
			}

			err = os.WriteFile(filepath.Join(configDir, "salt.txt"), byteToHex(salt), 0644)
			if err != nil {
				m.err = err
				return m, nil
			}

			m.salt = salt
			m.err = errors.New("secret key generated. press ctrl+ → to login page")

		case tea.KeyCtrlRight:
			m.pageIndex = loginPage
			return m, nil

		}
	case error:
		if errors.Is(msg, os.ErrExist) {
			m.err = errors.New("secret key exist")
			m.pageIndex = loginPage
			return m, nil
		}
		m.err = msg
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
