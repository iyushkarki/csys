package cmd

import (
	"fmt"
	"os"

	"github.com/iyushkarki/csys/internal/display"
	"github.com/iyushkarki/csys/internal/system"
	"github.com/spf13/cobra"
)

var (
	scanPath  string
	scanLimit int
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Analyze directory storage usage",
	Long:  `Scan a directory to see a breakdown of file types and top space consumers.`,
	Run: func(cmd *cobra.Command, args []string) {
		if scanPath == "" {
			var err error
			scanPath, err = os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				return
			}
		}

		result, err := system.ScanDirectory(scanPath)
		if err != nil {
			fmt.Printf("Error scanning directory: %v\n", err)
			return
		}

		fmt.Println(display.RenderScanResult(result))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&scanPath, "path", "p", "", "Directory to scan (default: current)")
	scanCmd.Flags().IntVarP(&scanLimit, "limit", "l", 10, "Number of top items to show")

	scanCmd.AddCommand(scanDiskCmd)
}

var scanDiskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Show usage of all disk partitions",
	Long:  `Scan and display storage usage for all mounted disk partitions.`,
	Run: func(cmd *cobra.Command, args []string) {

		info, err := system.GetFullDiskInfo()
		if err != nil {
			fmt.Printf("Error getting disk info: %v\n", err)
			return
		}

		fmt.Println(display.RenderDiskUsage(info))
	},
}
