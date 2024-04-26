package main

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hrishikesh/hash/database"
)

func passwordListUI(m Model) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("All passwords · total %d\n\n", len(m.allPassword)))
	for i, v := range m.allPassword {
		if i == m.passwordIndex {
			sb.WriteString(listItemSelectedHeaderStyle.Render(v.Name))
			sb.WriteString("\n")
			sb.WriteString(listItemSelectedStyle.Inherit(listItemStyle).Render(v.Email))
			sb.WriteString("\n\n")
		} else {
			sb.WriteString(listItemHeaderStyle.Render(v.Name))
			sb.WriteString("\n")
			sb.WriteString(listItemStyle.Render(v.Email))
			sb.WriteString("\n\n")
		}

	}

	sb.WriteString("ctrl+n Create new password ctrl+u update ctrl+d delete ")

	return sb.String()
}

func onPaswordListUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			if m.passwordIndex >= len(m.allPassword)-1 {
				m.passwordIndex = 0
			} else {
				m.passwordIndex++
			}

		case tea.KeyUp:
			if m.passwordIndex == 0 {
				m.passwordIndex = len(m.allPassword) - 1
			} else {
				m.passwordIndex--
			}
		case tea.KeyEnter:
			m.pageIndex = viewPasswordIndex

		case tea.KeyCtrlN:
			m.pageIndex = newPasswords

		case tea.KeyCtrlU:
			m.pageIndex = updatePasswords
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
			cmds := make([]tea.Cmd, 3)

			for i, v := range m.txtInputs {
				switch i {
				case 0:
					v.SetValue(m.allPassword[m.passwordIndex].Name)
				case 1:
					v.SetValue(m.allPassword[m.passwordIndex].Email)
				case 2:
					v.SetValue(string(m.allPassword[m.passwordIndex].Password))
				}

				m.txtInputs[i], cmds[i] = v.Update(msg)
			}
			return m, tea.Batch(cmds...)
		case tea.KeyCtrlD:
			err := database.DeletePassword(m.allPassword[m.passwordIndex].ID)
			if err != nil {
				m.err = err
				return m, nil
			}
			newAllPassword := make([]database.PWItem, 0)
			newAllPassword = append(newAllPassword, m.allPassword[:m.passwordIndex]...)
			newAllPassword = append(newAllPassword, m.allPassword[m.passwordIndex+1:]...)
			m.allPassword = newAllPassword
			return m, nil
		}
	}
	return m, nil
}
