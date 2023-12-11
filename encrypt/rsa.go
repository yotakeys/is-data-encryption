package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func EncryptRSA(data string, publicKey string) string {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return ""
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return ""
	}

	encData, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(data))
	if err != nil {
		return ""
	}

	sEnc := base64.StdEncoding.EncodeToString(encData)

	return sEnc
}

func DecryptRSA(data string, privateKey string) string {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return ""
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}

	sDec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}

	decData, err := rsa.DecryptPKCS1v15(rand.Reader, priv, []byte(sDec))
	if err != nil {
		return ""
	}

	return string(decData)
}

func ParsePrivateKeyFromPEM(keyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func ParsePublicKeyFromPEM(pemData string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Assert that the parsed key is an RSA public key
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("parsed key is not an RSA public key")
	}

	return publicKey, nil
}
