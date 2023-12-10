package helpers

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gin-gorm-clean-template/encrypt"
	"gin-gorm-clean-template/entity"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

const (
	PATH                = "storage"
	DataCommentKey      = "DataKeyF-02_"
	SignatureCommentKey = "SignatureKeyF-02_"
	PublicKeyCommentKey = "PublicKeyKeyF-02_"
	API_URL             = "www.isf.sre-its.com/static"
)

type (
	Signing struct {
		Name    string
		Email   string
		Release string
	}

	ReadContents struct {
		Data      []byte
		Signature []byte
		PublicKey []byte
	}
)

func UploadFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func UploadFileSuccess(file *multipart.FileHeader, path string) error {
	parts := strings.Split(path, "/")

	fileId := parts[2]
	directoryPath := fmt.Sprintf("%s/%s/%s", PATH, parts[0], parts[1])

	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if err := os.MkdirAll(directoryPath, 0777); err != nil {
			return err
		}
	}

	filePath := fmt.Sprintf("%s/%s", directoryPath, fileId)

	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	fileData, err := ioutil.ReadAll(uploadedFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, fileData, 0666)
	if err != nil {
		return err
	}

	return nil
}

func Write(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetExtension(filename string) string {
	parts := strings.Split(filename, ".")
	extension := parts[len(parts)-1]

	return extension
}

func WriteContent(content []byte, user entity.User, userName string) ([]byte, error) {
	var messagesAdded []byte
	hash := sha256.Sum256(content)

	signing := Signing{
		Name:    userName,
		Email:   user.Email,
		Release: time.Now().Format("2006-01-02 15:04:05"),
	}

	signingBytes, err := json.Marshal(signing)
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

	messagesAdded = append(messagesAdded, []byte("\n%"+DataCommentKey)...)
	messagesAdded = append(messagesAdded, data_signature...)

	messagesAdded = append(messagesAdded, []byte("\n%"+SignatureCommentKey)...)
	messagesAdded = append(messagesAdded, signingBytes...)

	messagesAdded = append(messagesAdded, []byte("\n%"+PublicKeyCommentKey)...)
	messagesAdded = append(messagesAdded, []byte(user.PublicKey)...)

	return messagesAdded, nil
}
