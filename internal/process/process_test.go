package process_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/jmora/portdor/internal/process"
)

func startSleepProcess(t *testing.T) *os.Process {
	t.Helper()
	cmd := exec.Command("sleep", "60")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start sleep process: %v", err)
	}
	t.Cleanup(func() { cmd.Process.Kill() })
	return cmd.Process
}

func TestStop(t *testing.T) {
	proc := startSleepProcess(t)

	if err := process.Stop(proc.Pid); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if process.IsAlive(proc.Pid) {
		t.Error("expected process to be stopped after SIGTERM")
	}
}

func TestKill(t *testing.T) {
	proc := startSleepProcess(t)

	if err := process.Kill(proc.Pid); err != nil {
		t.Fatalf("Kill: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if process.IsAlive(proc.Pid) {
		t.Error("expected process to be dead after SIGKILL")
	}
}

func TestIsAlive(t *testing.T) {
	proc := startSleepProcess(t)

	if !process.IsAlive(proc.Pid) {
		t.Error("expected process to be alive")
	}

	proc.Kill()
	time.Sleep(50 * time.Millisecond)

	if process.IsAlive(proc.Pid) {
		t.Error("expected process to be dead after kill")
	}
}

func TestStopInvalidPID(t *testing.T) {
	err := process.Stop(-1)
	if err == nil {
		t.Fatal("expected error for invalid PID, got nil")
	}
}
