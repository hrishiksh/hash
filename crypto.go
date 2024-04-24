package main

import (
	"crypto/rand"
	"crypto/sha512"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/pbkdf2"
)

func generateSalt(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, 10)
}

func varifySaltAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func generateSecretKey(password []byte, salt []byte) [32]byte {
	// In 2023, OWASP recommended to use 600,000 iterations for PBKDF2-HMAC-SHA256 and 210,000 for PBKDF2-HMAC-SHA512
	key := pbkdf2.Key(password, salt, 210_000, 32, sha512.New)

	var secretkey [32]byte
	copy(secretkey[:], key)

	return secretkey
}

func encryptMessage(message []byte, secretKey [32]byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	encryptMsg := secretbox.Seal(nonce[:], message, &nonce, &secretKey)

	return encryptMsg, nil
}

func decryptMsg(encodedMsg []byte, secretKey [32]byte) ([]byte, bool) {
	var nonce [24]byte
	copy(nonce[:], encodedMsg[:24])

	return secretbox.Open(nil, encodedMsg[24:], &nonce, &secretKey)

}
