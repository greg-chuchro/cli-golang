package completion_providers

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileSystemCompletionProvider struct{}

func NewFileSystemCompletionProvider() *FileSystemCompletionProvider {
	return &FileSystemCompletionProvider{}
}

func (p *FileSystemCompletionProvider) GetCompletions(ctx context.Context, completionContext ICompletionProviderContext) []string {
	partialPath := completionContext.PartialInput()
	if partialPath == "" {
		partialPath = ""
	}

	directory := ""
	searchPrefix := ""

	if partialPath == "" {
		directory = "."
		searchPrefix = ""
	} else {
		fullPath := filepath.Clean(partialPath)
		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			directory = fullPath
			searchPrefix = ""
		} else {
			directory = filepath.Dir(fullPath)
			searchPrefix = filepath.Base(partialPath)
		}
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return []string{}
	}

	results := []string{}

	entries, err := os.ReadDir(directory)
	if err == nil {
		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(searchPrefix)) {
				results = append(results, name)
			}
		}
	}

	patterns := []string{"*"}
	filter := completionContext.Filter()
	if filter != "" {
		patterns = strings.Split(filter, ";")
	}

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		entries, err := os.ReadDir(directory)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(searchPrefix)) {
				if matched, _ := filepath.Match(pattern, name); matched {
					results = append(results, name)
				}
			}
		}
	}

	sort.Strings(results)
	return uniqueStrings(results)
}
