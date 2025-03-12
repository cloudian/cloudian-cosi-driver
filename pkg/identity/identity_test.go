package identity

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

func TestDriverGetInfo(t *testing.T) {
	t.Parallel()

	server := Server{}
	expectedName := "cloudian-cosi-driver"

	response, err := server.DriverGetInfo(context.TODO(), &spec.DriverGetInfoRequest{})
	actualName := response.GetName()

	require.Equal(t, expectedName, actualName, fmt.Sprintf("expected driver name '%s', got '%s'", expectedName, actualName))
	require.NoError(t, err)
}
