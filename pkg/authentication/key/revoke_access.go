package key

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cloudian/cosi-driver/pkg/clients/admin/api"
	klog "k8s.io/klog/v2"
)

type UserDeletor interface {
	GetUser(ctx context.Context, params *api.GetUserParams, reqEditors ...api.RequestEditorFn) (*http.Response, error)
	DeleteUser(ctx context.Context, params *api.DeleteUserParams, reqEditors ...api.RequestEditorFn) (*http.Response, error)
}

type PolicyDeletor interface {
	GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
	DeleteBucketPolicy(ctx context.Context, params *s3.DeleteBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.DeleteBucketPolicyOutput, error)
	PutBucketPolicy(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)
}

// Revoke bucket access using keys
func RevokeBucketAccess(ctx context.Context, adminClient UserDeletor, s3Client PolicyDeletor, user, group, bucket string) error {
	klog.Infof("Deleting S3 user %s, group %s", user, group)

	canonicalUserID, err := DeleteUser(ctx, adminClient, user, group)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	klog.Infof("Deleting bucket policy for bucket %s", bucket)

	if err := DeletePolicy(ctx, s3Client, canonicalUserID, bucket); err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	return nil
}

func DeleteUser(ctx context.Context, client UserDeletor, user, group string) (string, error) {
	response, err := client.GetUser(ctx, &api.GetUserParams{UserId: user, GroupId: group})
	if err != nil {
		return "", fmt.Errorf("get user api call returned an error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	canonicalUserID, err := getCanonicalUserID(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get canonical user ID: %w", err)
	}

	response, err = client.DeleteUser(ctx, &api.DeleteUserParams{UserId: user, GroupId: group})
	if err != nil {
		return "", fmt.Errorf("delete user api call returned an error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return canonicalUserID, nil
}

func DeletePolicy(ctx context.Context, client PolicyDeletor, canonicalUser, bucket string) error {
	getParams := &s3.GetBucketPolicyInput{
		Bucket: &bucket,
	}

	bucketPolicy, err := client.GetBucketPolicy(ctx, getParams)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchBucketPolicy") {
			return nil
		}

		return fmt.Errorf("failed to get bucket policy: %w", err)
	}

	var policyJSON map[string]interface{}

	err = json.Unmarshal([]byte(*bucketPolicy.Policy), &policyJSON)
	if err != nil {
		return fmt.Errorf("failed to parse bucket policy JSON: %w", err)
	}

	err = removePolicyStatementsForUser(policyJSON, canonicalUser)
	if err != nil {
		return fmt.Errorf("failed to remove pre-existing policy statements for canonical user '%s': %w", canonicalUser, err)
	}

	statements, ok := policyJSON["Statement"].([]interface{})
	if !ok {
		return errors.New("failed to assert policy statements as []interface{}")
	}

	if len(statements) == 0 {
		deleteParams := &s3.DeleteBucketPolicyInput{
			Bucket: &bucket,
		}

		if _, err := client.DeleteBucketPolicy(ctx, deleteParams); err != nil {
			return fmt.Errorf("failed to delete bucket policy: %w", err)
		}
	} else {
		updatedPolicy, err := json.Marshal(policyJSON)
		if err != nil {
			return fmt.Errorf("failed to marshal updated policy: %w", err)
		}

		updatedPolicyStr := string(updatedPolicy)

		putParams := &s3.PutBucketPolicyInput{
			Bucket: &bucket,
			Policy: &updatedPolicyStr,
		}

		klog.Infof("Bucket policy for bucket %s contains statements affecting other users, updating policy instead of deleting", bucket)

		if _, err := client.PutBucketPolicy(ctx, putParams); err != nil {
			return fmt.Errorf("failed to put bucket policy: %w", err)
		}
	}

	return nil
}
