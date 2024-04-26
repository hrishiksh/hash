package main

import (
	"errors"
	"fmt"

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

type Model struct {
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

func initialModel() Model {

	m := Model{
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

func (m Model) Init() tea.Cmd {
	return readSaltFileTUI
}

func (m Model) View() string {
	switch m.pageIndex {
	case loginMasterPassword:
		return loginUI(m)
	case allPasswords:
		return passwordListUI(m)
	case newPasswords:
		return newPasswordUI(m)
	case viewPasswordIndex:
		return viewPasswordUI(m)
	case updatePasswords:
		return updatePasswordUI(m)
	}
	return fmt.Sprintln("No page selected")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		}

	}

	switch m.pageIndex {
	case loginMasterPassword:
		return onLoginUpdate(msg, m)
	case allPasswords:
		return onPaswordListUpdate(msg, m)
	case viewPasswordIndex:
		return onviewPasswordUpdate(msg, m)
	case newPasswords:
		return onNewPasswordUpdate(msg, m)
	case updatePasswords:
		return onPasswordUpdate(msg, m)
	}
	return m, nil
}
