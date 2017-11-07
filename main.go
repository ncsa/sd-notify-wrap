package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
)

var readydelay time.Duration

func init() {
	flag.DurationVar(&readydelay, "delay", 2*time.Second, "How long before notifying systemd that the proccess is ready")
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage %s cmd [args]\n", os.Args[0])
		os.Exit(1)
	}
	prog := flag.Args()[0]
	args := flag.Args()[1:]

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
		time.Sleep(readydelay)
		fmt.Fprintf(os.Stderr, "Survived for %s!\n", readydelay)
		daemon.SdNotify(false, "READY=1") //
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run %s: %v\n", prog, err)
		code := err.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus()
		os.Exit(code)
	}
	os.Exit(0)
}
