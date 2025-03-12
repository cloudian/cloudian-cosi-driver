package bucket

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	klog "k8s.io/klog/v2"
)

type BucketDeletor interface { //nolint:revive
	DeleteBucket(ctx context.Context, params *s3.DeleteBucketInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
}

// Delete bucket via S3 API
func DeleteBucket(ctx context.Context, client BucketDeletor, name string) error {
	params := &s3.DeleteBucketInput{
		Bucket: &name,
	}

	if _, err := client.DeleteBucket(ctx, params); err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}

	klog.Infof("Bucket %s deleted", name)

	return nil
}
