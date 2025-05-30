package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	flag "github.com/spf13/pflag"
)


type Options struct {
	Cluster   string
	Task      string
	Container string
	Port      []int
	LocalPort []int
	Debug     bool
}

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}

	return "unknown"
}

func parseArgs() *Options {
	opts := &Options{}

	flag.StringVarP(&opts.Cluster, "cluster", "c", "", "ECS cluster name")
	flag.StringVarP(&opts.Task, "task", "t", "", "ECS task ID.")
	flag.StringVarP(&opts.Container, "container", "n", "", "Container name in ECS task.")
	flag.IntSliceVarP(&opts.LocalPort, "local-port", "l", []int{}, "Client local port.")
	flag.IntSliceVarP(&opts.Port, "port", "p", []int{}, "Target remote port.")
	flag.BoolVarP(&opts.Debug, "debug", "d", false, "Only print the commands that would be run.")

	version := flag.BoolP("version", "v", false, "Print version information.")
	help := flag.BoolP("help", "?", false, "Print help information.")

	flag.CommandLine.SortFlags = false

	flag.Parse()

	if *version {
		fmt.Println(getVersion())
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if opts.Cluster == "" {
		log.Fatal("'--cluster' is required")
	}

	if opts.Task == "" {
		log.Fatal("'--task' is required")
	}

	if len(opts.Port) == 0 {
		log.Fatal("'--port' is required")
	}

	if len(opts.LocalPort) == 0 {
		log.Fatal("'--local-port' is required")
	}

	if len(opts.Port) != len(opts.LocalPort) {
		log.Fatal("for multiple ports, the local and remote port list should be the same length")
	}

	// make sure ports are all uint16
	for _, p := range append(opts.Port, opts.LocalPort...) {
		if p < 0 || p >= 65535 {
			log.Fatal("ports must be between 0 and 65536")
		}
	}

	return opts
}
