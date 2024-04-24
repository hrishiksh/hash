package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hrishikesh/hash/database"
)

func main() {

	err := database.InitDB("password.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		return
	}

	// salt, err := generateSalt([]byte("hello-world"))
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// err = os.WriteFile("salt.txt", byteToHex(salt), 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// saltContent, err := os.ReadFile("salt.txt")
	// if err != nil {
	// 	if errors.Is(err, os.ErrNotExist) {
	// 		log.Fatal("file not exist")
	// 		return
	// 	}
	// 	log.Fatal(err)
	// 	return
	// }

	// saltbyte, err := hexToByte(saltContent)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// secretKey := generateSecretKey([]byte("hello-world"), saltbyte)

	// encryptHex, err := encryptMessage([]byte("I am a good boy"), secretKey)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// err = os.WriteFile("encryptmsg.txt", encryptHex, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// encryptMsg, err := os.ReadFile("encryptmsg.txt")
	// if err != nil {
	// 	if errors.Is(err, os.ErrNotExist) {
	// 		log.Fatal("file not exist")
	// 		return
	// 	}
	// 	log.Fatal(err)
	// 	return
	// }

	// encryptByte, err := hexToByte(encryptMsg)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// msg, ok := decryptMsg(encryptByte, secretKey)
	// if !ok {
	// 	log.Fatal("something went wrong")
	// 	return
	// }

	// fmt.Printf("%s\n", msg)

}
