package custommodule

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	S3Client *s3.S3
	Bucket   string
}

// NewS3 create new S3 instance
func NewS3(accesskey string, secretkey string, host string, bucketname string, region string) (S3, error) {
	creds := credentials.NewStaticCredentials(accesskey, secretkey, "")
	_, err := creds.Get()
	if err != nil {
		return S3{}, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	session, _ := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(host),
		Credentials:      creds,
		S3ForcePathStyle: aws.Bool(true),
		HTTPClient:       client,
	})

	s3Client := s3.New(session)
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketname),
	})
	if err != nil {
		return S3{}, err
	}

	return S3{S3Client: s3Client, Bucket: bucketname}, err
}

// UploadJPEG uploads jpeg to s3
func (x *S3) UploadJPEG(file *os.File) (string, error) {
	fileinfo, _ := file.Stat()
	buffer := make([]byte, fileinfo.Size())
	file.Read(buffer)

	// Try new
	// Create an uploader with the session and default options
	response, err := x.S3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(x.Bucket),
		Key:    aws.String(fileinfo.Name()),
		Body:   bytes.NewReader(buffer),
	})

	return fmt.Sprintf("%v", response), err
}
