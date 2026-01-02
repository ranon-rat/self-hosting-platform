//go:build !linux && !darwin && !freebsd && !openbsd && !netbsd

package executionerServices

import "os/exec"

func setSysProcAttr(cmd *exec.Cmd) {
	// no-op en windows
}
