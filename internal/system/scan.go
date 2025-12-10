package system

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileItem struct {
	Name      string
	Path      string
	Size      int64
	IsDir     bool
	Extension string
}

type TypeBreakdown struct {
	Extension string
	Size      int64
	Count     int
}

type ScanResult struct {
	RootPath      string
	TotalSize     int64
	FileCount     int
	DirCount      int
	Items         []FileItem
	TypeBreakdown []TypeBreakdown
}

func ScanDirectory(path string) (*ScanResult, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	result := &ScanResult{
		RootPath: absPath,
		Items:    make([]FileItem, 0),
	}

	// Maps to aggregate sizes
	typeSizes := make(map[string]int64)
	typeCounts := make(map[string]int)

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		item := FileItem{
			Name:  entry.Name(),
			Path:  filepath.Join(absPath, entry.Name()),
			IsDir: entry.IsDir(),
		}

		if entry.IsDir() {
			size, _, err := getDirSize(item.Path)
			if err == nil {
				item.Size = size
				result.DirCount++
			}
		} else {
			item.Size = info.Size()
			item.Extension = strings.ToLower(filepath.Ext(item.Name))
			result.FileCount++

			ext := item.Extension
			if ext == "" {
				ext = "no-ext"
			}
			typeSizes[ext] += item.Size
			typeCounts[ext]++
		}

		result.Items = append(result.Items, item)
		result.TotalSize += item.Size
	}

	// Sort items by size (descending)
	sort.Slice(result.Items, func(i, j int) bool {
		return result.Items[i].Size > result.Items[j].Size
	})

	for ext, size := range typeSizes {
		result.TypeBreakdown = append(result.TypeBreakdown, TypeBreakdown{
			Extension: ext,
			Size:      size,
			Count:     typeCounts[ext],
		})
	}

	sort.Slice(result.TypeBreakdown, func(i, j int) bool {
		return result.TypeBreakdown[i].Size > result.TypeBreakdown[j].Size
	})

	return result, nil
}

func getDirSize(path string) (int64, int, error) {
	var size int64
	var count int
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
			count++
		}
		return nil
	})
	return size, count, err
}
