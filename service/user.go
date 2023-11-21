package service

import (
	"bytes"
	"context"
	"gin-gorm-clean-template/common"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/entity"
	"gin-gorm-clean-template/helpers"
	"gin-gorm-clean-template/repository"
	"html/template"
	"os"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type UserService interface {
	RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error)
	GetAllUser(ctx context.Context) ([]entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	Verify(ctx context.Context, email string, password string) (bool, error)
	CheckUser(ctx context.Context, email string) (bool, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) error
	MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error)
	SendEmailEncrypt(ctx context.Context, UserId uuid.UUID, email string) (entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func (us *userService) RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error) {
	user := entity.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(userDTO))
	if err != nil {
		return user, err
	}
	user.SymmetricKeyAes = helpers.RandStringBytesRmndr(24)
	user.SymmetricKeyRc4 = helpers.RandStringBytesRmndr(24)
	user.SymmetricKeyDes = helpers.RandStringBytesRmndr(8)

	bitSize := 4096

	privateKey, err := helpers.GeneratePrivateKey(bitSize)
	if err != nil {
		return user, err
	}

	publicKeyBytes, err := helpers.GeneratePublicKey(&privateKey.PublicKey)
	if err != nil {
		return user, err
	}

	privateKeyBytes := helpers.EncodePrivateKeyToPEM(privateKey)

	user.PrivateKey = string(privateKeyBytes)
	user.PublicKey = string(publicKeyBytes)

	return us.userRepository.RegisterUser(ctx, user)
}

func (us *userService) GetAllUser(ctx context.Context) ([]entity.User, error) {
	return us.userRepository.GetAllUser(ctx)
}

func (us *userService) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	return us.userRepository.FindUserByEmail(ctx, email)
}

func (us *userService) Verify(ctx context.Context, email string, password string) (bool, error) {
	res, err := us.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	CheckPassword, err := helpers.CheckPassword(res.Password, []byte(password))
	if err != nil {
		return false, err
	}
	if res.Email == email && CheckPassword {
		return true, nil
	}
	return false, nil
}

func (us *userService) CheckUser(ctx context.Context, email string) (bool, error) {
	result, err := us.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if result.Email == "" {
		return false, nil
	}
	return true, nil
}

func (us *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return us.userRepository.DeleteUser(ctx, userID)
}

func (us *userService) UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) error {
	user := entity.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(userDTO))
	if err != nil {
		return err
	}
	return us.userRepository.UpdateUser(ctx, user)
}

func (us *userService) MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	return us.userRepository.FindUserByID(ctx, userID)
}

func (us *userService) SendEmailEncrypt(ctx context.Context, UserId uuid.UUID, email string) (entity.User, error) {
	user, err := us.userRepository.FindUserByID(ctx, UserId)
	if err != nil {
		return user, err
	}

	draftEmail, err := buildEmail(user.Email, email)
	if err != nil {
		return user, err
	}

	err = common.SendMail(email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return user, err
	}

	return user, nil
}

func buildEmail(requestedEmail string, requestingEmail string) (map[string]string, error) {
	readHtml, err := os.ReadFile("view/request_data.html")
	if err != nil {
		return nil, err
	}

	data := struct {
		EmailRequestingUser string
		EmailRequestedUser  string
	}{
		EmailRequestingUser: requestingEmail,
		EmailRequestedUser:  requestedEmail,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return nil, err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return nil, err
	}

	draftEmail := map[string]string{
		"subject": "Request Email",
		"body":    strMail.String(),
	}

	return draftEmail, nil
}
