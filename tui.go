package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hrishikesh/hash/database"
)

const (
	loginMasterPassword = iota
	allPasswords
	newPasswords
	viewPasswordIndex
	updatePasswords
	deletePasswords
)

var (
	btnStyle                    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#7D56F4")).Padding(0, 1).Bold(true)
	activeBtnStyle              = btnStyle.Copy().Background(lipgloss.Color("#5539a8"))
	faintText                   = lipgloss.NewStyle().Faint(true)
	displayMsgStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef233c")).Bold(true)
	listItemStyle               = lipgloss.NewStyle().Padding(0, 1).Faint(true)
	listItemHeaderStyle         = listItemStyle.Copy().Bold(true).UnsetFaint()
	listItemSelectedStyle       = listItemStyle.Copy().BorderStyle(lipgloss.NormalBorder()).BorderLeft(true).BorderLeftForeground(lipgloss.Color("#7D56F4"))
	listItemSelectedHeaderStyle = listItemSelectedStyle.Copy().Bold(true).UnsetFaint()
)

type model struct {
	pageIndex                int
	salt                     []byte
	masterPasswordFocusIndex int
	masterPassword           textinput.Model
	secretKey                [32]byte
	focusIndex               int
	txtInputs                []textinput.Model
	allPassword              []database.PWItem
	passwordIndex            int
	err                      error
	displayMsg               string
}

type saltFound []byte

func initialModel() model {

	m := model{
		pageIndex:     loginMasterPassword,
		focusIndex:    0,
		txtInputs:     make([]textinput.Model, 3),
		err:           errors.New(""),
		passwordIndex: 0,
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

	// Getting all the password
	var err error
	m.allPassword, err = database.ReadAllPasswords()
	if err != nil {
		m.err = err
	}

	return m
}

func (m model) Init() tea.Cmd {
	return readSaltFileTUI
}

func (m model) View() string {
	switch m.pageIndex {
	case loginMasterPassword:
		return loginMaterPasswordView(m)
	case allPasswords:
		return allPasswordView(m)
	case newPasswords:
		return newPasswordView(m)
	case viewPasswordIndex:
		return viewPasswordPage(m)
	}
	return fmt.Sprintln("No page selected")
}

func loginMaterPasswordView(m model) string {
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

func allPasswordView(m model) string {
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

	sb.WriteString("ctrl+n Create new password")

	return sb.String()
}

func viewPasswordPage(m model) string {
	var sb strings.Builder
	sb.WriteString(m.allPassword[m.passwordIndex].Name)
	sb.WriteString(fmt.Sprintf("\n\nEmail: %s\n", m.allPassword[m.passwordIndex].Email))
	sb.WriteString(fmt.Sprintf("Password: %s", m.allPassword[m.passwordIndex].Password))
	sb.WriteString("\n\n")
	sb.WriteString("(ctrl+a decrypt · ← back)")
	return sb.String()

}

func newPasswordView(m model) string {
	var sb strings.Builder
	sb.WriteString("🔑 Add new password\n\n")
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
	sb.WriteString(displayMsgStyle.Render(m.displayMsg))
	sb.WriteString(faintText.Render("\n(press ctrl+c or esc to exit)\n"))
	sb.WriteString("\n")
	return sb.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		}

	}

	switch m.pageIndex {
	case loginMasterPassword:
		return authMasterPasswordUpdate(msg, m)
	case allPasswords:
		return allPasswordUpdate(msg, m)
	case viewPasswordIndex:
		return viewPasswordUpdate(msg, m)
	case newPasswords:
		return newPasswordUpdate(msg, m)
	}
	return m, nil
}

func authMasterPasswordUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
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

func allPasswordUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
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
		}
	}
	return m, nil
}

func viewPasswordUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
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

func newPasswordUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
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

		case tea.KeyEnter:
			if m.focusIndex >= 2 {
				encryptedPassword, err := encryptMessage([]byte(m.txtInputs[2].Value()), m.secretKey)
				if err != nil {
					return m, func() tea.Msg { return err }
				}
				err = database.AddNewPassword(m.txtInputs[0].Value(), m.txtInputs[1].Value(), byteToHex(encryptedPassword))
				if err != nil {
					return m, func() tea.Msg { return err }
				}
				return m, tea.Quit
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
