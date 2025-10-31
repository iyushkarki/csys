package system

import (
	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryInfo struct {
	Total       uint64
	Available   uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func GetMemoryInfo() (*MemoryInfo, error) {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryInfo{
		Total:       memStats.Total,
		Available:   memStats.Available,
		Used:        memStats.Used,
		Free:        memStats.Free,
		UsedPercent: memStats.UsedPercent,
	}, nil
}
