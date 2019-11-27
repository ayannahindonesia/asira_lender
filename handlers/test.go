package handlers

import (
	"asira_lender/asira"
	"encoding/base64"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type Payload struct {
	Image string `json:"image"`
}

func S3test(c echo.Context) error {
	defer c.Request().Body.Close()

	payload := Payload{}

	payloadRules := govalidator.MapData{
		"image": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &payload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	unbased, _ := base64.StdEncoding.DecodeString(payload.Image)

	filename := randString(4) + strconv.FormatInt(time.Now().Unix(), 10)
	file, _ := os.Create(filename + ".jpeg")
	defer file.Close()
	file.Write(unbased)
	file.Sync()

	great, err := asira.App.S3.PutObjectJPEG(file)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "noooo")
	}

	return c.JSON(http.StatusOK, great)
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
