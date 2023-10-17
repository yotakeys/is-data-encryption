package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func AESEncrypt(stringToEncrypt string) (encryptedString string, data map[string]interface{}, err error) {
	elapsedTimer := timerWithReturn("AESEncrypt")
	defer elapsedTimer()

	time.Sleep(1 * time.Second)

	key := []byte(os.Getenv("ENCRYPT_KEY"))
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	data = map[string]interface{}{
		"key":       os.Getenv("ENCRYPT_KEY"),
		"plaintext": string(plaintext),
		"block":     fmt.Sprintf("%d", block.BlockSize()),
		"aes-gcm":   fmt.Sprintf("%v", aesGCM),
		"nonce":     hex.EncodeToString(nonce),
		"mode":      "AES",
		"elapsed":   elapsedTimer().String(),
	}

	return fmt.Sprintf("%x", ciphertext), data, err
}

func AESDecrypt(encryptedString string) (decryptedString string, err error) {
	defer func() {
		if r := recover(); r != nil {
			decryptedString = ""
			err = errors.New("error in decrypting")
		}
	}()

	key := []byte(os.Getenv("ENCRYPT_KEY"))
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", nil
	}

	return string(plaintext), nil
}
