package key

import (
	"context"
	"testing"

	"github.com/cloudian/cosi-driver/pkg/clients/admin"
	"github.com/cloudian/cosi-driver/pkg/clients/s3"
	"github.com/stretchr/testify/require"
)

func TestRevokeBucketAccess(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "revoke_bucket_access_success",
			expectedError: false,
		},
		{
			name:          "revoke_bucket_access_policy_found_put_policy_success",
			expectedError: false,
		},
		{
			name:          "revoke_bucket_access_policy_found_delete_policy_success",
			expectedError: false,
		},
		{
			name:          "revoke_bucket_access_get_user_fail",
			expectedError: true,
		},
		{
			name:          "revoke_bucket_access_delete_user_fail",
			expectedError: true,
		},
		{
			name:          "revoke_bucket_access_get_policy_fail",
			expectedError: true,
		},
		{
			name:          "revoke_bucket_access_policy_found_put_policy_fail",
			expectedError: true,
		},
		{
			name:          "revoke_bucket_access_policy_found_delete_policy_fail",
			expectedError: true,
		},
	}

	adminClient := admin.ClientMock{}
	s3client := s3.ClientMock{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := RevokeBucketAccess(context.TODO(), adminClient, s3client, tc.name, "mockGroup", tc.name)

			if tc.expectedError {
				require.ErrorContains(t, err, tc.name)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
