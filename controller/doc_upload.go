package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/bkrmdahal/go_restAPI/utils"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {

	//r is an *http.Request object
	file, h, err := r.FormFile("data") // or whatever the form key is
	if err != nil {
		// handle error
	}

	//select Region to use.
	conf := aws.Config{Region: aws.String("us-east-1")}
	sess := session.New(&conf)
	svc := s3manager.NewUploader(sess)

	utils.Log.Info("Uploading file to S3...")
	_, errs := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String("testing-go"),
		Key:    aws.String(filepath.Base(h.Filename)),
		Body:   file,
	})
	if errs != nil {
		fmt.Println("error", errs)
		os.Exit(1)
	}

	utils.Log.Info("Successfully uploaded .", h.Filename)
	w.Write([]byte("upload success"))
}
