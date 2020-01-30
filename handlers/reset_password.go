package handlers

import (
	"asira_lender/asira"
	"asira_lender/models"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	// ResetRequestPayload container type
	ResetRequestPayload struct {
		Email string `json:"email"`
	}
	// ResetVerifyPayload container type
	ResetVerifyPayload struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	// Filter find interface
	Filter struct {
		Email string `json:"email"`
	}
)

func encrypt(text string, passphrase string) (string, error) {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), err
}

func decrypt(encryptedText string, passphrase string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(encryptedText)

	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("cannot decrypt")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}

// generateResetPassToken format [timestamp]|[expire_at]|[identifier]
func generateResetPassToken(identifier string) (string, error) {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	// expires in 5 mins
	expireAt := strconv.FormatInt(time.Now().Add(time.Minute*time.Duration(5)).Unix(), 10)

	// use jwt temporary
	// jwt := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	passphrase := asira.App.Config.GetString(fmt.Sprintf("%s.passphrase", asira.App.ENV))
	rawToken := now + "|" + expireAt + "|" + identifier

	return encrypt(rawToken, passphrase)
}

// UserResetPasswordRequest reset user's password
func UserResetPasswordRequest(c echo.Context) error {
	defer c.Request().Body.Close()
	var (
		resetRequestPayload ResetRequestPayload
		user                models.User
	)

	payloadRules := govalidator.MapData{
		"email": []string{"required", "email"},
	}

	validate := validateRequestPayload(c, payloadRules, &resetRequestPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	db := asira.App.DB

	err := db.Table("users").
		Select("*").
		Where("email = ?", resetRequestPayload.Email).
		First(&user).Error

	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("Email %s tidak ditemukan", resetRequestPayload.Email))
	}

	token, err := generateResetPassToken(fmt.Sprintf("%v:%v", resetRequestPayload.Email, user.ID))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Terjadi kesalahan")
	}

	message := ""
	switch c.QueryParam("system") {
	case "core":
		message = fmt.Sprintf("link reset password : %v", asira.App.Config.GetString(fmt.Sprintf("%s.core_url", asira.App.ENV))+"/ubahpassword?token="+token)
		break
	default:
		message = fmt.Sprintf("link reset password : %v", asira.App.Config.GetString(fmt.Sprintf("%s.dashboard_url", asira.App.ENV))+"/ubahpassword?token="+token)
		break
	}

	err = SendMail("Forgot Password Request", message, resetRequestPayload.Email)
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Terjadi kesalahan")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Instruksi telah dikirimkan ke email anda %v", resetRequestPayload.Email))
}

// UserResetPasswordVerify reset pass with confirmed token
func UserResetPasswordVerify(c echo.Context) error {
	defer c.Request().Body.Close()
	var (
		resetVerifyPayload ResetVerifyPayload
	)

	payloadRules := govalidator.MapData{
		"token":    []string{"required"},
		"password": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &resetVerifyPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Hambatan validasi")
	}

	d, err := decrypt(resetVerifyPayload.Token, asira.App.Config.GetString(fmt.Sprintf("%s.passphrase", asira.App.ENV)))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Terjadi kesalahan")
	}

	splits := strings.Split(d, "|")
	if len(splits) != 3 {
		return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Token tidak valid")
	}
	t, _ := strconv.ParseInt(splits[0], 10, 64)
	e, _ := strconv.ParseInt(splits[1], 10, 64)

	if time.Now().Unix() >= t && time.Now().Unix() <= e {
		splits2 := strings.Split(splits[2], ":")
		if len(splits2) != 2 {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Token tidak valid")
		}
		user := models.User{}
		err := user.FilterSearchSingle(&Filter{
			Email: splits2[0],
		})
		if err != nil {
			return returnInvalidResponse(http.StatusNotFound, err, "Token tidak valid")
		}
		user.ChangePassword(resetVerifyPayload.Password)
		user.Save()
	} else {
		return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Token tidak valid")
	}

	return c.JSON(http.StatusOK, "Password telah diganti")
}
