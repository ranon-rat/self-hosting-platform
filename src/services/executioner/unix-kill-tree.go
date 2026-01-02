//go:build unix

package executionerServices

import (
	"os/exec"
	"syscall"
)

func killTree(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		return
	}

	syscall.Kill(-pgid, syscall.SIGTERM)
}
