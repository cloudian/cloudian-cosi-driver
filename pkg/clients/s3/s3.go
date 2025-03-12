package s3

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	aws "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cloudian/cosi-driver/pkg/config"
	klog "k8s.io/klog/v2"
)

// Initializes a client for use with HyperStore S3 API
func NewClient(ctx context.Context, config config.Config) (*s3.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.DisableTLSCertificateChecking}, //nolint:gosec
	}
	httpClient := &http.Client{Transport: tr}

	cfg, err := aws.LoadDefaultConfig(ctx,
		aws.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.Credentials.AccessKey,
				config.Credentials.SecretKey,
				"ignored",
			),
		),
		aws.WithHTTPClient(httpClient),
		aws.WithRegion(config.Region),
		aws.WithBaseEndpoint(config.Endpoints.S3),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load S3 config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	klog.Info("S3 Client created")

	return s3Client, nil
}
