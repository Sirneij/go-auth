package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (app *application) uploadFileToS3(r *http.Request) (*string, error) {
	file, handler, err := r.FormFile("thumbnail")
	if err != nil {
		app.logError(r, err)
		return nil, err
	}
	defer file.Close()

	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		app.logError(r, err)
		return nil, err
	}

	// Encode bytes in base64
	s := base64.StdEncoding.EncodeToString(b)

	fileName := fmt.Sprintf("%s_%s", s, handler.Filename)
	key := fmt.Sprintf("%s%s", app.config.awsConfig.s3_key_prefix, fileName)

	_, err = app.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(app.config.awsConfig.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		app.logError(r, err)
		return nil, err
	}

	s3_url := fmt.Sprintf("%s/%s", app.config.awsConfig.BaseURL, key)

	return &s3_url, nil
}

func (app *application) deleteFileFromS3(r *http.Request) (bool, error) {
	thumbnailURL := r.FormValue("thumbnail_url")
	var objectIds []types.ObjectIdentifier

	word := "media"
	substrings := strings.Split(thumbnailURL, word)
	key := fmt.Sprintf("%s%s", word, substrings[len(substrings)-1])

	objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	_, err := app.S3Client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
		Bucket: aws.String(app.config.awsConfig.BucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		app.logError(r, err)
		return false, err
	}
	return true, nil
}
