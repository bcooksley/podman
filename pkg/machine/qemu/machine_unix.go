//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

package qemu

import (
	"bytes"
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func isProcessAlive(pid int) bool {
	err := unix.Kill(pid, syscall.Signal(0))
	if err == nil || err == unix.EPERM {
		return true
	}
	return false
}

func checkProcessStatus(processHint string, pid int, stderrBuf *bytes.Buffer) error {
	var status syscall.WaitStatus
	pid, err := syscall.Wait4(pid, &status, syscall.WNOHANG, nil)
	if err != nil {
		return fmt.Errorf("failed to read qem%su process status: %w", processHint, err)
	}
	if pid > 0 {
		// child exited
		return fmt.Errorf("%s exited unexpectedly with exit code %d, stderr: %s", processHint, status.ExitStatus(), stderrBuf.String())
	}
	return nil
}

func pathsFromVolume(volume string) []string {
	return strings.SplitN(volume, ":", 3)
}

func extractTargetPath(paths []string) string {
	if len(paths) > 1 {
		return paths[1]
	}
	return paths[0]
}

func sigKill(pid int) error {
	return unix.Kill(pid, unix.SIGKILL)
}

func findProcess(pid int) (int, error) {
	if err := unix.Kill(pid, 0); err != nil {
		if err == unix.ESRCH {
			return -1, nil
		}
		return -1, fmt.Errorf("pinging QEMU process: %w", err)
	}
	return pid, nil
}
