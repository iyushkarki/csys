package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iyushkarki/csys/internal/display"
	"github.com/iyushkarki/csys/internal/system"
	"github.com/spf13/cobra"
)

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: display.PortsShort,
	Long:  display.PortsLong,
	Run: func(cmd *cobra.Command, args []string) {
		runPortsList()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: display.ListShort,
	Long:  display.ListLong,
	Run: func(cmd *cobra.Command, args []string) {
		runPortsList()
	},
}

var killCmd = &cobra.Command{
	Use:   "kill <port> [port2] [port3]...",
	Short: display.KillShort,
	Long:  display.KillLong,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		runPortsKill(args, force)
	},
}

func init() {
	rootCmd.AddCommand(portsCmd)
	portsCmd.AddCommand(listCmd)
	portsCmd.AddCommand(killCmd)
	killCmd.Flags().BoolP("force", "f", false, "Force kill with SIGKILL")
}

func runPortsList() {
	ports, err := system.GetListeningPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting listening ports: %v\n", err)
		return
	}

	output := display.FormatPortsList(ports)
	fmt.Println(output)
}

func runPortsKill(portStrings []string, force bool) {
	var ports []int
	for _, portStr := range portStrings {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid port '%s' (must be a number)\n", portStr)
			fmt.Fprintf(os.Stderr, "Use 'csys ports kill --help' for usage examples\n")
			return
		}
		if port <= 0 || port > 65535 {
			fmt.Fprintf(os.Stderr, "Error: Port %d out of range (must be 1-65535)\n", port)
			fmt.Fprintf(os.Stderr, "Use 'csys ports kill --help' for usage examples\n")
			return
		}
		ports = append(ports, port)
	}

	for _, port := range ports {
		portInfo, err := system.GetProcessOnPort(port)
		if err != nil {
			fmt.Println(display.FormatPortNotFound(port))
			continue
		}

		if !force {
			fmt.Println(display.FormatKillConfirmation(portInfo))
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(text)) != "y" {
				fmt.Println("Kill cancelled")
				continue
			}
		}

		err = system.KillProcessOnPort(port, force)
		if err != nil {
			fmt.Println(display.FormatKillError(port, err))
		} else {
			fmt.Println(display.FormatKillSuccess(port, portInfo.ProcessName))
		}
	}
}
