package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	FILENAME       = "teste.txt"
	BUCKET_NAME    = "teste-golang-se"
	KEY_NAME       = "teste.txt"
	REGION         = "sa-east-1"
	LOCAL_FILENAME = "testeDownload.txt"
)

func iniciar() *s3.S3 {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(REGION),
		Credentials: credentials.NewStaticCredentials("AKIAR6BYM7NK7XFR2VMF", "ClneU5rUg+c0R5/3jyNlVpkzttNpb79JIfjazhy1", ""),
	})
	if err != nil {
		panic(err)
	}

	s3session := s3.New(sess)

	return s3session
}

func UploadFile(s3session *s3.S3) {
	file, err := os.Open(FILENAME)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Println("Uploading file")
	_, err = s3session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(KEY_NAME),
		Body:   file,
	})
	if err != nil {
		panic(err)
	}

	// resp, err := s3session.HeadObject(&s3.HeadObjectInput{
	// 	Bucket: aws.String(BUCKET_NAME),
	// 	Key:    aws.String(KEY_NAME),
	// })
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println("Uploaded file")

}

func DownloadFile(s3session *s3.S3) {
	req, _ := s3session.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(KEY_NAME),
	})
	urlStr, err := req.Presign(15 * time.Minute) // 15 minutos de validade da URL
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(LOCAL_FILENAME)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
}

func main() {

	s3session := iniciar()
	UploadFile(s3session)
	DownloadFile(s3session)
}
