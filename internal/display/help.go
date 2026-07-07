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
  csys clean          Interactive cache cleanup (TUI)
  csys clean modules  Find stale project artifacts
  csys clean -s       Clean all safe caches instantly
  csys clean -n       Dry run — show reclaimable space

NETWORK
  csys ports        List listening ports
  csys ports kill   Kill process on port

GIT
  cs gst            Quick git status
  cs glog [n]       Show commit log (default: 10)
  cs gsync          Reset to origin/main
  cs gclean         Delete all branches except current
  cs gsoft [n]      Soft reset HEAD~n
  cs gac "msg"      Add + commit
  cs gco <branch>   Switch to branch
  cs gcb <branch>   Create + switch to branch
  cs gbrn <name>    Rename current branch
  cs gpush          Push to remote
  cs gpull          Pull latest changes
  cs gfp            Force push (with lease)
  cs gundo          Undo last commit
  cs gwip           Quick WIP commit
  cs gamend         Amend last commit
  cs gamend "msg"   Amend with new message
  cs grb            Rebase onto main`

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

	CleanShort = "Find and clean developer caches to free disk space"
	CleanLong  = `Scan ~30 known locations (package manager caches, Xcode, simulators,
app leftovers, updater junk) and clean them interactively or automatically.

Every target shows what it is, what happens after cleaning, and when it was
last used. SAFE caches regenerate automatically and are deleted directly;
CAREFUL items are moved to the Trash so they stay recoverable.

EXAMPLES:
  csys clean          Interactive TUI: pick what to clean
  csys clean --safe   Clean the whole safe tier, no prompt
  csys clean -n       Dry run: only show what is reclaimable
  csys clean --json   Machine-readable output
  csys clean --nuke   Delete permanently instead of trashing
  csys clean --trash  Move even safe caches to Trash`

	CleanModulesShort = "Find stale project build artifacts (node_modules, target, …)"
	CleanModulesLong  = `Walk your project folders for build artifacts — node_modules, cargo target,
.venv, Pods, .next, .turbo — and show them with each project's last activity.
Artifacts are moved to the Trash and restore with a single install command.

EXAMPLES:
  csys clean modules                 Scan ~/Documents, ~/dev, ~/code, ~/projects
  csys clean modules -p ~/work      Scan a specific root
  csys clean modules --older 3m     Only projects untouched for 3+ months`
)
