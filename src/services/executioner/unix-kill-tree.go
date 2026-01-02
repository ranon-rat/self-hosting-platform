//go:build unix

package executionerServices

import (
	"os/exec"
	"syscall"
	"time"
)

func stopCmd(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		return
	}

	syscall.Kill(-pgid, syscall.SIGTERM)

	time.AfterFunc(2*time.Second, func() {
		syscall.Kill(-pgid, syscall.SIGKILL)
	})
}
