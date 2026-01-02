//go:build windows

package executionerServices

import "os/exec"

func setSysProcAttr(cmd *exec.Cmd) {
	// no-op en windows
}

func killTree(cmd *exec.Cmd) {

}
