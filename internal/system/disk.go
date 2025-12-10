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
	Label      string
	Category   string // "primary" or "system"
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

func GetFullDiskInfo() (*DiskInfo, error) {
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

		label, category := getDiskLabelAndCategory(partition.Mountpoint)

		diskPartition := DiskPartition{
			Mountpoint: partition.Mountpoint,
			Device:     partition.Device,
			Total:      usage.Total,
			Used:       usage.Used,
			Free:       usage.Free,
			Percent:    usage.UsedPercent,
			Label:      label,
			Category:   category,
		}

		if category == "primary" {
			// Prepend primary disks
			diskPartitions = append([]DiskPartition{diskPartition}, diskPartitions...)
		} else {
			diskPartitions = append(diskPartitions, diskPartition)
		}
	}

	return &DiskInfo{Partitions: diskPartitions}, nil
}

func getDiskLabelAndCategory(mountpoint string) (string, string) {
	switch mountpoint {
	case "/":
		return "System Root", "primary"
	case "/System/Volumes/Data", "/home":
		return "User Data", "primary"
	}

	// External drives (Mac: /Volumes/..., Linux: /media/...)
	if len(mountpoint) > 9 && mountpoint[:9] == "/Volumes/" {
		return "External: " + mountpoint[9:], "primary"
	}
	if len(mountpoint) > 7 && mountpoint[:7] == "/media/" {
		return "External: " + mountpoint[7:], "primary"
	}

	// System volumes
	return mountpoint, "system"
}
