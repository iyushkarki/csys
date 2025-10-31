package system

import (
	"github.com/shirou/gopsutil/v3/cpu"
)

func GetCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}

	if len(percentages) > 0 {
		return percentages[0], nil
	}

	return 0, nil
}

func GetCPUCount() (int, error) {
	count, err := cpu.Counts(false)
	if err != nil {
		return 0, err
	}
	return count, nil
}
