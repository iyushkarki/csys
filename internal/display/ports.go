package display

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/iyushkarki/csys/internal/system"
)

var (
	portHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	portLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	portNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	portProcessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
)

func FormatPortsList(ports []system.PortInfo) string {
	if len(ports) == 0 {
		return borderStyle.Render(
			portHeaderStyle.Render("◈ LISTENING PORTS") + "\n" +
				"  No ports currently listening",
		)
	}

	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Port < ports[j].Port
	})

	var content string
	header := portHeaderStyle.Render("◈ LISTENING PORTS")
	content += header + "\n\n"

	for i, port := range ports {
		portColor := getPortTypeColor(port.Port)
		portNum := portColor.Render(fmt.Sprintf("%5d", port.Port))
		protocol := formatProtocol(port.Protocol)
		processName := portProcessStyle.Render(truncate(port.ProcessName, 25))
		pid := portLabelStyle.Render(fmt.Sprintf("[PID: %d]", port.PID))
		memory := normalStyle.Render(humanize.IBytes(port.Memory))

		content += fmt.Sprintf("  %d  ⟳ %s  %s  %s  %s  %s\n",
			i+1,
			portNum,
			protocol,
			processName,
			pid,
			memory,
		)
	}

	return borderStyle.Render(content)
}

func FormatKillConfirmation(portInfo *system.PortInfo) string {
	portNum := portNumberStyle.Render(fmt.Sprintf("%d", portInfo.Port))
	protocol := formatProtocol(portInfo.Protocol)
	processName := portProcessStyle.Render(truncate(portInfo.ProcessName, 30))
	pid := portLabelStyle.Render(fmt.Sprintf("[PID: %d]", portInfo.PID))
	memory := normalStyle.Render(humanize.IBytes(portInfo.Memory))

	var content string
	content += portHeaderStyle.Render("⚠ KILL CONFIRMATION") + "\n\n"
	content += fmt.Sprintf("  ⟳ %s  %s  %s  %s  %s\n\n",
		portNum,
		protocol,
		processName,
		pid,
		memory,
	)
	content += labelStyle.Render("  Confirm termination? [y/N]: ")

	return borderStyle.Render(content)
}

func FormatKillSuccess(port int, processName string) string {
	var content string
	content += successStyle.Render(fmt.Sprintf("✓ Port %d (%s) killed successfully", port, processName))
	return content
}

func FormatKillError(port int, err error) string {
	var content string
	content += errorStyle.Render(fmt.Sprintf("✗ Failed to kill port %d: %v", port, err))
	return content
}

func FormatPortNotFound(port int) string {
	var content string
	content += errorStyle.Render(fmt.Sprintf("✗ Port %d is not listening", port))
	return content
}

func formatProtocol(proto string) string {
	if proto == "tcp" {
		return normalStyle.Render("→ TCP")
	}
	return warningStyle.Render("⚡ UDP")
}

func getPortTypeColor(port int) lipgloss.Style {
	if port < 1024 {
		return criticalStyle
	} else if port <= 10000 {
		return warningStyle
	}
	return normalStyle
}
