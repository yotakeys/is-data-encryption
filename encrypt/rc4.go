package encrypt

import (
	// ...
	"crypto/rc4"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

// RC4Encrypt melakukan enkripsi RC4
func RC4Encrypt(stringToEncrypt string) (encryptedString string, data map[string]interface{}, err error) {
	elapsedTimer := timerWithReturn("RC4Encrypt")
	defer elapsedTimer()

	time.Sleep(1 * time.Second)

	key := []byte(os.Getenv("ENCRYPT_KEY"))
	plaintext := []byte(stringToEncrypt)

	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return "", nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	cipher.XORKeyStream(ciphertext, plaintext)

	data = map[string]interface{}{
		"key":       os.Getenv("ENCRYPT_KEY"),
		"plaintext": string(plaintext),
		"block":     fmt.Sprintf("%v", cipher),
		"mode":      "RC4",
		"elapsed":   elapsedTimer().String(),
	}

	return fmt.Sprintf("%x", ciphertext), data, err
}

// RC4Decrypt melakukan dekripsi RC4
func RC4Decrypt(encryptedString string) (decryptedString string, err error) {
	key := []byte(os.Getenv("ENCRYPT_KEY"))
	enc, _ := hex.DecodeString(encryptedString)

	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(enc))
	cipher.XORKeyStream(plaintext, enc)

	return string(plaintext), nil
}
