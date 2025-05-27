package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	ecsexecpf "github.com/winebarrel/ecs-exec-pf"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func main() {
	opts := parseArgs()

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatalf("failed to load SDK config: %s", err)
	}

	containerId, err := ecsexecpf.GetContainerId(cfg, opts.Cluster, opts.Task, opts.Container)

	if err != nil {
		log.Fatalf("failed to get container ID: %s", err)
	}

	err = ecsexecpf.StartSession(opts.Cluster, opts.Task, containerId, opts.Port, opts.LocalPort)

	if err != nil {
		log.Fatalf("failed to start session: %s", err)
	}
}
