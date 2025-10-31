package system

import (
	"sort"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	PID    int32
	Name   string
	Memory uint64
}

func GetTopProcessesByMemory(count int) ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var procInfos []ProcessInfo

	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			continue
		}

		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}

		procInfos = append(procInfos, ProcessInfo{
			PID:    p.Pid,
			Name:   name,
			Memory: memInfo.RSS,
		})
	}

	sort.Slice(procInfos, func(i, j int) bool {
		return procInfos[i].Memory > procInfos[j].Memory
	})

	if len(procInfos) > count {
		return procInfos[:count], nil
	}

	return procInfos, nil
}
