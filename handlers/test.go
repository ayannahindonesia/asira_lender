package handlers

import (
	"asira_lender/asira"
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo"
)

func S3test(c echo.Context) error {
	defer c.Request().Body.Close()

	// body := map[string]interface{}{}
	// json.NewDecoder(c.Request().Body).Decode(&body)

	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	imagename := string(b) + strconv.FormatInt(time.Now().Unix(), 10) + ".jpeg"

	// unbased := base64.NewDecoder(base64.StdEncoding, strings.NewReader(body["image"].(string)))
	// buff := bytes.Buffer{}
	// buff.ReadFrom(unbased)

	// reader := bytes.NewReader(unbased)
	// img, _ := jpeg.Decode(reader)

	// file, err := os.Create(string(b) + time.Now().String() + ".jpeg")
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusInternalServerError, err, "noooo")
	// }
	// defer file.Close()

	// err = jpeg.Encode(file, img, nil)
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusInternalServerError, err, "noooo")
	// }

	file, _ := os.OpenFile("img/download.jpeg", os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	log.Printf("jancoook : %s", imagename)

	great, err := asira.App.S3.PutObjectJPEG(file, imagename)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "noooo")
	}

	return c.JSON(http.StatusOK, great)
}

func S3test2(c echo.Context) error {
	s, err := session.NewSession(&aws.Config{Region: aws.String("id-tbs")})
	if err != nil {
		log.Fatal(err)
	}

	// Open the file for use
	file, err := os.Open("./download.jpeg")
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "noooo")
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	greate, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("bucket-ayannah"),
		Key:                  aws.String("img/download.jpeg"),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	return c.JSON(http.StatusOK, greate)
}
