//go:build unix

package executionerServices

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
)

func killTree(cmd *exec.Cmd) {
	fmt.Println("executing from", runtime.GOOS)
	if cmd.Process == nil {
		return
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		return
	}

	syscall.Kill(-pgid, syscall.SIGTERM)
}
