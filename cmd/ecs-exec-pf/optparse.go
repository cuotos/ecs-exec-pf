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

func parseArgs() (*Options, error) {
	return parseArgsWithFlagSet(flag.CommandLine, os.Args[1:])
}

// flag set is used so in testing we can swap it out with inputs that we control
func parseArgsWithFlagSet(flagSet *flag.FlagSet, args []string) (*Options, error) {
	opts := &Options{}

	flagSet.StringVarP(&opts.Cluster, "cluster", "c", "", "ECS cluster name")
	flagSet.StringVarP(&opts.Task, "task", "t", "", "ECS task ID.")
	flagSet.StringVarP(&opts.Container, "container", "n", "", "Container name in ECS task.")
	flagSet.IntSliceVarP(&opts.LocalPort, "local-port", "l", []int{}, "Client local port.")
	flagSet.IntSliceVarP(&opts.Port, "port", "p", []int{}, "Target remote port.")
	flagSet.BoolVarP(&opts.Debug, "debug", "d", false, "Only print the commands that would be run.")

	version := flagSet.BoolP("version", "v", false, "Print version information.")
	help := flagSet.BoolP("help", "?", false, "Print help information.")

	flagSet.SortFlags = false

	if err := flagSet.Parse(args); err != nil {
		log.Fatal(err)
	}

	if *version {
		fmt.Println(getVersion())
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if opts.Cluster == "" {
		return nil, fmt.Errorf("'--cluster' is required")
	}

	if opts.Task == "" {
		return nil, fmt.Errorf("'--task' is required")
	}

	if len(opts.Port) == 0 {
		return nil, fmt.Errorf("'--port' is required")
	}

	if len(opts.LocalPort) == 0 {
		return nil, fmt.Errorf("'--local-port' is required")
	}

	if len(opts.Port) != len(opts.LocalPort) {
		return nil, fmt.Errorf("for multiple ports, the local and remote port list should be the same length")
	}

	// make sure ports are all uint16
	for _, p := range append(opts.Port, opts.LocalPort...) {
		if p < 0 || p >= 65535 {
			return nil, fmt.Errorf("ports must be between 0 and 65536")
		}
	}

	return opts, nil
}
