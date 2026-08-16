package completion_providers

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileCompletionProvider struct{}

func NewFileCompletionProvider() *FileCompletionProvider {
	return &FileCompletionProvider{}
}

func (p *FileCompletionProvider) GetCompletions(ctx context.Context, completionContext ICompletionProviderContext) []string {
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

	patterns := []string{"*"}
	filter := completionContext.Filter()
	if filter != "" {
		patterns = strings.Split(filter, ";")
	}

	files := []string{}
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
					files = append(files, name)
				}
			}
		}
	}

	sort.Strings(files)
	return uniqueStrings(files)
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(input))
	for _, s := range input {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
