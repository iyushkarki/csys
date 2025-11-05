package system

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type PortInfo struct {
	Port        int
	Protocol    string
	State       string
	ProcessName string
	PID         int32
	Memory      uint64
}

func GetListeningPorts() ([]PortInfo, error) {
	conns, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}

	var ports []PortInfo
	seen := make(map[int]bool)

	for _, conn := range conns {
		if conn.Status != "LISTEN" {
			continue
		}

		port := int(conn.Laddr.Port)
		if port == 0 || seen[port] {
			continue
		}
		seen[port] = true

		processName := "unknown"
		memory := uint64(0)

		if conn.Pid != 0 {
			p, err := process.NewProcess(conn.Pid)
			if err == nil {
				if name, err := p.Name(); err == nil {
					processName = name
				}
				if memInfo, err := p.MemoryInfo(); err == nil {
					memory = memInfo.RSS
				}
			}
		}

		protocol := "tcp"
		if conn.Type == 17 {
			protocol = "udp"
		}

		ports = append(ports, PortInfo{
			Port:        port,
			Protocol:    protocol,
			State:       conn.Status,
			ProcessName: processName,
			PID:         conn.Pid,
			Memory:      memory,
		})
	}

	return ports, nil
}

func GetProcessOnPort(port int) (*PortInfo, error) {
	ports, err := GetListeningPorts()
	if err != nil {
		return nil, err
	}

	for _, p := range ports {
		if p.Port == port {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("port %d not found", port)
}

func KillProcessOnPort(port int, force bool) error {
	portInfo, err := GetProcessOnPort(port)
	if err != nil {
		return err
	}

	if portInfo.PID == 0 {
		return fmt.Errorf("no process found for port %d", port)
	}

	pid := int(portInfo.PID)
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}

	if force {
		err = proc.Signal(syscall.SIGKILL)
	} else {
		err = proc.Signal(syscall.SIGTERM)
		if err == nil {
			time.Sleep(1 * time.Second)
			if proc2, err2 := os.FindProcess(pid); err2 == nil {
				if err3 := proc2.Signal(syscall.Signal(0)); err3 == nil {
					proc.Signal(syscall.SIGKILL)
				}
			}
		}
	}

	if err != nil {
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}

	return nil
}
