//go:build unix

package executionerServices

import (
	"os/exec"
	"syscall"
	"time"
)
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
		// Fallback: matar el proceso directamente
		cmd.Process.Kill()
		return
	}

	// Enviar SIGTERM a TODO el grupo de procesos
	syscall.Kill(-pgid, syscall.SIGTERM)

	// Después de 2 segundos, enviar SIGKILL si aún existe
	time.AfterFunc(2*time.Second, func() {
		// Verificar si el proceso aún existe antes de SIGKILL
		if err := syscall.Kill(-pgid, 0); err == nil {
			// El proceso aún existe, forzar terminación
			syscall.Kill(-pgid, syscall.SIGKILL)
		}
	})
}
