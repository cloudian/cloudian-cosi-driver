package key

import (
	"context"
	"testing"

	"github.com/cloudian/cosi-driver/pkg/clients/admin"
	"github.com/cloudian/cosi-driver/pkg/clients/s3"
	"github.com/stretchr/testify/require"
)

func TestGrantBucketAccess(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "grant_bucket_access_success",
			expectedError: false,
		},
		{
			name:          "grant_bucket_access_policy_found_success",
			expectedError: false,
		},
		{
			name:          "grant_bucket_access_create_user_fail",
			expectedError: true,
		},
		{
			name:          "grant_bucket_access_create_credentials_fail",
			expectedError: true,
		},
		{
			name:          "grant_bucket_access_get_policy_fail",
			expectedError: true,
		},
		{
			name:          "grant_bucket_access_put_policy_fail",
			expectedError: true,
		},
		{
			name:          "grant_bucket_access_policy_found_put_policy_fail",
			expectedError: true,
		},
	}

	adminClient := admin.ClientMock{}
	s3client := s3.ClientMock{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			accessKey, secretKey, err := GrantBucketAccess(context.TODO(), adminClient, s3client, tc.name, "mockGroup", tc.name)

			if tc.expectedError {
				require.ErrorContains(t, err, tc.name)
			} else {
				require.NoError(t, err)
				require.Equal(t, accessKey, "mockAccessKey")
				require.Equal(t, secretKey, "mockSecretKey")
			}
		})
	}
}
