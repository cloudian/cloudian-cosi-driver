package provisioner

import (
	"context"
	"testing"

	"github.com/cloudian/cosi-driver/pkg/clients/admin"
	"github.com/cloudian/cosi-driver/pkg/clients/iam"
	"github.com/cloudian/cosi-driver/pkg/clients/s3"
	"github.com/cloudian/cosi-driver/pkg/config"
	"github.com/stretchr/testify/require"
	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

func TestDriverCreateBucket(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "create_bucket_success",
			expectedError: false,
		},
		{
			name:          "create_bucket_fail",
			expectedError: true,
		},
	}

	client := s3.ClientMock{}
	region := "mock_region"
	cfg := config.Config{Region: region}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := Server{
				S3Client: client,
				Config:   cfg,
			}

			response, err := server.DriverCreateBucket(context.TODO(), &spec.DriverCreateBucketRequest{Name: tc.name})

			if tc.expectedError {
				require.Nil(t, response)
				require.Error(t, err)
			} else {
				require.Equal(t, response.GetBucketId(), tc.name)
				require.Equal(t, response.GetBucketInfo().GetS3().GetRegion(), region)
				require.NoError(t, err)
			}
		})
	}
}

func TestDriverDeleteBucket(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "delete_bucket_success",
			expectedError: false,
		},
		{
			name:          "delete_bucket_fail",
			expectedError: true,
		},
	}

	client := s3.ClientMock{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := Server{
				S3Client: client,
			}

			response, err := server.DriverDeleteBucket(context.TODO(), &spec.DriverDeleteBucketRequest{BucketId: tc.name})

			require.Nil(t, response)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDriverGrantBucketAccess(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		authenticationType spec.AuthenticationType
		expectedError      bool
	}{
		{
			name:               "grant_bucket_access_iam_success",
			authenticationType: spec.AuthenticationType_IAM,
			expectedError:      false,
		},
		{
			name:               "grant_bucket_access_iam_fail",
			authenticationType: spec.AuthenticationType_IAM,
			expectedError:      true,
		},
		{
			name:               "grant_bucket_access_key_success",
			authenticationType: spec.AuthenticationType_Key,
			expectedError:      false,
		},
		{
			name:               "grant_bucket_access_key_fail",
			authenticationType: spec.AuthenticationType_Key,
			expectedError:      true,
		},
	}

	cfg := config.Config{
		Region: "mock-region",
		Endpoints: struct {
			S3    string
			IAM   string
			Admin string
		}{
			S3: "mockS3Endpoint",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := Server{
				IAMClient:   iam.ClientMock{},
				S3Client:    s3.ClientMock{},
				AdminClient: admin.ClientMock{},
				Config:      cfg,
			}

			request := &spec.DriverGrantBucketAccessRequest{
				BucketId:           tc.name,
				Name:               tc.name,
				AuthenticationType: tc.authenticationType,
			}

			response, err := server.DriverGrantBucketAccess(context.TODO(), request)

			prefix := "key_"
			if tc.authenticationType == spec.AuthenticationType_IAM {
				prefix = "iam_"
			}

			expectedResponse := spec.DriverGrantBucketAccessResponse{
				AccountId: prefix + tc.name,
				Credentials: map[string]*spec.CredentialDetails{
					"s3": {
						Secrets: map[string]string{
							"authenticationType": tc.authenticationType.String(),
							"endpoint":           cfg.Endpoints.S3,
							"region":             cfg.Region,
							"accessKeyID":        "mockAccessKey",
							"accessSecretKey":    "mockSecretKey",
						},
					},
				},
			}

			if tc.expectedError {
				require.Nil(t, response)
				require.Error(t, err)
			} else {
				require.Equal(t, expectedResponse.GetAccountId(), response.GetAccountId())
				require.Equal(t, expectedResponse.GetCredentials(), response.GetCredentials())
				require.NoError(t, err)
			}
		})
	}
}

func TestDriverRevokeBucketAccess(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		authenticationType spec.AuthenticationType
		expectedError      bool
	}{
		{
			name:               "revoke_bucket_access_iam_success",
			authenticationType: spec.AuthenticationType_IAM,
			expectedError:      false,
		},
		{
			name:               "revoke_bucket_access_iam_fail",
			authenticationType: spec.AuthenticationType_IAM,
			expectedError:      true,
		},
		{
			name:               "revoke_bucket_access_key_success",
			authenticationType: spec.AuthenticationType_Key,
			expectedError:      false,
		},
		{
			name:               "revoke_bucket_access_key_fail",
			authenticationType: spec.AuthenticationType_Key,
			expectedError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := Server{
				IAMClient:   iam.ClientMock{},
				S3Client:    s3.ClientMock{},
				AdminClient: admin.ClientMock{},
			}

			prefix := "key_"
			if tc.authenticationType == spec.AuthenticationType_IAM {
				prefix = "iam_"
			}

			request := &spec.DriverRevokeBucketAccessRequest{
				BucketId:  tc.name,
				AccountId: prefix + tc.name,
			}

			response, err := server.DriverRevokeBucketAccess(context.TODO(), request)

			if tc.expectedError {
				require.Nil(t, response)
				require.Error(t, err)
			} else {
				require.NotNil(t, response)
				require.NoError(t, err)
			}
		})
	}
}
