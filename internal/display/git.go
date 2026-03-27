package display

const (
	GSyncShort = "Hard reset to origin branch"
	GSyncLong  = `Fetch and hard reset to origin branch.

EXAMPLES:
  cs gsync           Reset to origin/main
  cs gsync develop   Reset to origin/develop`

	GCleanShort = "Delete all local branches except current"
	GCleanLong  = `Delete all local branches except the current one.

EXAMPLES:
  cs gclean          Delete with confirmation
  cs gclean --force  Delete without confirmation`

	GSoftShort = "Soft reset HEAD~n for re-committing"
	GSoftLong  = `Soft reset to undo commits while keeping changes staged.

EXAMPLES:
  cs gsoft       Soft reset HEAD~1
  cs gsoft 3     Soft reset HEAD~3`

	GAcShort = "Add all changes and commit"
	GAcLong  = `Stage all changes and commit with message.

EXAMPLES:
  cs gac "fix bug"
  cs gac "add feature"`

	GCoShort = "Switch to an existing branch"
	GCoLong  = `Switch to an existing branch.

EXAMPLES:
  cs gco main             Switch to main
  cs gco feature/auth     Switch to feature/auth`

	GPushShort = "Push to remote"
	GPushLong  = `Push current branch to remote.

EXAMPLES:
  cs gpush`

	GPullShort = "Pull latest changes"
	GPullLong  = `Pull latest changes from remote.

EXAMPLES:
  cs gpull`

	GFpShort = "Force push (with lease)"
	GFpLong  = `Force push current branch using --force-with-lease for safety.
This prevents overwriting others' work while allowing force push after rebase/amend.

EXAMPLES:
  cs gfp`

	GUndoShort = "Undo last commit, keep changes staged"
	GUndoLong  = `Undo the last commit but keep all changes staged.

EXAMPLES:
  cs gundo`

	GWipShort = "Quick WIP commit"
	GWipLong  = `Stage all changes and create a WIP commit.

EXAMPLES:
  cs gwip`

	GAmendShort = "Amend last commit (optionally change message)"
	GAmendLong  = `Stage all changes and amend the last commit.
Without arguments, keeps the existing message.
With a message argument, updates the commit message.

EXAMPLES:
  cs gamend              Amend without changing message
  cs gamend "new msg"    Amend with new message`

	GRbShort = "Rebase branch onto base branch"
	GRbLong  = `Rebase current or specified branch onto a base branch.

EXAMPLES:
  cs grb                     Rebase current branch onto origin/main
  cs grb develop             Rebase current branch onto origin/develop
  cs grb main feature/auth   Checkout feature/auth, rebase onto origin/main`

	GLogShort = "Show commit log with graph"
	GLogLong  = `Show a formatted commit log.

EXAMPLES:
  cs glog        Show last 10 commits
  cs glog 20     Show last 20 commits`

	GStShort = "Show git status"
	GStLong  = `Show working tree status (short format).

EXAMPLES:
  cs gst`

	GCbShort = "Create and switch to a new branch"
	GCbLong  = `Create a new branch and switch to it immediately.

EXAMPLES:
  cs gcb feature/auth      Create and switch to feature/auth
  cs gcb fix/login-bug     Create and switch to fix/login-bug`

	GBrnShort = "Rename current branch"
	GBrnLong  = `Rename the current branch to a new name (force rename).

EXAMPLES:
  cs gbrn feature/new-name    Rename current branch
  cs gbrn main                Rename current branch to main`
)
