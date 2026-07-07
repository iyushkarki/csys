package cleaners

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func (t *Target) Clean(opts Options) error {
	if t.cleanFn != nil {
		return t.cleanFn(t, opts)
	}

	toTrash := t.UsesTrash(opts)
	var firstErr error
	for _, p := range t.paths {
		var err error
		switch {
		case t.keepDir:
			err = removeContents(p, toTrash)
		case toTrash:
			err = moveToTrash(p)
		default:
			err = os.RemoveAll(p)
		}
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if len(t.postClean) > 0 && firstErr == nil {
		if _, err := exec.LookPath(t.postClean[0]); err == nil {
			_ = exec.Command(t.postClean[0], t.postClean[1:]...).Run()
		}
	}
	return firstErr
}

func moveToTrash(path string) error {
	if runtime.GOOS == "darwin" {
		return macTrash(path)
	}
	return xdgTrash(path)
}

func macTrash(path string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dest := trashDest(filepath.Join(home, ".Trash"), path)
	if err := os.Rename(path, dest); err == nil {
		return nil
	}
	if trash, err := exec.LookPath("trash"); err == nil {
		if out, err := exec.Command(trash, path).CombinedOutput(); err != nil {
			return fmt.Errorf("trash: %s", strings.TrimSpace(string(out)))
		}
		return nil
	}
	return fmt.Errorf("could not move %s to Trash", filepath.Base(path))
}

func xdgTrash(path string) error {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dataHome = filepath.Join(home, ".local", "share")
	}
	files := filepath.Join(dataHome, "Trash", "files")
	info := filepath.Join(dataHome, "Trash", "info")
	if err := os.MkdirAll(files, 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(info, 0o700); err != nil {
		return err
	}

	dest := trashDest(files, path)
	infoPath := filepath.Join(info, filepath.Base(dest)+".trashinfo")
	contents := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n",
		(&url.URL{Path: path}).EscapedPath(), time.Now().Format("2006-01-02T15:04:05"))
	if err := os.WriteFile(infoPath, []byte(contents), 0o600); err != nil {
		return err
	}
	if err := os.Rename(path, dest); err != nil {
		_ = os.Remove(infoPath)
		return fmt.Errorf("could not move %s to trash: %w", filepath.Base(path), err)
	}
	return nil
}

func trashDest(trashDir, path string) string {
	dest := filepath.Join(trashDir, filepath.Base(path))
	if _, err := os.Lstat(dest); err == nil {
		dest = fmt.Sprintf("%s %s", dest, time.Now().Format("15.04.05"))
	}
	return dest
}

func removeContents(dir string, toTrash bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var firstErr error
	for _, e := range entries {
		p := filepath.Join(dir, e.Name())
		var err error
		if toTrash {
			err = moveToTrash(p)
		} else {
			err = os.RemoveAll(p)
		}
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
