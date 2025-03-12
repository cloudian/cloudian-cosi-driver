package key

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemovePolicyStatementsForUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		expectedOutput map[string]interface{}
	}{
		{
			name: "statement_removed",
			expectedOutput: map[string]interface{}{
				"Statement": []interface{}{},
				"Version":   "2012-10-17",
			},
		},
		{
			name: "statement_not_removed",
			expectedOutput: map[string]interface{}{
				"Statement": []interface{}{
					map[string]interface{}{
						"Action": []interface{}{
							"s3:GetObject",
							"s3:GetObjectVersion",
							"s3:PutObject",
							"s3:DeleteObject",
							"s3:DeleteObjectVersion",
							"s3:ListBucket",
						},
						"Effect":    "Allow",
						"Principal": map[string]interface{}{"CanonicalUser": "statement_not_removed"},
						"Resource": []interface{}{
							"arn:aws:s3:::bucket",
							"arn:aws:s3:::bucket/*",
						},
					},
				},
				"Version": "2012-10-17",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var policyJSON map[string]interface{}

			err := json.Unmarshal([]byte(defaultPolicy(tc.name, "bucket")), &policyJSON)
			require.NoError(t, err)

			err = removePolicyStatementsForUser(policyJSON, "statement_removed")
			require.NoError(t, err)
			require.Equal(t, tc.expectedOutput, policyJSON)
		})
	}
}

func TestAddStatementToPolicy(t *testing.T) {
	t.Parallel()

	policyJSON := map[string]interface{}{
		"Statement": []interface{}{},
		"Version":   "2012-10-17",
	}

	expectedOutput := map[string]interface{}{
		"Statement": []interface{}{
			map[string]interface{}{
				"Action": []interface{}{
					"s3:GetObject",
					"s3:GetObjectVersion",
					"s3:PutObject",
					"s3:DeleteObject",
					"s3:DeleteObjectVersion",
					"s3:ListBucket",
				},
				"Effect":    "Allow",
				"Principal": map[string]interface{}{"CanonicalUser": "user"},
				"Resource": []interface{}{
					"arn:aws:s3:::bucket",
					"arn:aws:s3:::bucket/*",
				},
			},
		},
		"Version": "2012-10-17",
	}

	policyStatement := defaultPolicyStatement("user", "bucket")

	err := addStatementToPolicy(policyJSON, policyStatement)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, policyJSON)
}
