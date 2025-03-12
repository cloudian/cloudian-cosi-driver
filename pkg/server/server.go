package server

import (
	"context"
	"fmt"
	"net"

	"github.com/cloudian/cosi-driver/pkg/clients/admin"
	"github.com/cloudian/cosi-driver/pkg/clients/iam"
	"github.com/cloudian/cosi-driver/pkg/clients/s3"
	"github.com/cloudian/cosi-driver/pkg/config"
	"github.com/cloudian/cosi-driver/pkg/identity"
	"github.com/cloudian/cosi-driver/pkg/provisioner"
	"google.golang.org/grpc"
	klog "k8s.io/klog/v2"
	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

func Start(cosiConfig config.Config) error {
	grpcServer := newGRPCServer()
	ctx := context.Background()

	s3Client, err := s3.NewClient(ctx, cosiConfig)
	if err != nil {
		return fmt.Errorf("failed to create s3 client: %w", err)
	}

	iamClient, err := iam.NewClient(ctx, cosiConfig)
	if err != nil {
		return fmt.Errorf("failed to create iam client: %w", err)
	}

	adminClient, err := admin.NewClient(cosiConfig)
	if err != nil {
		return fmt.Errorf("failed to create admin client: %w", err)
	}

	identity := identity.Server{}
	provisioner := provisioner.Server{
		S3Client:    s3Client,
		IAMClient:   iamClient,
		AdminClient: adminClient,
		Config:      cosiConfig,
	}

	spec.RegisterIdentityServer(grpcServer, &identity)
	spec.RegisterProvisionerServer(grpcServer, &provisioner)

	listener, err := net.Listen("unix", "/var/lib/cosi/cosi.sock")
	if err != nil {
		err := fmt.Errorf("listener error %w", err)
		klog.Error(err)

		return err
	}

	klog.Info("Starting Cloudian COSI gRPC Server")

	err = grpcServer.Serve(listener)
	if err != nil {
		err := fmt.Errorf("cloudian COSI gRPC Server exited with: %w", err)
		klog.Error(err)

		return err
	}

	klog.Info("Cloudian COSI gRPC Server exited without error")

	return nil
}

func newGRPCServer() *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(loggingInterceptor()),
	}

	return grpc.NewServer(opts...)
}

func loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		klog.InfoS("GRPC call", "method", info.FullMethod, "request", req)

		resp, err := handler(ctx, req)
		if err != nil {
			klog.Errorf("GRPC error %v", err)
		} else {
			klog.InfoS("GRPC response", "method", info.FullMethod, "response", resp)
		}

		return resp, err
	}
}
