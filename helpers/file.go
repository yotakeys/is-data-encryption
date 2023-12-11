package helpers

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"gin-gorm-clean-template/encrypt"
	"gin-gorm-clean-template/entity"
	"time"
)

type (
	DigitalSignature struct {
		Name     string
		Email    string
		DateTime string
	}

	DigitalSignatureContent struct {
		Data      []byte
		Signature []byte
		PublicKey []byte
	}
)

func WriteDigitalSignature(content []byte, user entity.User, userName string) ([]byte, error) {
	var messagesAdded []byte
	hash := sha256.Sum256(content)

	digitalSignature := DigitalSignature{
		Name:     userName,
		Email:    user.Email,
		DateTime: time.Now().Format("2002-10-10 07:00:00"),
	}

	digitalSignatureBytes, err := json.Marshal(digitalSignature)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, err := encrypt.ParsePrivateKeyFromPEM(user.PrivateKey)
	if err != nil {
		return nil, err
	}

	data_signature, err := rsa.SignPKCS1v15(nil, rsaPrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}

	messagesAdded = append(messagesAdded, content...)
	messagesAdded = append(messagesAdded, []byte("\n%Hash")...)
	messagesAdded = append(messagesAdded, data_signature...)
	messagesAdded = append(messagesAdded, []byte("\n%Signature")...)
	messagesAdded = append(messagesAdded, digitalSignatureBytes...)
	messagesAdded = append(messagesAdded, []byte("\n%PublicKey")...)
	messagesAdded = append(messagesAdded, []byte(user.PublicKey)...)

	return messagesAdded, nil
}

func ReadDigitalSignature(content []byte) (DigitalSignatureContent, error) {
	var buffer []byte

	tokens := []string{"%Hash", "%Signature", "%PublicKey"}

	results := [3]string{}
	ctr := 0
	idx := 0

	stringFileContent := string(content)
	for idx < len(stringFileContent) {
		if ctr <= 2 &&
			idx+len(tokens[ctr]) < len(stringFileContent) &&
			stringFileContent[idx:idx+len(tokens[ctr])] == tokens[ctr] {

			if ctr > 0 && len(buffer) > 0 {
				results[ctr-1] = string(buffer[:len(buffer)-1])
				buffer = []byte{}
			}

			idx += len(stringFileContent[idx : idx+len(tokens[ctr])])
			ctr++
		}

		if ctr > 0 {
			buffer = append(buffer, stringFileContent[idx])
		}

		if idx == len(stringFileContent)-1 && ctr != 0 {
			results[ctr-1] = string(buffer)
		}

		idx++
	}

	return DigitalSignatureContent{
		Data:      []byte(results[0]),
		Signature: []byte(results[1]),
		PublicKey: []byte(results[2]),
	}, nil
}
