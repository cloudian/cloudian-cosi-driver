package key

import (
	"encoding/json"
	"fmt"
	"io"
)

func getCanonicalUserID(body io.ReadCloser) (string, error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	type User struct {
		CanonicalUserID string `json:"canonicalUserId"`
	}

	var canonicalUser User

	err = json.Unmarshal(bytes, &canonicalUser)
	if err != nil {
		return "", fmt.Errorf("failed to parse json response: %w", err)
	}

	return canonicalUser.CanonicalUserID, nil
}
