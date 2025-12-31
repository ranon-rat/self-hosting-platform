package executioner

import (
	"context"
	"os/exec"
)

type RunningProject struct {
	Cmd    *exec.Cmd
	Cancel context.CancelFunc
}
