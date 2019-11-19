package custommodule

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 main type
type S3 struct {
	Service *s3.S3
	Bucket  string
}

// NewS3 create new S3 instance
func NewS3(accesskey string, secretkey string, host string, bucketname string, region string) (S3, error) {
	creds := credentials.NewStaticCredentials(accesskey, secretkey, "")
	_, err := creds.Get()
	if err != nil {
		return S3{}, err
	}

	config := aws.NewConfig().WithRegion(region).WithCredentials(creds)
	x := S3{
		Service: s3.New(session.New(), config),
		Bucket:  bucketname,
	}
	return x, nil
}

// PutObjectJPEG uploads jpeg to s3
func (x *S3) PutObjectJPEG(file *os.File) (string, error) {
	fileinfo, _ := file.Stat()
	params := &s3.PutObjectInput{
		Bucket:        aws.String(x.Bucket),
		Key:           aws.String(file.Name()),
		Body:          file,
		ContentLength: aws.Int64(fileinfo.Size()),
		ContentType:   aws.String("jpeg"),
	}

	response, err := x.Service.PutObject(params)
	if err != nil {
		return "", err
	}

	return response.String(), nil
}
