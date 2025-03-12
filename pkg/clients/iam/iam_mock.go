package iam

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

// Mock client

type ClientMock struct {
	BucketAccessGrantorMock
	BucketAccessRevokerMock
}

// Grant bucket access

type BucketAccessGrantorMock struct{}

func (BucketAccessGrantorMock) CreateUser(_ context.Context, params *iam.CreateUserInput, _ ...func(*iam.Options)) (*iam.CreateUserOutput, error) {
	if *params.UserName == "grant_bucket_access_create_user_fail" {
		return nil, errors.New("grant_bucket_access_create_user_fail")
	}

	return &iam.CreateUserOutput{}, nil
}

func (BucketAccessGrantorMock) PutUserPolicy(_ context.Context, params *iam.PutUserPolicyInput, _ ...func(*iam.Options)) (*iam.PutUserPolicyOutput, error) {
	if *params.UserName == "grant_bucket_access_create_policy_fail" {
		return nil, errors.New("grant_bucket_access_create_policy_fail")
	}

	return &iam.PutUserPolicyOutput{}, nil
}

func (BucketAccessGrantorMock) CreateAccessKey(_ context.Context, params *iam.CreateAccessKeyInput, _ ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error) {
	if strings.Contains(*params.UserName, "fail") {
		return nil, errors.New("grant_bucket_access_create_access_keys_fail")
	}

	accessKeyID, secretAccessKey := "mockAccessKey", "mockSecretKey"

	accessKey := &types.AccessKey{
		AccessKeyId:     &accessKeyID,
		SecretAccessKey: &secretAccessKey,
	}

	return &iam.CreateAccessKeyOutput{AccessKey: accessKey}, nil
}

// Revoke bucket access

type BucketAccessRevokerMock struct{}

func (BucketAccessRevokerMock) ListUserPolicies(_ context.Context, params *iam.ListUserPoliciesInput, _ ...func(*iam.Options)) (*iam.ListUserPoliciesOutput, error) {
	if *params.UserName == "revoke_bucket_access_list_policies_fail" {
		return nil, errors.New("revoke_bucket_access_list_policies_fail")
	}

	return &iam.ListUserPoliciesOutput{PolicyNames: []string{"mockPolicy"}}, nil
}

func (BucketAccessRevokerMock) DeleteUserPolicy(_ context.Context, params *iam.DeleteUserPolicyInput, _ ...func(*iam.Options)) (*iam.DeleteUserPolicyOutput, error) {
	if *params.UserName == "revoke_bucket_access_delete_policy_fail" {
		return nil, errors.New("revoke_bucket_access_delete_policy_fail")
	}

	return &iam.DeleteUserPolicyOutput{}, nil
}

func (BucketAccessRevokerMock) ListAccessKeys(_ context.Context, params *iam.ListAccessKeysInput, _ ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error) {
	if *params.UserName == "revoke_bucket_access_list_access_keys_fail" {
		return nil, errors.New("revoke_bucket_access_list_access_keys_fail")
	}

	accessKeyMetadata := []types.AccessKeyMetadata{{AccessKeyId: params.UserName, UserName: params.UserName}}

	return &iam.ListAccessKeysOutput{AccessKeyMetadata: accessKeyMetadata}, nil
}

func (BucketAccessRevokerMock) DeleteAccessKey(_ context.Context, params *iam.DeleteAccessKeyInput, _ ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error) {
	if *params.UserName == "revoke_bucket_access_delete_access_key_fail" {
		return nil, errors.New("revoke_bucket_access_delete_access_key_fail")
	}

	return &iam.DeleteAccessKeyOutput{}, nil
}

func (BucketAccessRevokerMock) DeleteUser(_ context.Context, params *iam.DeleteUserInput, _ ...func(*iam.Options)) (*iam.DeleteUserOutput, error) {
	if strings.Contains(*params.UserName, "fail") {
		return nil, errors.New("revoke_bucket_access_delete_user_fail")
	}

	return &iam.DeleteUserOutput{}, nil
}
