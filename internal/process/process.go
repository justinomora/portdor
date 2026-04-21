package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func Stop(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find process %d: %w", pid, err)
	}
	return proc.Signal(syscall.SIGTERM)
}

func Kill(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find process %d: %w", pid, err)
	}
	return proc.Signal(syscall.SIGKILL)
}

func IsAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	// Attempt a non-blocking wait to reap any zombie so signal(0) reflects
	// true liveness. If wait4 returns this pid, the process has exited.
	var ws syscall.WaitStatus
	wpid, _ := syscall.Wait4(pid, &ws, syscall.WNOHANG, nil)
	if wpid == pid {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

func Restart(pid int, command, cwd string) (int, error) {
	if pid > 0 && IsAlive(pid) {
		Stop(pid)
		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)
			if !IsAlive(pid) {
				break
			}
		}
		if IsAlive(pid) {
			Kill(pid)
		}
	}

	return Spawn(command, cwd)
}

func Spawn(command, cwd string) (int, error) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = cwd
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("spawn %q: %w", command, err)
	}
	go cmd.Wait()

	return cmd.Process.Pid, nil
}
