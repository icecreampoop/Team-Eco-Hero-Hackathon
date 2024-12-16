package backend

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Client s3.S3
	bucketName string
)

func init() {
	// Initialize a new AWS session for DigitalOcean Spaces
	accessKey := os.Getenv("DIGITAL_OCEAN_ACCESS_KEY")
	secretKey := os.Getenv("DIGITAL_OCEAN_SECRET_KEY")
	region := os.Getenv("DIGITAL_OCEAN_REGION")
	bucketName = os.Getenv("DIGITAL_OCEAN_BUCKET_NAME")
	endpoint := "https://" + region + ".digitaloceanspaces.com"

	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true), // Required for DigitalOcean Spaces
	}))

	s3Client = *s3.New(sess)
}

func UploadFile(fileName string, fileContent []byte) error {
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(fileContent),
		ACL:    aws.String("public-read"), // Optional: Make the file publicly accessible
	})
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
		return err
	}

	fmt.Printf("File uploaded successfully to spaces: %s\n", fileName)
	return nil
}

func DeleteFile(fileName string) error {
	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		log.Fatalf("Failed to delete file: %v", err)
		return err
	}

	fmt.Printf("File deleted from spaces successfully: %s\n", fileName)
	return nil
}