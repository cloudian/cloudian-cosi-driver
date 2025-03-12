package config

import (
	"os"
	"strings"

	klog "k8s.io/klog/v2"
)

type Config struct {
	Region    string
	Endpoints struct {
		S3    string
		IAM   string
		Admin string
	}
	Credentials struct {
		AccessKey string
		SecretKey string
		Group     string
	}
	SystemAdmin struct {
		Username string
		Password string
	}
	DisableTLSCertificateChecking bool
}

func GetFromEnv() Config {
	c := Config{}
	c.Region = os.Getenv("S3_REGION")
	c.Endpoints.S3 = os.Getenv("S3_ENDPOINT")
	c.Endpoints.IAM = os.Getenv("IAM_ENDPOINT")
	c.Endpoints.Admin = os.Getenv("ADMIN_ENDPOINT")
	c.Credentials.AccessKey = os.Getenv("S3_ACCESS_KEY")
	c.Credentials.SecretKey = os.Getenv("S3_SECRET_KEY")
	c.Credentials.Group = os.Getenv("GROUP")
	c.SystemAdmin.Username = os.Getenv("ADMIN_USER")
	c.SystemAdmin.Password = os.Getenv("ADMIN_PASSWORD")
	c.DisableTLSCertificateChecking = strings.ToLower(os.Getenv("DISABLE_TLS_CERTIFICATE_CHECK")) == "true"

	klog.Info("Created config from env variables")

	return c
}
