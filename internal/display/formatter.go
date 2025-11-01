package display

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/iyushkarki/csys/internal/system"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#626262")).
			Padding(1, 2)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500"))

	criticalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	barFilled = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	barWarning = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500"))

	barCritical = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	barEmpty = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C3C3C"))

	processStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
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

	header := "SYSTEM OVERVIEW"
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
	var lines string

	// Disk
	if len(diskInfo.Partitions) > 0 {
		disk := diskInfo.Partitions[0]
		diskPercent := disk.Percent
		diskUsed := humanize.IBytes(disk.Used)
		diskTotal := humanize.IBytes(disk.Total)
		diskBar := createProgressBar(diskPercent, 20)
		diskPercStr := getColoredPercent(diskPercent)

		lines += fmt.Sprintf("◉ Disk    %s %s  %s / %s\n",
			diskBar,
			diskPercStr,
			diskUsed,
			diskTotal,
		)
	}

	// Memory
	memPercent := memInfo.UsedPercent
	memUsed := humanize.IBytes(memInfo.Used)
	memTotal := humanize.IBytes(memInfo.Total)
	memBar := createProgressBar(memPercent, 20)
	memPercStr := getColoredPercent(memPercent)

	lines += fmt.Sprintf("▣ Memory  %s %s  %s / %s\n",
		memBar,
		memPercStr,
		memUsed,
		memTotal,
	)

	// CPU
	cpuBar := createProgressBar(cpuPercent, 20)
	cpuPercStr := getColoredPercent(cpuPercent)

	lines += fmt.Sprintf("△ CPU     %s %s",
		cpuBar,
		cpuPercStr,
	)

	return lines
}

func formatProcessSection(procs []system.ProcessInfo) string {
	header := titleStyle.Render("▲ TOP MEMORY PROCESSES") + "\n"

	if len(procs) == 0 {
		return header + "  No processes found"
	}

	var processes string
	for i, proc := range procs {
		if i >= 5 {
			break
		}
		memSize := humanize.IBytes(proc.Memory)
		processes += fmt.Sprintf("  %d  %s  %s\n",
			i+1,
			processStyle.Render(truncate(proc.Name, 35)),
			normalStyle.Render(memSize),
		)
	}

	return header + processes
}

func createProgressBar(percent float64, width int) string {
	if percent > 100 {
		percent = 100
	}
	if percent < 0 {
		percent = 0
	}

	filled := int(float64(width) * percent / 100)
	empty := width - filled

	var barStyle lipgloss.Style
	if percent >= 90 {
		barStyle = barCritical
	} else if percent >= 70 {
		barStyle = barWarning
	} else {
		barStyle = barFilled
	}

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	filledPart := barStyle.Render(bar)

	emptyBar := ""
	for i := 0; i < empty; i++ {
		emptyBar += "░"
	}
	emptyPart := barEmpty.Render(emptyBar)

	return filledPart + emptyPart
}

func getColoredPercent(percent float64) string {
	color := getColorForPercent(percent)
	return color.Render(fmt.Sprintf("%.0f%%", percent))
}

func getColorForPercent(percent float64) lipgloss.Style {
	if percent >= 90 {
		return criticalStyle
	} else if percent >= 70 {
		return warningStyle
	}
	return normalStyle
}

func truncate(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
