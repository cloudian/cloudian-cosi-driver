package iam

import (
	"context"
	"testing"

	"github.com/cloudian/cosi-driver/pkg/clients/iam"
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
			name:          "grant_bucket_access_create_user_fail",
			expectedError: true,
		},
		{
			name:          "grant_bucket_access_create_policy_fail",
			expectedError: true,
		},
		{
			name:          "grant_bucket_access_create_access_keys_fail",
			expectedError: true,
		},
	}

	client := iam.ClientMock{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			accessKey, secretKey, err := GrantBucketAccess(context.TODO(), client, tc.name, "mockBucketID")

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
