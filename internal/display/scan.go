package display

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/iyushkarki/csys/internal/system"
)

var (
	scanHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	sizeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	dirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6")).
			Bold(true)

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A3A3A3"))

	barFullChar  = "█"
	barEmptyChar = "░"
	barWidth     = 20
)

func RenderScanResult(result *system.ScanResult) string {
	var content string

	// Header
	content += scanHeaderStyle.Render("◈ DIRECTORY SCAN") + "\n"
	content += pathStyle.Render(result.RootPath) + "\n\n"

	// Summary
	content += fmt.Sprintf("Total Size: %s  •  Files: %d  •  Dirs: %d\n\n",
		sizeStyle.Render(humanize.IBytes(uint64(result.TotalSize))),
		result.FileCount,
		result.DirCount,
	)

	// Type Breakdown
	if len(result.TypeBreakdown) > 0 {
		content += scanHeaderStyle.Render("◈ TYPE BREAKDOWN") + "\n"
		var breakdown []string
		for i, tb := range result.TypeBreakdown {
			if i >= 5 {
				break
			}
			ext := tb.Extension
			if ext == "no-ext" {
				ext = "misc"
			}
			breakdown = append(breakdown, fmt.Sprintf("%s: %s", ext, humanize.IBytes(uint64(tb.Size))))
		}
		content += strings.Join(breakdown, "  |  ") + "\n\n"
	}

	content += scanHeaderStyle.Render("◈ TOP SPACE CONSUMERS") + "\n"

	maxSize := int64(0)
	if len(result.Items) > 0 {
		maxSize = result.Items[0].Size
	}

	for i, item := range result.Items {
		if i >= 10 {
			break
		}

		percent := float64(item.Size) / float64(maxSize)
		barLen := int(percent * float64(barWidth))
		bar := strings.Repeat(barFullChar, barLen) + strings.Repeat(barEmptyChar, barWidth-barLen)

		barStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#5A5A5A")) // Default gray
		if item.IsDir {
			barStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6")) // Blue for dirs
		} else if percent > 0.5 {
			barStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EAB308")) // Yellow for large files
		}

		name := item.Name
		if item.IsDir {
			name = dirStyle.Render(name + "/")
		} else {
			name = fileStyle.Render(name)
		}

		size := humanize.IBytes(uint64(item.Size))

		content += fmt.Sprintf("%2d. %-30s  %10s  %s\n",
			i+1,
			truncate(name, 30),
			size,
			barStyle.Render(bar),
		)
	}

	return borderStyle.Render(content)
}

func RenderDiskUsage(info *system.DiskInfo) string {
	var content string

	var primaryDisks []system.DiskPartition
	var systemDisks []system.DiskPartition

	for _, disk := range info.Partitions {
		if disk.Category == "primary" {
			primaryDisks = append(primaryDisks, disk)
		} else {
			systemDisks = append(systemDisks, disk)
		}
	}

	if len(primaryDisks) > 0 {
		content += scanHeaderStyle.Render("◈ PRIMARY STORAGE") + "\n\n"
		for _, disk := range primaryDisks {
			content += fmt.Sprintf("%s  %s\n",
				dirStyle.Render(disk.Label),
				pathStyle.Render(disk.Device),
			)

			percent := disk.Percent
			bar := createProgressBar(percent, 30)
			percStr := getColoredPercent(percent)

			used := humanize.IBytes(disk.Used)
			total := humanize.IBytes(disk.Total)
			free := humanize.IBytes(disk.Free)

			content += fmt.Sprintf("%s %s\n", bar, percStr)
			content += fmt.Sprintf("%s used  •  %s free  •  %s total\n\n",
				sizeStyle.Render(used),
				fileStyle.Render(free),
				fileStyle.Render(total),
			)
		}
	}

	if len(systemDisks) > 0 {
		content += scanHeaderStyle.Render("◈ SYSTEM VOLUMES") + "\n"
		for _, disk := range systemDisks {
			// Compact view: Name (Usage)
			name := truncate(disk.Mountpoint, 30)
			if disk.Label != disk.Mountpoint {
				name = disk.Label
			}

			used := humanize.IBytes(disk.Used)
			percent := getColoredPercent(disk.Percent)

			content += fmt.Sprintf("  • %-30s  %s used (%s)\n",
				fileStyle.Render(name),
				sizeStyle.Render(used),
				percent,
			)
		}
	}

	return borderStyle.Render(content)
}
