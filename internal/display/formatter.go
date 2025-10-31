package display

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/iyushkarki/csys/internal/system"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	borderStyle  = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(1, 2)
	metricStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	processStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	barFgStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
)

func FormatSystemOverview(
	diskInfo *system.DiskInfo,
	memInfo *system.MemoryInfo,
	cpuPercent float64,
	topProcs []system.ProcessInfo,
) string {
	return FormatSystemOverviewWithTime(diskInfo, memInfo, cpuPercent, topProcs, time.Time{})
}

func FormatSystemOverviewWithTime(
	diskInfo *system.DiskInfo,
	memInfo *system.MemoryInfo,
	cpuPercent float64,
	topProcs []system.ProcessInfo,
	timestamp time.Time,
) string {
	var content string

	header := "📊  SYSTEM OVERVIEW"
	if !timestamp.IsZero() {
		header += fmt.Sprintf("  (Last updated: %s)", timestamp.Format("15:04:05"))
	}
	content += titleStyle.Render(header) + "\n\n"

	content += formatMetricsSection(diskInfo, memInfo, cpuPercent)
	content += "\n\n"
	content += formatProcessSection(topProcs)

	return borderStyle.Render(content)
}

func formatMetricsSection(diskInfo *system.DiskInfo, memInfo *system.MemoryInfo, cpuPercent float64) string {
	diskPercent := 0
	diskUsed := "0"
	diskTotal := "0"
	if len(diskInfo.Partitions) > 0 {
		disk := diskInfo.Partitions[0]
		diskPercent = int(disk.Percent)
		diskUsed = humanize.Bytes(disk.Used)
		diskTotal = humanize.Bytes(disk.Total)
	}

	memPercent := int(memInfo.UsedPercent)
	memUsed := humanize.Bytes(memInfo.Used)
	memTotal := humanize.Bytes(memInfo.Total)
	cpuPercentInt := int(cpuPercent)

	diskBar := createProgressBar(diskPercent, 10)
	memBar := createProgressBar(memPercent, 10)
	cpuBar := createProgressBar(cpuPercentInt, 10)

	diskPercStr := getColoredPercent(diskPercent)
	memPercStr := getColoredPercent(memPercent)
	cpuPercStr := getColoredPercent(cpuPercentInt)

	lines := fmt.Sprintf("  💾  Disk      %6s / %6s   %s  %s\n", diskUsed, diskTotal, diskBar, diskPercStr)
	lines += fmt.Sprintf("  🧠  Memory    %6s / %6s   %s  %s\n", memUsed, memTotal, memBar, memPercStr)
	lines += fmt.Sprintf("  ⚡  CPU        %6.1f%% / 100%%   %s  %s", cpuPercent, cpuBar, cpuPercStr)

	return metricStyle.Render(lines)
}

func getColoredPercent(percent int) string {
	color := getColorForPercent(percent)
	return color.Render(fmt.Sprintf("%3d%%", percent))
}

func formatProcessSection(procs []system.ProcessInfo) string {
	header := labelStyle.Render("📈  TOP MEMORY PROCESSES:")
	var processes string

	if len(procs) == 0 {
		processes = "\n  No processes found"
	} else {
		processes = "\n"
		for i, proc := range procs {
			if i >= 5 {
				break
			}
			memMB := float64(proc.Memory) / 1024 / 1024
			line := fmt.Sprintf("  %d  • %-27s  %8.1f MB\n",
				i+1,
				truncate(proc.Name, 27),
				memMB,
			)
			processes += line
		}
	}

	return header + processStyle.Render(processes)
}

func createProgressBar(percent int, width int) string {
	if percent > 100 {
		percent = 100
	}
	if percent < 0 {
		percent = 0
	}

	filled := (percent*width + 50) / 100
	bar := ""

	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return barFgStyle.Render(bar)
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

func getColorForPercent(percent int) lipgloss.Style {
	if percent >= 80 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	} else if percent >= 60 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
}
