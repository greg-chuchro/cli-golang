package completion_providers

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type DirectoryCompletionProvider struct{}

func NewDirectoryCompletionProvider() *DirectoryCompletionProvider {
	return &DirectoryCompletionProvider{}
}

func (p *DirectoryCompletionProvider) GetCompletions(ctx context.Context, completionContext ICompletionProviderContext) []string {
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

	directories := []string{}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return []string{}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(searchPrefix)) {
			filter := completionContext.Filter()
			if filter == "" || matchesGlob(name, filter) {
				directories = append(directories, name)
			}
		}
	}

	sort.Strings(directories)
	return directories
}

func matchesGlob(name string, pattern string) bool {
	if pattern == "*" {
		return true
	}

	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		middle := strings.Trim(pattern, "*")
		return strings.Contains(strings.ToLower(name), strings.ToLower(middle))
	}

	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(strings.ToLower(name), strings.ToLower(suffix))
	}

	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix))
	}

	return strings.EqualFold(name, pattern)
}
