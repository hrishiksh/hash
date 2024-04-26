package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hrishikesh/hash/database"
)

func newPasswordUI(m Model) string {
	var sb strings.Builder
	sb.WriteString(headlineStyle.Render("Add new password\n"))
	sb.WriteString("\n")
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
	sb.WriteString("\n")
	sb.WriteString(faintTextStyle.Render(m.err.Error()))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%s %s", boldFaintTextStyle.Render("ctrl+←"), faintTextStyle.Render("back ")))
	sb.WriteString(fmt.Sprintf("%s %s", boldFaintTextStyle.Render("・ ctrl+c/esc"), faintTextStyle.Render("exit")))
	return sb.String()
}

func onNewPasswordUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			if m.focusIndex >= len(m.txtInputs) {
				m.focusIndex = 0
			} else {
				m.focusIndex++
			}

		case tea.KeyUp:
			if m.focusIndex == 0 {
				m.focusIndex = len(m.txtInputs)
			} else {
				m.focusIndex--
			}

		case tea.KeyCtrlLeft:
			m.pageIndex = allPasswordsPage

		case tea.KeyEnter:
			if m.focusIndex >= 2 {
				encryptedPassword, err := encryptMessage([]byte(m.txtInputs[2].Value()), m.secretKey)
				if err != nil {
					return m, func() tea.Msg { return err }
				}
				id, err := database.AddNewPassword(m.txtInputs[0].Value(), m.txtInputs[1].Value(), byteToHex(encryptedPassword))
				if err != nil {
					return m, func() tea.Msg { return err }
				}
				m.allPassword = append(m.allPassword, database.PWItem{ID: id, Name: m.txtInputs[0].Value(), Email: m.txtInputs[1].Value(), Password: byteToHex(encryptedPassword)})
				m.pageIndex = allPasswordsPage
				return m, nil
			} else {
				m.focusIndex++
			}
		}

	case error:
		m.err = msg
		return m, nil
	}

	cmds := make([]tea.Cmd, 3)
	for i, v := range m.txtInputs {
		if i == m.focusIndex {
			v.Focus()
		} else {
			v.Blur()
		}
		m.txtInputs[i], cmds[i] = v.Update(msg)

	}

	return m, tea.Batch(cmds...)
}
