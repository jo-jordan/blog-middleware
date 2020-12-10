package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/lzjlxebr/blog-middleware/service"
	"testing"
)

var svc *s3.S3

func TestListBuckets(t *testing.T) {
	if svc == nil {
		svc = service.CreateS3Client("us-east-2")
	}
	service.ListBuckets(svc)
}

func TestListObject(t *testing.T) {
	if svc == nil {
		svc = service.CreateS3Client("ap-northeast-2")
	}
	service.ListObject(svc, "edgeless-blog-bucket")
}

func TestUploadObject(t *testing.T) {
	sess := service.CreateSession("ap-northeast-2")

	service.UploadObject(sess, "edgeless-blog-bucket", "/Users/lzjlxebr/Documents/interview/books/PolePositionClientServer.pdf")
}

func TestBatchUploadObject(t *testing.T) {
	sess := service.CreateSession("ap-northeast-2")

	objects := []s3manager.BatchUploadObject{
		{
			Object: &s3manager.UploadInput{
				Key:    aws.String("key"),
				Bucket: aws.String("bucket"),
			},
		},
	}

	service.BatchUploadObject(sess, objects)
}
