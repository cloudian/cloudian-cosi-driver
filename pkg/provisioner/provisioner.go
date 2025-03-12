package provisioner

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudian/cosi-driver/pkg/authentication/iam"
	"github.com/cloudian/cosi-driver/pkg/authentication/key"
	"github.com/cloudian/cosi-driver/pkg/bucket"
	"github.com/cloudian/cosi-driver/pkg/config"
	klog "k8s.io/klog/v2"
	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

type S3Client interface {
	bucket.BucketCreator
	bucket.BucketDeletor
	key.PolicyCreator
	key.PolicyDeletor
}

type IAMClient interface {
	iam.BucketAccessGrantor
	iam.BucketAccessRevoker
}

type AdminClient interface {
	key.UserCreator
	key.UserDeletor
}

type Server struct {
	S3Client    S3Client
	IAMClient   IAMClient
	AdminClient AdminClient
	Config      config.Config
}

func (s *Server) DriverCreateBucket(ctx context.Context, request *spec.DriverCreateBucketRequest) (*spec.DriverCreateBucketResponse, error) {
	params := request.GetParameters()
	storagePolicyID := params["storagePolicyID"]

	klog.Infof("Creating bucket %s", request.GetName())

	if err := bucket.CreateBucket(ctx, s.S3Client, request.GetName(), storagePolicyID); err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return &spec.DriverCreateBucketResponse{
		BucketId: request.GetName(),
		BucketInfo: &spec.Protocol{
			Type: &spec.Protocol_S3{
				S3: &spec.S3{
					Region: s.Config.Region,
				},
			},
		},
	}, nil
}

func (s *Server) DriverDeleteBucket(ctx context.Context, request *spec.DriverDeleteBucketRequest) (resp *spec.DriverDeleteBucketResponse, err error) {
	klog.Infof("Deleting bucket %s", request.GetBucketId())

	if err := bucket.DeleteBucket(ctx, s.S3Client, request.GetBucketId()); err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return resp, err
}

func (s *Server) DriverGrantBucketAccess(ctx context.Context, request *spec.DriverGrantBucketAccessRequest) (*spec.DriverGrantBucketAccessResponse, error) {
	var accessKey, secretKey, accountID string

	var err error

	if request.GetAuthenticationType() == spec.AuthenticationType_IAM {
		klog.Infof("Granting IAM access to bucket %s", request.GetBucketId())

		accountID = "iam_" + request.GetName()
		accessKey, secretKey, err = iam.GrantBucketAccess(ctx, s.IAMClient, accountID, request.GetBucketId())
	} else {
		klog.Infof("Granting Key access to bucket %s", request.GetBucketId())

		accountID = "key_" + request.GetName()
		accessKey, secretKey, err = key.GrantBucketAccess(ctx, s.AdminClient, s.S3Client, accountID, s.Config.Credentials.Group, request.GetBucketId())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to grant bucket access: %w", err)
	}

	credentialsMap := map[string]*spec.CredentialDetails{
		"s3": {
			Secrets: map[string]string{
				"authenticationType": request.GetAuthenticationType().String(),
				"endpoint":           s.Config.Endpoints.S3,
				"region":             s.Config.Region,
				"accessKeyID":        accessKey,
				"accessSecretKey":    secretKey,
			},
		},
	}

	return &spec.DriverGrantBucketAccessResponse{
		AccountId:   accountID,
		Credentials: credentialsMap,
	}, nil
}

func (s *Server) DriverRevokeBucketAccess(ctx context.Context, request *spec.DriverRevokeBucketAccessRequest) (*spec.DriverRevokeBucketAccessResponse, error) {
	var err error

	if strings.HasPrefix(request.GetAccountId(), "iam") {
		klog.Infof("Revoking IAM access to bucket %s", request.GetBucketId())

		err = iam.RevokeBucketAccess(ctx, s.IAMClient, request.GetAccountId())

		klog.Infof("IAM access revoked for bucket %s", request.GetBucketId())
	} else {
		klog.Infof("Revoking Key access to bucket %s", request.GetBucketId())

		err = key.RevokeBucketAccess(ctx, s.AdminClient, s.S3Client, request.GetAccountId(), s.Config.Credentials.Group, request.GetBucketId())

		klog.Infof("Key access revoked for bucket %s", request.GetBucketId())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to revoke bucket access: %w", err)
	}

	return &spec.DriverRevokeBucketAccessResponse{}, nil
}
