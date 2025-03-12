package bucket

import (
	"context"
	"testing"

	"github.com/cloudian/cosi-driver/pkg/clients/s3"
	"github.com/stretchr/testify/require"
)

func TestDeleteBucket(t *testing.T) {
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

			err := DeleteBucket(context.TODO(), client, tc.name)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
