package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/iyushkarki/csys/internal/display"
	"github.com/iyushkarki/csys/internal/system"
	"github.com/spf13/cobra"
)

var Version = "dev"
var liveMode bool

var rootCmd = &cobra.Command{
	Use:               "csys",
	Short:             display.RootShort,
	Long:              display.RootLong,
	Version:           Version,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run: func(cmd *cobra.Command, args []string) {
		if liveMode {
			runLiveMode()
		} else {
			runSnapshot()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&liveMode, "live", "l", false, "Enable live monitoring mode (updates every 2 seconds)")
}

func runSnapshot() {
	runSnapshotWithTime(time.Time{})
}

func runSnapshotWithTime(timestamp time.Time) {
	diskInfo, err := system.GetDiskInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting disk info: %v\n", err)
		return
	}

	memInfo, err := system.GetMemoryInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting memory info: %v\n", err)
		return
	}

	cpuPercent, err := system.GetCPUUsage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting CPU info: %v\n", err)
		return
	}

	topProcs, err := system.GetTopProcessesByMemory(5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting process info: %v\n", err)
		return
	}

	if timestamp.IsZero() {
		output := display.FormatSystemOverview(diskInfo, memInfo, cpuPercent, topProcs)
		fmt.Println(output)
	} else {
		output := display.FormatSystemOverviewWithTime(diskInfo, memInfo, cpuPercent, topProcs, timestamp)
		fmt.Println(output)
	}
}

func runLiveMode() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		clearScreen()
		runSnapshotWithTime(time.Now())
		<-ticker.C
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
