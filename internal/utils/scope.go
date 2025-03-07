package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func GetScopesFromHistory() ([]string, error) {
	output, err := ExecGit("log", "--pretty=format:%s")
	if err != nil {
		return nil, err
	}

	scopeRegex := regexp.MustCompile(`^[\w]+\(([\w-]+)\)(?:!)?:`)
	scopesMap := make(map[string]bool)
	commitLines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range commitLines {
		if line == "" {
			continue
		}

		matches := scopeRegex.FindStringSubmatch(line)
		if len(matches) > 1 && matches[1] != "" {
			scopesMap[matches[1]] = true
		}
	}

	if len(scopesMap) == 0 {
		return nil, nil
	}

	scopes := make([]string, 0, len(scopesMap))
	for scope := range scopesMap {
		scopes = append(scopes, scope)
	}

	sort.Strings(scopes)
	return scopes, nil
}

func GetFilesFromDirectory(maxDepth int) ([]string, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	skipDirs := map[string]bool{
		"node_modules": true,
		"dist":         true,
		"build":        true,
		"bin":          true,
		"vendor":       true,
		"target":       true,
		".git":         true,
		".github":      true,
		".vscode":      true,
		".idea":        true,
		"coverage":     true,
		"tmp":          true,
		"temp":         true,
	}

	// Store all file paths
	var filePaths []string

	// Walk through all directories recursively
	err = filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path from the project root
		relPath, err := filepath.Rel(path, currentPath)
		if err != nil {
			return err
		}

		// Skip the root directory
		if relPath == "." {
			return nil
		}

		// Check if this directory or any of its parent directories should be skipped
		pathParts := strings.Split(relPath, string(os.PathSeparator))
		for _, part := range pathParts {
			if skipDirs[part] {
				return filepath.SkipDir
			}
		}

		// Check if we've reached the maximum depth
		if maxDepth > 0 && len(pathParts) > maxDepth {
			return nil
		}

		ignored, err := isGitIgnored(relPath)
		if err == nil && ignored {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Only include files, not directories
		if !info.IsDir() {
			filePaths = append(filePaths, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(filePaths)
	return filePaths, nil
}

func isGitIgnored(path string) (bool, error) {
	_, err := ExecGit("check-ignore", "-q", path)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil // Exit code 1 means the file is not ignored
			}
		}
		return false, err
	}

	return true, nil // Exit code 0 means the file is ignored
}
