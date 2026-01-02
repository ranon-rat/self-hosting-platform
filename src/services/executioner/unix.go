//go:build unix

package executionerServices

import (
	"os/exec"
	"syscall"
	"time"
)

func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
}

func stopCmd(cmd *exec.Cmd) {

	if cmd == nil || cmd.Process == nil {
		return
	}
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		// we can just kill the process in case of any weird behaviour
		cmd.Process.Kill()
		return
	}
	// i am checking if there is any weird behaviour
	syscall.Kill(-pgid, syscall.SIGTERM)

	time.AfterFunc(2*time.Second, func() {
		// i check if something exists
		if err := syscall.Kill(-pgid, 0); err == nil {
			// then i kill it if it still exists
			syscall.Kill(-pgid, syscall.SIGKILL)
		}
	})
}
