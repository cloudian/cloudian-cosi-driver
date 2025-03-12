package bucket

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	klog "k8s.io/klog/v2"
)

type BucketCreator interface { //nolint:revive
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
}

// Create bucket via S3 API
func CreateBucket(ctx context.Context, client BucketCreator, name, storagePolicyID string) error {
	params := &s3.CreateBucketInput{
		Bucket: &name,
	}

	options := func(o *s3.Options) {
		o.HTTPClient = &http.Client{
			Transport: roundTripper{http.DefaultTransport, storagePolicyID},
		}
	}

	if _, err := client.CreateBucket(ctx, params, options); err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", name, err)
	}

	klog.Infof("Bucket %s created", name)

	return nil
}

type roundTripper struct {
	rt              http.RoundTripper
	storagePolicyID string
}

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("x-gmt-policyid", r.storagePolicyID) //nolint:canonicalheader

	return r.rt.RoundTrip(req) //nolint:wrapcheck
}
