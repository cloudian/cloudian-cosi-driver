package iam

import (
	"context"
	"testing"

	"github.com/cloudian/cosi-driver/pkg/clients/iam"
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
			name:          "revoke_bucket_access_list_policies_fail",
			expectedError: true,
		},
		{
			name:          "revoke_bucket_access_delete_policy_fail",
			expectedError: true,
		},
		{
			name:          "revoke_bucket_access_list_access_keys_fail",
			expectedError: true,
		},
		{
			name:          "revoke_bucket_access_delete_access_key_fail",
			expectedError: true,
		},
		{
			name:          "revoke_bucket_access_delete_user_fail",
			expectedError: true,
		},
	}

	client := iam.ClientMock{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := RevokeBucketAccess(context.TODO(), client, tc.name)

			if tc.expectedError {
				require.ErrorContains(t, err, tc.name)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
