package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func costFilepath() string {
	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		dir = filepath.Join(homeDir, ".local", "share")
	}
	return filepath.Join(dir, "kommit", "cost.json")
}

type Cost struct {
	RepoName string
	Cost     float64
}

func UpdateCost(cost float64) error {
	costFilePath := costFilepath()
	if costFilePath == "" {
		return fmt.Errorf("could not determine cost file path")
	}

	err := os.MkdirAll(filepath.Dir(costFilePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	repoName, err := getRepoName()
	if err != nil {
		return fmt.Errorf("failed to get repository name: %w", err)
	}

	costs := make(map[string]float64)
	costData, err := os.ReadFile(costFilePath)
	if err == nil && len(costData) > 0 {
		if err := json.Unmarshal(costData, &costs); err != nil {
			costs = make(map[string]float64)
		}
	}
	costs[repoName] += cost

	updatedData, err := json.MarshalIndent(costs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal costs: %w", err)
	}

	if err := os.WriteFile(costFilePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write cost file: %w", err)
	}

	return nil
}

func getRepoName() (string, error) {
	fullPath, err := ExecGit("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	fullPath = strings.TrimSpace(fullPath)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fullPath, nil
	}

	if strings.HasPrefix(fullPath, homeDir) {
		return strings.Replace(fullPath, homeDir, "~", 1), nil
	}
	return fullPath, nil
}
