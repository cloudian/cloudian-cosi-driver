package admin

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/cloudian/cosi-driver/pkg/clients/admin/api"
	"github.com/cloudian/cosi-driver/pkg/config"
	klog "k8s.io/klog/v2"
)

// Initializes a client for use with HyperStore Admin API
func NewClient(config config.Config) (*api.ClientWithResponses, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.DisableTLSCertificateChecking}, //nolint:gosec
	}
	httpClient := &http.Client{Transport: tr}

	apiClient, err := api.NewClientWithResponses(
		config.Endpoints.Admin,
		api.WithRequestEditorFn(
			func(_ context.Context, req *http.Request) error {
				req.SetBasicAuth(config.SystemAdmin.Username, config.SystemAdmin.Password)

				return nil
			},
		),
		api.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin API client: %w", err)
	}

	klog.Info("Admin Client created")

	return apiClient, nil
}
