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
	loginPage = iota
	allPasswordsPage
	newPasswordPage
	viewPasswordPage
	updatePasswordPage
)

var (
	boldTextStyle               = lipgloss.NewStyle().Bold(true)
	boldFaintTextStyle          = boldTextStyle.Copy().Faint(true)
	faintTextStyle              = lipgloss.NewStyle().Faint(true)
	headlineStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff006e")).Bold(true)
	boldtextWBackground         = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#8338ec")).Padding(0, 1)
	btnStyle                    = boldTextStyle.Copy().Foreground(lipgloss.Color("#FFFFFF")).Padding(0, 1).Border(lipgloss.NormalBorder())
	activeBtnStyle              = btnStyle.Copy().Foreground(lipgloss.Color("#8338ec")).BorderForeground(lipgloss.Color("#8338ec"))
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
}

type saltFound []byte

func initialModel() Model {

	m := Model{
		pageIndex:     loginPage,
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
	case loginPage:
		return loginUI(m)
	case allPasswordsPage:
		return passwordListUI(m)
	case newPasswordPage:
		return newPasswordUI(m)
	case viewPasswordPage:
		return viewPasswordUI(m)
	case updatePasswordPage:
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
	case loginPage:
		return onLoginUpdate(msg, m)
	case allPasswordsPage:
		return onPaswordListUpdate(msg, m)
	case viewPasswordPage:
		return onviewPasswordUpdate(msg, m)
	case newPasswordPage:
		return onNewPasswordUpdate(msg, m)
	case updatePasswordPage:
		return onPasswordUpdate(msg, m)
	}
	return m, nil
}
