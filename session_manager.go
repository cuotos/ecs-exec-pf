package ecsexecpf

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
)

func StartSession(ctx context.Context, cluster string, taskId string, containerId string, port int, localPort int, debug bool) error {
	target := fmt.Sprintf("ecs:%s_%s_%s", cluster, taskId, containerId)
	params := fmt.Sprintf(`{"portNumber":["%d"],"localPortNumber":["%d"]}`, port, localPort)

	cmdWithArgs := []string{
		"aws", "ssm", "start-session",
		"--target", target,
		"--document-name", "AWS-StartPortForwardingSession",
		"--parameters", params,
	}

	if debug {
		fmt.Println(strings.Join(cmdWithArgs, " "))
		return nil
	}
	return runCommand(ctx, cmdWithArgs)
}

func runCommand(ctx context.Context, cmdWithArgs []string) error {
	cmd := exec.CommandContext(ctx, cmdWithArgs[0], cmdWithArgs[1:]...)

	outReader, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	errReader, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig)

	go func() {
		for {
			s := <-sig
			_ = cmd.Process.Signal(s)
		}
	}()

	go func() {
		_, _ = io.Copy(os.Stdout, outReader)
		wg.Done()
	}()

	go func() {
		_, _ = io.Copy(os.Stderr, errReader)
		wg.Done()
	}()

	err = cmd.Start()

	if err != nil {
		return err
	}

	err = cmd.Wait()

	if err != nil {
		return err
	}

	wg.Wait()

	return nil
}
