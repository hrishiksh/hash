package main

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func viewPasswordUI(m Model) string {
	var sb strings.Builder
	sb.WriteString(m.allPassword[m.passwordIndex].Name)
	sb.WriteString(fmt.Sprintf("\n\nEmail: %s\n", m.allPassword[m.passwordIndex].Email))
	sb.WriteString(fmt.Sprintf("Password: %s", m.allPassword[m.passwordIndex].Password))
	sb.WriteString("\n\n")
	sb.WriteString("(ctrl+a decrypt · ← back)")
	return sb.String()

}

func onviewPasswordUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlA:
			encryptedPasswordHex := m.allPassword[m.passwordIndex].Password
			encryptedPassword, err := hexToByte(encryptedPasswordHex)
			if err != nil {
				m.err = err
				return m, nil
			}
			decryptPassword, ok := decryptMsg(encryptedPassword, m.secretKey)
			if !ok {
				m.err = errors.New("not decrypted")
				return m, nil
			}

			m.allPassword[m.passwordIndex].Password = decryptPassword
			return m, nil

		case tea.KeyLeft:
			m.pageIndex = allPasswords
		}
	}

	return m, nil

}
