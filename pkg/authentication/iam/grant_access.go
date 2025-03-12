package iam

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	klog "k8s.io/klog/v2"
)

type BucketAccessGrantor interface {
	CreateUser(ctx context.Context, params *iam.CreateUserInput, optFns ...func(*iam.Options)) (*iam.CreateUserOutput, error)
	PutUserPolicy(ctx context.Context, params *iam.PutUserPolicyInput, optFns ...func(*iam.Options)) (*iam.PutUserPolicyOutput, error)
	CreateAccessKey(ctx context.Context, params *iam.CreateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error)
}

// Grant bucket access using IAM
func GrantBucketAccess(ctx context.Context, client BucketAccessGrantor, name, bucketID string) (accessKey, secretKey string, err error) {
	klog.Infof("Creating IAM user %s", name)

	if _, err := client.CreateUser(ctx, &iam.CreateUserInput{UserName: &name}); err != nil {
		return "", "", fmt.Errorf("failed to create iam user: %w", err)
	}

	klog.Infof("Creating IAM policy for bucket %s, user %s", bucketID, name)

	if err := CreatePolicy(ctx, client, name, bucketID); err != nil {
		return "", "", fmt.Errorf("failed to create iam policy: %w", err)
	}

	klog.Infof("Creating access key for IAM user %s", name)

	output, err := client.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{UserName: &name})
	if err != nil {
		return "", "", fmt.Errorf("failed to create access key for iam user: %w", err)
	}

	klog.Infof("IAM access granted for bucket %s, user %s", bucketID, name)

	return *output.AccessKey.AccessKeyId, *output.AccessKey.SecretAccessKey, nil
}

// Create policy via IAM API
func CreatePolicy(ctx context.Context, client BucketAccessGrantor, userName, bucketName string) error {
	policy := defaultPolicy(bucketName)

	params := &iam.PutUserPolicyInput{
		UserName:       &userName,
		PolicyName:     &userName,
		PolicyDocument: &policy,
	}

	if _, err := client.PutUserPolicy(ctx, params); err != nil {
		return fmt.Errorf("iam client failed to put user policy: %w", err)
	}

	return nil
}
