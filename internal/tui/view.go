package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/iyushkarki/csys/internal/cleaners"
	"github.com/iyushkarki/csys/internal/display"
)

var (
	accent = lipgloss.Color("#7D56F4")
	green  = lipgloss.Color("#04B575")
	amber  = lipgloss.Color("#FFA500")
	red    = lipgloss.Color("#FF0000")
	gray   = lipgloss.Color("#626262")
	dim    = lipgloss.Color("#3C3C3C")

	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(accent)
	sizeStyle   = lipgloss.NewStyle().Bold(true).Foreground(green)
	grayStyle   = lipgloss.NewStyle().Foreground(gray)
	dimStyle    = lipgloss.NewStyle().Foreground(dim)
	amberStyle  = lipgloss.NewStyle().Foreground(amber)
	redStyle    = lipgloss.NewStyle().Bold(true).Foreground(red)
	greenStyle  = lipgloss.NewStyle().Foreground(green)
	cursorStyle = lipgloss.NewStyle().Bold(true).Foreground(accent)

	panelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(gray).
			Padding(0, 1)
)

const nameWidth = 36

func (m *Model) View() string {
	if m.phase == done {
		return m.doneView()
	}

	var b strings.Builder
	b.WriteString(m.headerView() + "\n\n")
	b.WriteString(m.listView())
	b.WriteString("\n" + m.detailView())
	b.WriteString("\n" + m.footerView() + "\n")
	return b.String()
}

func (m *Model) headerView() string {
	header := titleStyle.Render("✦ csys clean")
	if m.scanning {
		header += "  " + m.spin.View() + grayStyle.Render(" scanning…")
	}
	stats := fmt.Sprintf("found %s · selected %s",
		humanize.IBytes(m.foundSize()),
		sizeStyle.Render(humanize.IBytes(m.selectedSize())))
	if m.phase == cleaning {
		stats = m.spin.View() + grayStyle.Render(" cleaning…")
	}
	gap := m.width - lipgloss.Width(header) - lipgloss.Width(stats) - 2
	if gap < 2 {
		gap = 2
	}
	return " " + header + strings.Repeat(" ", gap) + stats
}

func (m *Model) listView() string {
	var lines []string
	cursorLine := 0
	lastTier := cleaners.Tier(-1)

	for i, it := range m.items {
		if it.target.Tier != lastTier {
			if it.target.Tier == cleaners.Safe {
				lines = append(lines, " "+greenStyle.Render("● SAFE")+grayStyle.Render(" — regenerates automatically"))
			} else {
				if len(lines) > 0 {
					lines = append(lines, "")
				}
				lines = append(lines, " "+amberStyle.Render("● CAREFUL")+grayStyle.Render(" — read the note first"))
			}
			lastTier = it.target.Tier
		}
		if i == m.cursor {
			cursorLine = len(lines)
		}
		lines = append(lines, m.rowView(i, it))
	}

	maxRows := m.height - 14
	if maxRows < 5 {
		maxRows = 5
	}
	if len(lines) > maxRows {
		start := cursorLine - maxRows/2
		if start < 0 {
			start = 0
		}
		if start+maxRows > len(lines) {
			start = len(lines) - maxRows
		}
		lines = lines[start : start+maxRows]
	}
	return strings.Join(lines, "\n") + "\n"
}

func (m *Model) rowView(i int, it *item) string {
	cursor := "  "
	if i == m.cursor && m.phase == selecting {
		cursor = cursorStyle.Render("▸ ")
	}

	marker := "[ ]"
	switch {
	case it.status == statusWorking:
		marker = " " + m.spin.View() + " "
	case it.status == statusDone:
		marker = greenStyle.Render(" ✓ ")
	case it.status == statusFailed:
		marker = redStyle.Render(" ✗ ")
	case it.selected:
		marker = cursorStyle.Render("[x]")
	}

	name := display.Truncate(it.target.Name, nameWidth)
	nameStyled := lipgloss.NewStyle().Width(nameWidth).Render(name)
	if !it.selected && it.status == statusIdle {
		nameStyled = grayStyle.Width(nameWidth).Render(name)
	}

	row := fmt.Sprintf(" %s%s %s %9s  %s",
		cursor, marker, nameStyled,
		humanize.IBytes(it.target.Size),
		m.badgeView(it))
	return row
}

func (m *Model) badgeView(it *item) string {
	if it.status == statusFailed {
		return redStyle.Render(display.Truncate(it.err.Error(), 40))
	}
	used := it.target.LastUsed
	label := "used " + cleaners.Ago(used)
	switch {
	case used.IsZero():
		return dimStyle.Render("—")
	case time.Since(used) < 7*24*time.Hour:
		return amberStyle.Render(label)
	default:
		return grayStyle.Render(label)
	}
}

func (m *Model) detailView() string {
	it := m.currentItem()
	if it == nil {
		return ""
	}
	t := it.target

	header := lipgloss.NewStyle().Bold(true).Render(t.Name)
	if !t.LastUsed.IsZero() {
		header += grayStyle.Render(" · used " + cleaners.Ago(t.LastUsed))
	}
	if t.Note != "" {
		header += "  " + amberStyle.Render(t.Note)
	}

	label := func(s string) string { return grayStyle.Render(fmt.Sprintf("%-6s", s)) }
	body := []string{header}
	if t.Explain.What != "" {
		body = append(body, label("What")+t.Explain.What)
	}
	if t.Explain.After != "" {
		body = append(body, label("After")+t.Explain.After)
	}
	body = append(body, label("Will")+t.Action(m.opts))

	width := m.width - 4
	if width < 40 {
		width = 40
	}
	return panelStyle.Width(width).Render(strings.Join(body, "\n"))
}

func (m *Model) footerView() string {
	if m.phase == cleaning {
		return grayStyle.Render("  cleaning — hang tight…")
	}
	keys := []string{"↑/↓ move", "space toggle", "a all safe", "u none", "enter clean", "q quit"}
	return grayStyle.Render("  " + strings.Join(keys, " · "))
}

func (m *Model) doneView() string {
	if len(m.items) == 0 {
		return panelStyle.Render(titleStyle.Render("✦ ALL CLEAN")+"\n\n"+
			greenStyle.Render("Nothing to reclaim — your system is already tidy.")) + "\n"
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("✦ CLEANED UP") + "\n\n")

	if m.freed > 0 {
		b.WriteString(fmt.Sprintf("Freed %s\n", sizeStyle.Render(humanize.IBytes(m.freed))))
	}
	if m.trashed > 0 {
		b.WriteString(fmt.Sprintf("Moved %s to Trash %s\n",
			sizeStyle.Render(humanize.IBytes(m.trashed)),
			grayStyle.Render("(recoverable — empty Trash to reclaim)")))
	}
	for _, it := range m.items {
		if it.status == statusFailed {
			b.WriteString(redStyle.Render("✗ ") + it.target.Name + grayStyle.Render(" — "+it.err.Error()) + "\n")
		}
	}

	if m.diskBefore != nil && m.diskAfter != nil &&
		len(m.diskBefore.Partitions) > 0 && len(m.diskAfter.Partitions) > 0 {
		before, after := m.diskBefore.Partitions[0], m.diskAfter.Partitions[0]
		b.WriteString(grayStyle.Render(fmt.Sprintf("\nDisk: %s used (%.0f%%) → %s used (%.0f%%)",
			humanize.IBytes(before.Used), before.Percent,
			humanize.IBytes(after.Used), after.Percent)) + "\n")
	}

	if lifetime := cleaners.LifetimeFreed(); lifetime > 0 {
		b.WriteString(grayStyle.Render("Lifetime cleaned with csys: ") +
			sizeStyle.Render(humanize.IBytes(lifetime)) + "\n")
	}

	return panelStyle.Render(strings.TrimRight(b.String(), "\n")) + "\n" +
		grayStyle.Render("  q to quit") + "\n"
}
