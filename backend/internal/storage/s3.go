// Package storage wraps AWS S3 for user-uploaded images. Objects are stored
// privately (no public bucket access required); read access is granted through
// short-lived presigned GET URLs generated on demand.
package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"calorie-counter/internal/config"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Default is the process-wide storage handle, initialised once at startup by
// Init. It is nil when S3 is not configured; callers must check IsEnabled.
var Default *S3Service

// S3Service holds the S3 client and target bucket.
type S3Service struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
}

// Init builds Default from config.Values. It is a no-op (leaving Default nil)
// when the required storage variables are absent, so the rest of the app still
// boots — only the photo endpoints degrade.
//
// It works with both real AWS S3 and S3-compatible providers (e.g. MinIO on
// Railway): when S3Endpoint is set, the client targets that host and switches
// to path-style addressing, which MinIO requires.
func Init() error {
	c := config.Values
	if c.AWSS3Bucket == "" || c.AWSAccessKeyID == "" || c.AWSSecretAccessKey == "" {
		return nil
	}

	// Region is mandatory for the AWS signer but arbitrary for MinIO; default
	// to us-east-1 so a MinIO deployment doesn't need to set one.
	region := c.AWSRegion
	if region == "" {
		region = "us-east-1"
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(c.AWSAccessKeyID, c.AWSSecretAccessKey, ""),
		),
	)
	if err != nil {
		return fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if c.S3Endpoint != "" {
			o.BaseEndpoint = awssdk.String(c.S3Endpoint)
			o.UsePathStyle = true // required by MinIO and most non-AWS providers
		}
	})
	Default = &S3Service{
		client:  client,
		presign: s3.NewPresignClient(client),
		bucket:  c.AWSS3Bucket,
	}

	// Create the bucket on first boot so no manual console step is needed. This
	// is idempotent: an existing bucket we own is treated as success.
	if err := Default.ensureBucket(context.Background()); err != nil {
		return fmt.Errorf("ensure bucket %q: %w", c.AWSS3Bucket, err)
	}
	log.Printf("S3 storage ready (bucket %q)", c.AWSS3Bucket)
	return nil
}

// ensureBucket creates the configured bucket if it does not already exist.
func (s *S3Service) ensureBucket(ctx context.Context) error {
	// Fast path: the bucket already exists and we can reach it.
	if _, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: awssdk.String(s.bucket),
	}); err == nil {
		return nil
	}

	_, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: awssdk.String(s.bucket),
	})
	if err != nil {
		// Another instance (or a previous boot) may have created it first.
		var owned *s3types.BucketAlreadyOwnedByYou
		var exists *s3types.BucketAlreadyExists
		if errors.As(err, &owned) || errors.As(err, &exists) {
			return nil
		}
		return err
	}
	log.Printf("created bucket %q", s.bucket)
	return nil
}

// IsEnabled reports whether S3 is configured and ready.
func IsEnabled() bool { return Default != nil }

// Upload stores the bytes at key with the given content type.
func (s *S3Service) Upload(ctx context.Context, key, contentType string, data []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      awssdk.String(s.bucket),
		Key:         awssdk.String(key),
		Body:        bytes.NewReader(data),
		ContentType: awssdk.String(contentType),
	})
	return err
}

// PresignGet returns a temporary URL that grants read access to key.
func (s *S3Service) PresignGet(ctx context.Context, key string, expiry time.Duration) (string, error) {
	req, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: awssdk.String(s.bucket),
		Key:    awssdk.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// Delete removes the object at key.
func (s *S3Service) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: awssdk.String(s.bucket),
		Key:    awssdk.String(key),
	})
	return err
}
