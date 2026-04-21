package registry

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type HealthInfo struct {
	Alive  bool
	MemPct float64
	CPUPct float64
}

func (r *Registry) Check(name string) (HealthInfo, error) {
	svc, err := r.Get(name)
	if err != nil {
		return HealthInfo{}, err
	}
	if svc.PID == 0 {
		return HealthInfo{Alive: false}, nil
	}

	info := checkPID(svc.PID)

	status := StatusStopped
	if info.Alive {
		status = StatusRunning
	} else if svc.Status == StatusRunning {
		status = StatusCrashed
	}
	r.SetPID(name, svc.PID, status)

	return info, nil
}

func (r *Registry) CheckAll() {
	for _, svc := range r.List() {
		r.Check(svc.Name)
	}
}

func checkPID(pid int) HealthInfo {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return HealthInfo{Alive: false}
	}
	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return HealthInfo{Alive: false}
	}

	info := HealthInfo{Alive: true}
	info.MemPct, info.CPUPct = getProcStats(pid)
	return info
}

func getProcStats(pid int) (memPct, cpuPct float64) {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "%mem=,%cpu=")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) >= 2 {
		fmt.Sscanf(parts[0], "%f", &memPct)
		fmt.Sscanf(parts[1], "%f", &cpuPct)
	}
	return
}

func PIDForPort(port int) int {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("lsof", "-iTCP:"+strconv.Itoa(port), "-sTCP:LISTEN", "-n", "-P", "-t")
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("ss -tlnp | grep ':%d ' | grep -oP 'pid=\\K[0-9]+'", port))
	}
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return pid
}
