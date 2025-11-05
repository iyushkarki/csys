package display

const (
	RootShort = "Beautiful system monitoring CLI"
	RootLong  = `csys - Monitor system metrics beautifully.

Quick Start:
  csys              System overview
  csys --live       Live monitoring
  csys ports        List listening ports
  csys ports kill   Kill process on port
  csys ports -h     Help for ports command`

	PortsShort = "Manage and monitor network ports"
	PortsLong  = `List listening ports or terminate processes.

EXAMPLES:
  csys ports              List all ports
  csys ports kill 3000              Kill single port
  csys ports kill 3000 8080         Kill multiple ports (space-separated)
  csys ports kill 3000 --force      Force kill without confirmation`

	ListShort = "List all listening ports with process info"
	ListLong  = `Display all listening ports with port number, protocol, process name, PID, memory.

EXAMPLES:
  csys ports list
  csys ports              (shorthand)`

	KillShort = "Kill process(es) running on specific port(s)"
	KillLong  = `Terminate process(es) on one or more ports.

By default: SIGTERM → wait 1s → SIGKILL if needed
Use --force (-f) to skip confirmation and force kill immediately.

EXAMPLES:
  csys ports kill 3000              Kill single port (with confirmation)
  csys ports kill 3000 8080         Kill multiple ports (space-separated)
  csys ports kill 3000 --force      Force kill without confirmation
  csys ports kill 3000 -f           Shorthand: -f for --force`
)
