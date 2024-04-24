package main

import tea "github.com/charmbracelet/bubbletea"

func readSaltFileTUI() tea.Msg {
	stringByte, err := readSaltFile()
	if err != nil {
		return err
	}
	return saltFound(stringByte)
}
