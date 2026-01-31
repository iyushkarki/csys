package display

const (
	GShort = "Git shortcuts for common workflows"
	GLong  = `Git shortcuts for common workflows.

Usage:
  csys g [command]

Available Commands:
  sync [branch]   Hard reset to origin (default: main)
  clean           Delete all local branches except current
  soft [n]        Soft reset HEAD~n (default: 1)
  ac "msg"        Add all + commit
  acp "msg"       Add + commit + push
  undo            Undo last commit (keep changes staged)
  wip             Quick WIP commit
  amend           Amend last commit without editing message
  rb [base] [branch]  Rebase branch onto base (default: main)`

	GSyncShort = "Hard reset to origin branch"
	GSyncLong  = `Fetch and hard reset to origin branch.

EXAMPLES:
  csys g sync           Reset to origin/main
  csys g sync develop   Reset to origin/develop`

	GCleanShort = "Delete all local branches except current"
	GCleanLong  = `Delete all local branches except the current one.

EXAMPLES:
  csys g clean          Delete with confirmation
  csys g clean --force  Delete without confirmation`

	GSoftShort = "Soft reset HEAD~n for re-committing"
	GSoftLong  = `Soft reset to undo commits while keeping changes staged.

EXAMPLES:
  csys g soft       Soft reset HEAD~1
  csys g soft 3     Soft reset HEAD~3`

	GAcShort = "Add all changes and commit"
	GAcLong  = `Stage all changes and commit with message.

EXAMPLES:
  csys g ac "fix bug"
  csys g ac "add feature"`

	GAcpShort = "Add, commit, and push"
	GAcpLong  = `Stage all changes, commit, and push to remote.

EXAMPLES:
  csys g acp "fix bug"
  csys g acp "add feature"`

	GUndoShort = "Undo last commit, keep changes staged"
	GUndoLong  = `Undo the last commit but keep all changes staged.

EXAMPLES:
  csys g undo`

	GWipShort = "Quick WIP commit"
	GWipLong  = `Stage all changes and create a WIP commit.

EXAMPLES:
  csys g wip`

	GAmendShort = "Amend last commit without editing message"
	GAmendLong  = `Stage all changes and amend the last commit.

EXAMPLES:
  csys g amend`

	GRbShort = "Rebase branch onto base branch"
	GRbLong  = `Rebase current or specified branch onto a base branch.

EXAMPLES:
  csys g rb                     Rebase current branch onto origin/main
  csys g rb develop             Rebase current branch onto origin/develop
  csys g rb main feature/auth   Checkout feature/auth, rebase onto origin/main`

	GLogShort = "Show commit log with graph"
	GLogLong  = `Show a formatted commit log.

EXAMPLES:
  csys g log        Show last 10 commits
  csys g log 20     Show last 20 commits`

	GStShort = "Show git status"
	GStLong  = `Show working tree status (short format).

EXAMPLES:
  csys g st`
)
