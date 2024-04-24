package main

import (
	"encoding/hex"
	"os"
)

func byteToHex(src []byte) []byte {
	out := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(out, src)
	return out
}

func hexToByte(src []byte) ([]byte, error) {
	out := make([]byte, hex.DecodedLen(len(src)))

	_, err := hex.Decode(out, src)
	return out, err
}

func readSaltFile() ([]byte, error) {
	hexSalt, err := os.ReadFile("salt.txt")
	if err != nil {
		return nil, err
	}
	return hexToByte(hexSalt)
}
