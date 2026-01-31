package display

const (
	RootShort = "Beautiful system monitoring CLI"
	RootLong  = `csys - Beautiful system monitoring CLI

SYSTEM
  csys              System overview
  csys --live       Live monitoring (2s refresh)

STORAGE
  csys scan         Scan current directory
  csys scan disk    Show all disk partitions

NETWORK
  csys ports        List listening ports
  csys ports kill   Kill process on port

GIT
  csys g st         Quick git status
  csys g log [n]    Show commit log (default: 10)
  csys g sync       Reset to origin/main
  csys g clean      Delete all branches except current
  csys g soft [n]   Soft reset HEAD~n
  csys g ac "msg"   Add + commit
  csys g acp "msg"  Add + commit + push
  csys g undo       Undo last commit
  csys g wip        Quick WIP commit
  csys g amend      Amend last commit
  csys g rb         Rebase current branch onto main`

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

	ScanShort = "Analyze directory storage usage"
	ScanLong  = `Scan a directory to see a breakdown of file types and top space consumers.

EXAMPLES:
  csys scan                 Scan current directory
  csys scan --path ~/Downloads   Scan specific directory`

	ScanDiskShort = "Show usage of all disk partitions"
	ScanDiskLong  = `Scan and display storage usage for all mounted disk partitions.`
)
