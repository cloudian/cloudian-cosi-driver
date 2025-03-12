package key

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cloudian/cosi-driver/pkg/clients/admin/api"
	klog "k8s.io/klog/v2"
)

type UserCreator interface {
	PutUser(ctx context.Context, body api.PutUserJSONRequestBody, reqEditors ...api.RequestEditorFn) (*http.Response, error)
	PutUserCredentialsWithResponse(ctx context.Context, params *api.PutUserCredentialsParams, reqEditors ...api.RequestEditorFn) (*api.PutUserCredentialsResponse, error)
}

type PolicyCreator interface {
	GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
	PutBucketPolicy(ctx context.Context, params *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)
}

// Grant bucket access using keys
func GrantBucketAccess(ctx context.Context, adminClient UserCreator, s3Client PolicyCreator, user, group, bucket string) (accessKey, secretKey string, err error) {
	klog.Infof("Creating S3 user %s, group %s", user, group)

	canonicalUserID, err := CreateUser(ctx, adminClient, user, group)
	if err != nil {
		return "", "", fmt.Errorf("failed to create user: %w", err)
	}

	klog.Infof("Creating credentials for user %s", user)

	accessKey, secretKey, err = CreateCredentials(ctx, adminClient, user, group)
	if err != nil {
		return "", "", fmt.Errorf("failed to create user credentials: %w", err)
	}

	klog.Infof("Creating policy for bucket %s to allow user %s access", bucket, user)

	if err := CreatePolicy(ctx, s3Client, canonicalUserID, bucket); err != nil {
		return "", "", fmt.Errorf("failed to create policy: %w", err)
	}

	klog.Infof("Key access granted for bucket %s, user %s", bucket, user)

	return accessKey, secretKey, nil
}

func CreateUser(ctx context.Context, client UserCreator, user, group string) (string, error) {
	reqBody := api.UserInfo{
		Active:      "true",
		GroupId:     group,
		LdapEnabled: false,
		UserId:      user,
		UserType:    "User",
	}

	response, err := client.PutUser(ctx, reqBody)
	if err != nil {
		return "", fmt.Errorf("api returned an error: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	canonicalUserID, err := getCanonicalUserID(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get canonical user ID: %w", err)
	}

	klog.Infof("User %s has canonical user ID %s", user, canonicalUserID)

	return canonicalUserID, nil
}

func CreateCredentials(ctx context.Context, client UserCreator, user, group string) (accessKey, secretKey string, err error) {
	reqBody := &api.PutUserCredentialsParams{
		UserId:  user,
		GroupId: group,
	}

	response, err := client.PutUserCredentialsWithResponse(ctx, reqBody)
	if err != nil {
		return "", "", fmt.Errorf("api returned an error: %w", err)
	}
	defer response.HTTPResponse.Body.Close()

	if response.HTTPResponse.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code: %d", response.HTTPResponse.StatusCode)
	}

	return response.JSON200.AccessKey, response.JSON200.SecretKey, nil
}

func CreatePolicy(ctx context.Context, client PolicyCreator, canonicalUser, bucket string) error {
	policyStatement := defaultPolicyStatement(canonicalUser, bucket)
	policy := defaultPolicy(canonicalUser, bucket)

	getParams := &s3.GetBucketPolicyInput{
		Bucket: &bucket,
	}

	bucketPolicy, err := client.GetBucketPolicy(ctx, getParams)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchBucketPolicy") {
			putParams := &s3.PutBucketPolicyInput{
				Bucket: &bucket,
				Policy: &policy,
			}

			if _, err := client.PutBucketPolicy(ctx, putParams); err != nil {
				return fmt.Errorf("failed to create bucket policy: %w", err)
			}

			return nil
		}

		return fmt.Errorf("failed to get bucket policy: %w", err)
	}

	klog.Infof("Bucket policy found for bucket %s, updating bucket policy", bucket)

	if err = updateBucketPolicy(ctx, client, bucketPolicy, policyStatement, canonicalUser, bucket); err != nil {
		return fmt.Errorf("failed to update bucket policy: %w", err)
	}

	return nil
}
