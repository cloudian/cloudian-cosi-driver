package key

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCanonicalUserID(t *testing.T) {
	t.Parallel()

	responseString := `{
		"userId":"userID",
		"userType":"User",
		"groupId":"groupID",
		"active":"true",
		"canonicalUserId":"canonicalID",
		"ldapEnabled":false,
		"fileEndpoints":null
	}`

	readCloser := io.NopCloser(strings.NewReader(responseString))

	canonicalID, err := getCanonicalUserID(readCloser)
	require.NoError(t, err)
	require.Equal(t, "canonicalID", canonicalID)
}
