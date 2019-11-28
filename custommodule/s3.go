package custommodule

import (
	"asira_lender/asira"
	"bytes"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 main type
// type S3 struct {
// 	Service *s3.S3
// 	Bucket  string
// }

// NewS3 create new S3 instance
// func NewS3(accesskey string, secretkey string, host string, bucketname string, region string) (S3, error) {
// 	creds := credentials.NewStaticCredentials(accesskey, secretkey, "")
// 	_, err := creds.Get()
// 	if err != nil {
// 		return S3{}, err
// 	}

// 	config := aws.NewConfig().WithEndpoint(host).WithRegion(region).WithCredentials(creds)
// 	x := S3{
// 		Service: s3.New(session.New(), config),
// 		Bucket:  bucketname,
// 	}
// 	return x, nil
// }

// PutObjectJPEG uploads jpeg to s3
func PutObjectJPEG(file *os.File) (string, error) {
	s3Svs := s3.New(session.New(), asira.App.S3config)
	fileinfo, _ := file.Stat()
	buffer := make([]byte, fileinfo.Size())
	file.Read(buffer)
	log.Printf("buffer : %s", string(buffer))

	// Try new
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploaderWithClient(asira.App.S3.Service)
	response, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("bucket-ayannah"),
		Key:    aws.String("s0yHeS/eBvjmyshVckbNqVWnTKwhnaP6kYBpFBHk"),
		Body:   bytes.NewReader(buffer),
	})
	// ====

	// params := &s3.PutObjectInput{
	// 	Bucket:               aws.String(x.Bucket),
	// 	Key:                  aws.String(file.Name()),
	// 	ACL:                  aws.String("public-read"),
	// 	Body:                 bytes.NewReader(buffer),
	// 	ContentLength:        aws.Int64(fileinfo.Size()),
	// 	ContentType:          aws.String(http.DetectContentType(buffer)),
	// 	ContentDisposition:   aws.String("attachment"),
	// 	ServerSideEncryption: aws.String("AES256"),
	// 	StorageClass:         aws.String("INTELLIGENT_TIERING"),
	// }

	// response, err := x.Service.PutObject(params)
	if err != nil {
		return "", err
	}

	return response.String(), nil
}
