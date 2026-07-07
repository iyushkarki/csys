package cleaners

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type HistoryEntry struct {
	Time    time.Time `json:"time"`
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Size    uint64    `json:"size"`
	Trashed bool      `json:"trashed"`
}

func historyPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".csys", "history.jsonl")
}

func NewHistoryEntry(t *Target, opts Options) HistoryEntry {
	return HistoryEntry{
		Time:    time.Now(),
		ID:      t.ID,
		Name:    t.Name,
		Size:    t.Size,
		Trashed: t.UsesTrash(opts),
	}
}

func AppendHistory(entries []HistoryEntry) {
	path := historyPath()
	if path == "" || len(entries) == 0 {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range entries {
		_ = enc.Encode(e)
	}
}

func LifetimeFreed() uint64 {
	f, err := os.Open(historyPath())
	if err != nil {
		return 0
	}
	defer f.Close()

	var total uint64
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var e HistoryEntry
		if json.Unmarshal(scanner.Bytes(), &e) == nil {
			total += e.Size
		}
	}
	return total
}
