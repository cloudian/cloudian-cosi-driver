package main

import (
	"github.com/cloudian/cosi-driver/pkg/config"
	"github.com/cloudian/cosi-driver/pkg/server"
	klog "k8s.io/klog/v2"
)

func main() {
	c := config.GetFromEnv()

	err := server.Start(c)
	if err != nil {
		klog.Fatal(err)
	}

	klog.Info("Cloudian COSI Driver Exited")
}
