package ccopilot

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	repo_model "creepercoding.dev/models/repo"
	"creepercoding.dev/modules/gitrepo"
	"creepercoding.dev/modules/log"
	"creepercoding.dev/modules/setting"
)

func GenerateCommitMessage(ctx context.Context, repo *repo_model.Repository, path, content, branch, commitSummary string) (string, error) {
	cfg := setting.Config().Ccopilot
	if !cfg.Enabled.Value(ctx) {
		return "", fmt.Errorf("ccopilot is not enabled")
	}

	gitRepo, closer, err := gitrepo.RepositoryFromContextOrOpen(ctx, repo)
	if err != nil {
		return "", fmt.Errorf("open repo: %w", err)
	}
	defer closer.Close()

	var originalContent string
	if branch != "" {
		commit, err := gitRepo.GetBranchCommit(branch)
		if err == nil {
			orig, err := commit.GetFileContent(path, 2000)
			if err == nil {
				originalContent = orig
			}
		}
	}

	prompt := "Generate a concise, descriptive Git commit message for the following change.\n"
	if commitSummary != "" {
		prompt += fmt.Sprintf("User's intent: %s\n", commitSummary)
	}
	prompt += fmt.Sprintf("File: %s\n", path)
	if originalContent != "" {
		prompt += fmt.Sprintf("Original content (first 2000 chars):\n%s\n\n", truncateString(originalContent, 2000))
	}
	prompt += fmt.Sprintf("New content (first 2000 chars):\n%s\n", truncateString(content, 2000))
	prompt += "\nRespond with ONLY the commit message (subject line, optionally followed by a blank line and body). Keep the subject under 72 characters."

	messages := []chatMessage{
		{Role: "system", Content: "You are a helpful assistant that generates concise, descriptive Git commit messages."},
		{Role: "user", Content: prompt},
	}

	respJSON, err := queryAIMessages(ctx, messages, nil)
	if err != nil {
		return "", fmt.Errorf("AI query: %w", err)
	}

	var choice chatResponseChoice
	if err := json.Unmarshal([]byte(respJSON), &choice); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	msg := strings.TrimSpace(choice.Message.Content)
	log.Debug("ccopilot generated commit message: %s", msg)
	return msg, nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
