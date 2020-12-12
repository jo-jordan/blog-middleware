package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/lzjlxebr/blog-middleware/common"
	"log"
	"os"
	"time"
)

func CreateSession(region string) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(common.ACID, common.AS, ""),
		HTTPClient: common.NewHTTPClientWithSettings(common.HTTPClientSettings{
			Connect:          5 * time.Second,
			ExpectContinue:   1 * time.Second,
			IdleConn:         90 * time.Second,
			ConnKeepAlive:    30 * time.Second,
			MaxAllIdleConns:  100,
			MaxHostIdleConns: 10,
			ResponseHeader:   5 * time.Second,
			TLSHandshake:     5 * time.Second,
		}),
	})
	common.ErrorBus(err)

	return sess
}

func CreateS3Client(region string) *s3.S3 {
	svc := s3.New(CreateSession(region))

	return svc
}

func BatchDeleteObject(region string, objects []s3manager.BatchDeleteObject) {
	sess := CreateSession(region)
	deleter := s3manager.NewBatchDelete(sess)

	iter := &s3manager.DeleteObjectsIterator{Objects: objects}
	err := deleter.Delete(aws.BackgroundContext(), iter)
	common.ErrorBus(err)
}

func ListBuckets(svc *s3.S3) []*s3.Bucket {
	result, err := svc.ListBuckets(nil)
	common.ErrorBus(err)
	log.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
	return result.Buckets
}

func ListObject(svc *s3.S3, bucket string) []string {

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	common.ErrorBus(err)

	var buckets []string
	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
		buckets = append(buckets, *item.Key)
	}

	return buckets
}

func UploadObject(session *session.Session, bucket string, filename string) {
	file, err := os.Open(filename)
	common.ErrorBus(err)
	defer file.Close()

	uploader := s3manager.NewUploader(session)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	common.ErrorBus(err)

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}

func BatchUploadObject(region string, objects []s3manager.BatchUploadObject) {
	sess := CreateSession(region)
	uploader := s3manager.NewUploader(sess)
	iter := &s3manager.UploadObjectsIterator{Objects: objects}
	err := uploader.UploadWithIterator(aws.BackgroundContext(), iter)
	common.ErrorBus(err)
}
