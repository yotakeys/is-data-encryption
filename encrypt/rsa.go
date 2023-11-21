package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

func EncryptRSA(msg, publicKey string) string {
	parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		panic(err)
	}
	// To get back to an *rsa.PublicKey, we need to first upgrade to the
	// ssh.CryptoPublicKey interface
	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	pub := pubCrypto.(*rsa.PublicKey)

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub,
		[]byte(msg),
		nil)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes)
}

func DecryptRSA(data, priv string) (string, error) {
	data2, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode([]byte(priv))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	decrypted, err := rsa.DecryptOAEP(sha256.New(),
		rand.Reader, key, data2, nil)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
