package display

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/iyushkarki/csys/internal/cleaners"
	"github.com/iyushkarki/csys/internal/system"
)

func RenderCleanList(targets []*cleaners.Target) string {
	var content string

	content += scanHeaderStyle.Render("◈ RECLAIMABLE SPACE") + "\n\n"

	maxSize := uint64(0)
	var total uint64
	for _, t := range targets {
		if t.Size > maxSize {
			maxSize = t.Size
		}
		total += t.Size
	}

	lastTier := cleaners.Tier(-1)
	for _, t := range targets {
		if t.Tier != lastTier {
			if t.Tier == cleaners.Safe {
				content += normalStyle.Render("● SAFE") + labelStyle.Render("  regenerates automatically") + "\n"
			} else {
				content += "\n" + warningStyle.Render("● CAREFUL") + labelStyle.Render("  read the note first") + "\n"
			}
			lastTier = t.Tier
		}

		percent := float64(t.Size) / float64(maxSize)
		barLen := int(percent * float64(barWidth))
		bar := strings.Repeat(barFullChar, barLen) + strings.Repeat(barEmptyChar, barWidth-barLen)

		barStyled := barFilled.Render(bar)
		if t.Tier == cleaners.Careful {
			barStyled = barWarning.Render(bar)
		}

		line := fmt.Sprintf("  %-36s %10s  %s  %s",
			Truncate(t.Name, 36),
			sizeStyle.Render(humanize.IBytes(t.Size)),
			barStyled,
			labelStyle.Render("used "+cleaners.Ago(t.LastUsed)),
		)
		if t.Note != "" {
			line += "  " + labelStyle.Render(t.Note)
		}
		content += line + "\n"
	}

	content += "\n" + fmt.Sprintf("Total reclaimable: %s", sizeStyle.Render(humanize.IBytes(total)))

	return borderStyle.Render(content)
}

func FormatCleanProgress(name string, freed uint64, err error) string {
	if err != nil {
		return fmt.Sprintf("  %s %s — %v", criticalStyle.Render("✗"), name, err)
	}
	return fmt.Sprintf("  %s %s — freed %s",
		normalStyle.Render("✓"), name, sizeStyle.Render(humanize.IBytes(freed)))
}

func RenderCleanSummary(freed uint64, before, after *system.DiskInfo) string {
	var content string
	content += titleStyle.Render("✦ CLEANED UP") + "\n\n"
	content += fmt.Sprintf("Freed %s\n", sizeStyle.Render(humanize.IBytes(freed)))

	if before != nil && after != nil &&
		len(before.Partitions) > 0 && len(after.Partitions) > 0 {
		b, a := before.Partitions[0], after.Partitions[0]
		content += "\n"
		content += fmt.Sprintf("◉ Disk    %s %s  %s / %s\n",
			createProgressBar(a.Percent, 20),
			getColoredPercent(a.Percent),
			humanize.IBytes(a.Used),
			humanize.IBytes(a.Total),
		)
		content += labelStyle.Render(fmt.Sprintf("          was %s used (%.0f%%) → now %s used (%.0f%%)",
			humanize.IBytes(b.Used), b.Percent, humanize.IBytes(a.Used), a.Percent))
	}

	return borderStyle.Render(content)
}

func RenderNothingToClean() string {
	return borderStyle.Render(
		titleStyle.Render("✦ ALL CLEAN") + "\n\n" +
			normalStyle.Render("Nothing to reclaim — your caches are already tidy."))
}
