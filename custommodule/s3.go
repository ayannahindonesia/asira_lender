package custommodule

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 main type
type S3 struct {
	S3Config   *aws.Config
	S3Uploader *s3manager.Uploader
	Bucket     string
}

// NewS3 create new S3 instance
func NewS3(accesskey string, secretkey string, host string, bucketname string, region string) (S3, error) {
	creds := credentials.NewStaticCredentials(accesskey, secretkey, "")
	_, err := creds.Get()
	if err != nil {
		return S3{}, err
	}

	session, _ := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(host),
		Credentials: creds,
	})

	_, err = session.Config.Credentials.Get()

	config := aws.NewConfig().WithEndpoint(host).WithRegion(region).WithCredentials(creds)
	x := S3{
		S3Config:   config,
		S3Uploader: s3manager.NewUploader(session),
		Bucket:     bucketname,
	}
	return x, err
}

// PutObjectJPEG uploads jpeg to s3
func (x *S3) PutObjectJPEG(file *os.File) (string, error) {
	fileinfo, _ := file.Stat()
	buffer := make([]byte, fileinfo.Size())
	file.Read(buffer)

	// Try new
	// Create an uploader with the session and default options
	response, err := x.S3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(x.Bucket),
		Key:    aws.String(fileinfo.Name()),
		Body:   bytes.NewReader(buffer),
	})

	return response.UploadID, err
}
