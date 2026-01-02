//go:build windows

package executionerServices

import "os/exec"

func setSysProcAttr(cmd *exec.Cmd) {
	// no-op en windows
}

func stopCmd(cmd *exec.Cmd) {
	cmd.Cancel()
}
