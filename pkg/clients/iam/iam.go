package iam

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	aws "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/cloudian/cosi-driver/pkg/config"
	klog "k8s.io/klog/v2"
)

// Initializes a client for use with HyperStore IAM API
func NewClient(ctx context.Context, config config.Config) (*iam.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.DisableTLSCertificateChecking}, //nolint:gosec
	}
	httpClient := &http.Client{Transport: tr}

	cfg, err := aws.LoadDefaultConfig(ctx,
		aws.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.Credentials.AccessKey, config.Credentials.SecretKey, "ignored")),
		aws.WithHTTPClient(httpClient),
		aws.WithRegion(config.Region),
		aws.WithBaseEndpoint(config.Endpoints.IAM),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load IAM config: %w", err)
	}

	iamClient := iam.NewFromConfig(cfg)

	klog.Info("IAM Client created")

	return iamClient, nil
}
