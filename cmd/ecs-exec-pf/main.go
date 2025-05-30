package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	ecsexecpf "github.com/cuotos/ecs-exec-pf"
	"golang.org/x/sync/errgroup"
)

// init sets the log flags for the application.
func init() {
	log.SetFlags(0)
}

// main is the entry point for the ecs-exec-pf command-line tool.
// It parses arguments, loads AWS config, retrieves the container ID, and starts the ECS Exec session.
func main() {
	opts, err := parseArgs()
	if err != nil {
		log.Fatalf("failed to parse args: %s", err)
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load SDK config: %s", err)
	}

	containerId, err := ecsexecpf.GetContainerId(cfg, opts.Cluster, opts.Task, opts.Container)
	if err != nil {
		log.Fatalf("failed to get container ID: %s", err)
	}

	errGroup, ctx := errgroup.WithContext(context.Background())

	for i := range opts.LocalPort {
		errGroup.Go(func() error {
			return ecsexecpf.StartSession(ctx, opts.Cluster, opts.Task, containerId, opts.Port[i], opts.LocalPort[i], opts.Debug)
		})
	}
	if err := errGroup.Wait(); err != nil {
		log.Fatalf("failed to start session: %s", err)
	}
}
