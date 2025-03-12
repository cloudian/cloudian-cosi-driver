package s3

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Mock client

type ClientMock struct {
	BucketCreatorMock
	BucketDeletorMock
	PolicyModifierMock
}

// Bucket create

type BucketCreatorMock struct{}

func (BucketCreatorMock) CreateBucket(_ context.Context, params *s3.CreateBucketInput, _ ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	location := "/" + *params.Bucket

	if *params.Bucket == "create_bucket_fail" {
		return nil, errors.New("create_bucket_fail")
	}

	return &s3.CreateBucketOutput{Location: &location}, nil
}

// Bucket delete

type BucketDeletorMock struct{}

func (BucketDeletorMock) DeleteBucket(_ context.Context, params *s3.DeleteBucketInput, _ ...func(*s3.Options)) (*s3.DeleteBucketOutput, error) {
	if *params.Bucket == "delete_bucket_fail" {
		return nil, errors.New("delete_bucket_fail")
	}

	return &s3.DeleteBucketOutput{}, nil
}

// Policy modify

type PolicyModifierMock struct{}

func (PolicyModifierMock) GetBucketPolicy(_ context.Context, params *s3.GetBucketPolicyInput, _ ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
	if strings.Contains(*params.Bucket, "get_policy_fail") || strings.Contains(*params.Bucket, "bucket_access_key_fail") {
		return nil, errors.New("grant_bucket_access_get_policy_fail, revoke_bucket_access_get_policy_fail")
	}

	if strings.Contains(*params.Bucket, "policy_found") {
		mockPolicy := getMockPolicy(false)

		if strings.Contains(*params.Bucket, "delete_policy") {
			mockPolicy = getMockPolicy(true)
		}

		return &s3.GetBucketPolicyOutput{Policy: &mockPolicy}, nil
	}

	return nil, errors.New("NoSuchBucketPolicy")
}

func (PolicyModifierMock) PutBucketPolicy(_ context.Context, params *s3.PutBucketPolicyInput, _ ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
	if strings.Contains(*params.Bucket, "put_policy_fail") {
		return nil, errors.New(`
		grant_bucket_access_put_policy_fail, 
		grant_bucket_access_policy_found_put_policy_fail, 
		revoke_bucket_access_policy_found_put_policy_fail
		`)
	}

	return &s3.PutBucketPolicyOutput{}, nil
}

func (PolicyModifierMock) DeleteBucketPolicy(_ context.Context, params *s3.DeleteBucketPolicyInput, _ ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error) {
	if strings.Contains(*params.Bucket, "delete_policy_fail") {
		return nil, errors.New("revoke_bucket_access_policy_found_delete_policy_fail")
	}

	return &s3.DeleteBucketPolicyOutput{}, nil
}

func getMockPolicy(empty bool) string {
	if empty {
		return `{
			"Version": "2012-10-17",
			"Statement": []
		}`
	}

	return `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": [
					"s3:GetObject"
				],
				"Resource": [
					"arn:aws:s3:::mock-bucket/*"
				]
			},
			{
				"Effect": "Deny",
				"Principal": {
					"CanonicalUser": "mockUser"
				},
				"Action": [
					"s3:DeleteObject"
				],
				"Resource": [
					"arn:aws:s3:::mock-bucket/*"
				]
			}
		]
	}`
}
