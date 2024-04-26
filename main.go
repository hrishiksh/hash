package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hrishikesh/hash/database"
)

var configDir string

func main() {

	var err error
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
		return
	}

	// creating hash directory in user cofig directory
	configDir = filepath.Join(userConfigDir, "hash")
	err = os.Mkdir(configDir, 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
		return
	}

	err = database.InitDB(filepath.Join(configDir, "password.db"))
	if err != nil {
		log.Fatal(err)
		return
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
