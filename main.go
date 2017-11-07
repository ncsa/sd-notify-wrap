package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/coreos/go-systemd/daemon"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage %s cmd [args]\n", os.Args[0])
		os.Exit(1)
	}
	prog := os.Args[1]
	args := os.Args[2:]

	cmd := exec.Command(prog, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start %s: %v\n", prog, err)
		os.Exit(2)
	}
	fmt.Fprintf(os.Stderr, "Pid: %d\n", cmd.Process.Pid)

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Fprintf(os.Stderr, "Survived for 2 seconds!\n")
		daemon.SdNotify(false, "READY=1") //
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run %s: %v\n", prog, err)
		os.Exit(2)
	}
	os.Exit(0)
}
