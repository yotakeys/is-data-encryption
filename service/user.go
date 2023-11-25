package service

import (
	"bytes"
	"context"
	"errors"
	"gin-gorm-clean-template/common"
	"gin-gorm-clean-template/dto"
	"gin-gorm-clean-template/encrypt"
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
	SendEmailResponse(ctx context.Context, email string, response string) error
	AsymmetricEncrypt(ctx context.Context, requestedUserID uuid.UUID, requestingUserEmail string, response string) (entity.User, error)
	AsymmetricDecrypt(ctx context.Context, userID uuid.UUID, requestingUserEmail string) ([]dto.DecryptRSAResponseDTO, error)
	GetListRequestingUser(ctx context.Context, userID uuid.UUID) ([]entity.User, error)
}

type userService struct {
	userRepository    repository.UserRepository
	encryptRepository repository.EncryptRepository
}

func NewUserService(ur repository.UserRepository, er repository.EncryptRepository) UserService {
	return &userService{
		userRepository:    ur,
		encryptRepository: er,
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
	requestedUser, err := us.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return entity.User{}, errors.New("email user yang di request tidak terdaftar")
	}
	asymmetrics, err := us.userRepository.FindAsymmetricByUserID(ctx, UserId, requestedUser.ID)
	if err != nil || len(asymmetrics) > 0 {
		return entity.User{}, errors.New("user sudah melakukan request enkripsi")
	}

	encrypt, err := us.encryptRepository.GetFirstAESEncrpytedData(ctx, requestedUser.ID)
	if err != nil {
		return entity.User{}, err
	}
	if encrypt == (entity.Encrypt{}) {
		return entity.User{}, errors.New("user belum melakukan enkripsi")
	}

	user, err := us.userRepository.FindUserByID(ctx, UserId)
	if err != nil {
		return user, err
	}

	asymmetric := entity.Asymmetric{
		RequestedUserID:  requestedUser.ID,
		RequestingUserID: UserId,
	}

	_, err = us.userRepository.CreateAsymmetric(ctx, asymmetric)
	if err != nil {
		return user, err
	}

	draftEmail, err := buildEmailEncrypt(user.Email, email)
	if err != nil {
		return user, err
	}

	err = common.SendMail(email, draftEmail["subject"], draftEmail["body"])
	if err != nil {
		return user, err
	}

	return user, nil
}

func (us *userService) SendEmailResponse(ctx context.Context, email string, response string) error {
	if response != "accept" && response != "decline" {
		return errors.New("response must be accept or decline")
	}
	emailResponse, err := os.ReadFile("view/accept_response.html")
	if err != nil {
		return err
	}
	if response == "decline" {
		emailResponse, err = os.ReadFile("view/decline_response.html")
		if err != nil {
			return err
		}
	}

	err = common.SendMail(email, "Response Email Data", string(emailResponse))
	if err != nil {
		return err
	}

	return nil
}

func buildEmailEncrypt(requestedEmail string, requestingEmail string) (map[string]string, error) {
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

func (us *userService) AsymmetricEncrypt(ctx context.Context, requestedUserID uuid.UUID, requestingUserEmail string, response string) (entity.User, error) {
	us.SendEmailResponse(ctx, requestingUserEmail, response)

	requestedUser, err := us.userRepository.FindUserByID(ctx, requestedUserID)
	if err != nil {
		return requestedUser, errors.New("user tidak ditemukan")
	}
	requestingUser, err := us.userRepository.FindUserByEmail(ctx, requestingUserEmail)
	if err != nil {
		return requestingUser, errors.New("email user yang me request tidak terdaftar")
	}

	if response == "decline" {
		err := us.userRepository.DeleteRequestingUser(ctx, requestedUser.ID, requestingUser.ID)
		if err != nil {
			return requestedUser, err
		}

		return requestedUser, nil
	}

	dataRequested, err := us.encryptRepository.GetAllEncrypt(ctx, requestedUser.ID)
	if err != nil || len(dataRequested) == 0 {
		return requestedUser, errors.New("user yang di request tidak memiliki data")
	}

	encryptedData, err := us.encryptRepository.GetFirstAESEncrpytedData(ctx, requestedUser.ID)
	if err != nil {
		return requestedUser, err
	}

	decrypt_name, err := encrypt.AESDecrypt(encryptedData.Name)
	if err != nil {
		return requestedUser, err
	}
	decrypt_phone, err := encrypt.AESDecrypt(encryptedData.PhoneNumber)
	if err != nil {
		return requestedUser, err
	}
	decrypt_idcard, err := encrypt.AESDecrypt(encryptedData.IDCardUrl)
	if err != nil {
		return requestedUser, err
	}
	decrypt_cv, err := encrypt.AESDecrypt(encryptedData.CVUrl)
	if err != nil {
		return requestedUser, err
	}
	decrypt_video, err := encrypt.AESDecrypt(encryptedData.VideoUrl)
	if err != nil {
		return requestedUser, err
	}

	asymmetrics, err := us.userRepository.FindAsymmetricByUserID(ctx, requestingUser.ID, requestedUser.ID)
	if err != nil || len(asymmetrics) == 0 {
		return requestedUser, errors.New("user belum melakukan request enkripsi")
	}

	asymmetric := entity.Asymmetric{
		ID:               asymmetrics[0].ID,
		RequestingUserID: requestingUser.ID,
		RequestedUserID:  requestedUser.ID,
		Name:             encrypt.EncryptRSA(decrypt_name, requestingUser.PublicKey),
		PhoneNumber:      encrypt.EncryptRSA(decrypt_phone, requestingUser.PublicKey),
		IDCardUrl:        encrypt.EncryptRSA(decrypt_idcard, requestingUser.PublicKey),
		CVUrl:            encrypt.EncryptRSA(decrypt_cv, requestingUser.PublicKey),
		VideoUrl:         encrypt.EncryptRSA(decrypt_video, requestingUser.PublicKey),
	}

	_, err = us.userRepository.UpdateAsymmetric(ctx, asymmetric)
	if err != nil {
		return requestedUser, err
	}

	return requestedUser, nil
}

func (us *userService) AsymmetricDecrypt(ctx context.Context, userID uuid.UUID, requestingUserEmail string) ([]dto.DecryptRSAResponseDTO, error) {
	requestedUser, err := us.userRepository.FindUserByEmail(ctx, requestingUserEmail)
	if err != nil {
		return []dto.DecryptRSAResponseDTO{}, errors.New("email user yang di request tidak terdaftar")
	}

	requestingUser, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return []dto.DecryptRSAResponseDTO{}, err
	}

	asymmetrics, err := us.userRepository.FindAsymmetricByUserID(ctx, userID, requestedUser.ID)
	if err != nil || len(asymmetrics) == 0 {
		return []dto.DecryptRSAResponseDTO{}, errors.New("akses data ke user tidak ada")
	}

	var decryptResponse []dto.DecryptRSAResponseDTO

	for _, asymmetric := range asymmetrics {
		decryptResponse = append(decryptResponse, dto.DecryptRSAResponseDTO{
			Name:        encrypt.DecryptRSA(asymmetric.Name, requestingUser.PrivateKey),
			PhoneNumber: encrypt.DecryptRSA(asymmetric.PhoneNumber, requestingUser.PrivateKey),
			IDCard:      encrypt.DecryptRSA(asymmetric.IDCardUrl, requestingUser.PrivateKey),
			CV:          encrypt.DecryptRSA(asymmetric.CVUrl, requestingUser.PrivateKey),
			Video:       encrypt.DecryptRSA(asymmetric.VideoUrl, requestingUser.PrivateKey),
		})
	}

	return decryptResponse, nil
}

func (us *userService) GetListRequestingUser(ctx context.Context, userID uuid.UUID) ([]entity.User, error) {
	requestingUsers, err := us.userRepository.FindRequestingUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	listUsers := []entity.User{}
	for _, requestingUser := range requestingUsers {
		user, err := us.userRepository.FindUserByID(ctx, requestingUser.RequestingUserID)
		if err != nil {
			return nil, err
		}
		listUsers = append(listUsers, user)
	}

	return listUsers, nil
}
