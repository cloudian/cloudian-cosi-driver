package identity

import (
	"context"

	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

type Server struct{}

func (s *Server) DriverGetInfo(_ context.Context, _ *spec.DriverGetInfoRequest) (*spec.DriverGetInfoResponse, error) {
	return &spec.DriverGetInfoResponse{Name: "cloudian-cosi-driver"}, nil
}
