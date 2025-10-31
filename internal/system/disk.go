package system

import (
	"github.com/shirou/gopsutil/v3/disk"
)

type DiskInfo struct {
	Partitions []DiskPartition
}

type DiskPartition struct {
	Mountpoint string
	Device     string
	Total      uint64
	Used       uint64
	Free       uint64
	Percent    float64
}

func GetDiskInfo() (*DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var diskPartitions []DiskPartition

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		diskPartition := DiskPartition{
			Mountpoint: partition.Mountpoint,
			Device:     partition.Device,
			Total:      usage.Total,
			Used:       usage.Used,
			Free:       usage.Free,
			Percent:    usage.UsedPercent,
		}

		if partition.Mountpoint == "/" {
			diskPartitions = append([]DiskPartition{diskPartition}, diskPartitions...)
		} else {
			diskPartitions = append(diskPartitions, diskPartition)
		}
	}

	if len(diskPartitions) > 0 {
		return &DiskInfo{Partitions: diskPartitions[:1]}, nil
	}

	return &DiskInfo{Partitions: diskPartitions}, nil
}
