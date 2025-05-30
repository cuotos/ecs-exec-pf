package main

import (
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedOpts  *Options
		expectError   bool
		fatalContains string
	}{
		{
			name: "valid single port mapping",
			args: []string{"ecs-exec-pf", "-c", "my-cluster", "-t", "task123", "-n", "container1", "-p", "80", "-l", "8080"},
			expectedOpts: &Options{
				Cluster:   "my-cluster",
				Task:      "task123",
				Container: "container1",
				Port:      []int{80},
				LocalPort: []int{8080},
				Debug:     false,
			},
			expectError: false,
		},
		{
			name: "valid multiple port mappings",
			args: []string{"ecs-exec-pf", "-c", "my-cluster", "-t", "task123", "-p", "80", "-l", "8080", "-p", "443", "-l", "8443"},
			expectedOpts: &Options{
				Cluster:   "my-cluster",
				Task:      "task123",
				Container: "",
				Port:      []int{80, 443},
				LocalPort: []int{8080, 8443},
				Debug:     false,
			},
			expectError: false,
		},
		{
			name: "debug mode",
			args: []string{"ecs-exec-pf", "-c", "my-cluster", "-t", "task123", "-p", "80", "-l", "8080", "-d"},
			expectedOpts: &Options{
				Cluster:   "my-cluster",
				Task:      "task123",
				Container: "",
				Port:      []int{80},
				LocalPort: []int{8080},
				Debug:     true,
			},
			expectError: false,
		},
		{
			name:          "missing cluster",
			args:          []string{"ecs-exec-pf", "-t", "task123", "-p", "80", "-l", "8080"},
			expectError:   true,
			fatalContains: "--cluster",
		},
		{
			name:          "missing task",
			args:          []string{"ecs-exec-pf", "-c", "my-cluster", "-p", "80", "-l", "8080"},
			expectError:   true,
			fatalContains: "--task",
		},
		{
			name:          "missing port",
			args:          []string{"ecs-exec-pf", "-c", "my-cluster", "-t", "task123", "-l", "8080"},
			expectError:   true,
			fatalContains: "--port",
		},
		{
			name:          "missing local port",
			args:          []string{"ecs-exec-pf", "-c", "my-cluster", "-t", "task123", "-p", "80"},
			expectError:   true,
			fatalContains: "--local-port",
		},
		{
			name:          "unequal port counts",
			args:          []string{"ecs-exec-pf", "-c", "my-cluster", "-t", "task123", "-p", "80", "-p", "443", "-l", "8080"},
			expectError:   true,
			fatalContains: "the local and remote port list should be the same length",
		},
		{
			name:          "invalid port number too high",
			args:          []string{"ecs-exec-pf", "-c", "my-cluster", "-t", "task123", "-p", "70000", "-l", "8080"},
			expectError:   true,
			fatalContains: "between 0 and 65536",
		},
		{
			name:          "invalid port number negative",
			args:          []string{"ecs-exec-pf", "-c", "my-cluster", "-t", "task123", "-p", "-80", "-l", "8080"},
			expectError:   true,
			fatalContains: "between 0 and 65536",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			fs := flag.NewFlagSet(tc.name, flag.PanicOnError)

			opts, err := parseArgsWithFlagSet(fs, tc.args)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			// Verify options
			assert.Equal(t, tc.expectedOpts.Cluster, opts.Cluster)
			assert.Equal(t, tc.expectedOpts.Task, opts.Task)
			assert.Equal(t, tc.expectedOpts.Container, opts.Container)
			assert.Equal(t, tc.expectedOpts.Port, opts.Port)
			assert.Equal(t, tc.expectedOpts.LocalPort, opts.LocalPort)
			assert.Equal(t, tc.expectedOpts.Debug, opts.Debug)
		})
	}
}
