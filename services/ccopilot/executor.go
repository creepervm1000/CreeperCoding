package ccopilot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"creepercoding.dev/models/perm/access"
	perm_model "creepercoding.dev/models/perm"
	"creepercoding.dev/models/unit"
	repo_model "creepercoding.dev/models/repo"
	user_model "creepercoding.dev/models/user"
	"creepercoding.dev/modules/git"
	"creepercoding.dev/modules/gitrepo"
	"creepercoding.dev/modules/log"
	files_service "creepercoding.dev/services/repository/files"
	repo_service "creepercoding.dev/services/repository"
	"creepercoding.dev/services/pull"
	issues_model "creepercoding.dev/models/issues"
	"creepercoding.dev/modules/setting"
)

const maxToolIterations = 15

func agentChat(ctx context.Context, repo *repo_model.Repository, doer *user_model.User, messages []chatMessage, systemPrompt string) (string, []chatMessage, error) {
	perm, _ := access.GetDoerRepoPermission(ctx, repo, doer)
	canWrite := perm.CanAccess(perm_model.AccessModeWrite, unit.TypeCode)

	allMessages := make([]chatMessage, 0, len(messages)+2)
	allMessages = append(allMessages, chatMessage{Role: "system", Content: systemPrompt})
	allMessages = append(allMessages, messages...)

	tools := toolsDefinitions()

	for iteration := 0; iteration < maxToolIterations; iteration++ {
		rawJSON, err := queryAIMessages(ctx, allMessages, tools)
		if err != nil {
			return "", nil, fmt.Errorf("AI query failed at iteration %d: %w", iteration, err)
		}

		var choice chatResponseChoice
		if err := json.Unmarshal([]byte(rawJSON), &choice); err != nil {
			return "", nil, fmt.Errorf("failed to parse AI response: %w", err)
		}

		assistantMsg := choice.Message
		allMessages = append(allMessages, assistantMsg)

		if choice.FinishReason == "stop" || (choice.FinishReason == "" && len(assistantMsg.ToolCalls) == 0) {
			return assistantMsg.Content, allMessages, nil
		}

		if choice.FinishReason == "tool_calls" && len(assistantMsg.ToolCalls) > 0 {
			for _, tc := range assistantMsg.ToolCalls {
				result, err := executeToolCall(ctx, tc, repo, doer, canWrite)
				if err != nil {
					result = fmt.Sprintf("Error executing %s: %v", tc.Function.Name, err)
					log.Error("ccopilot: tool %s failed: %v", tc.Function.Name, err)
				}
				allMessages = append(allMessages, chatMessage{
					Role:       "tool",
					ToolCallID: tc.ID,
					Content:    result,
				})
			}
			continue
		}

		log.Warn("ccopilot: unexpected finish_reason=%s at iteration %d", choice.FinishReason, iteration)
		return assistantMsg.Content, allMessages, nil
	}

	return "", nil, fmt.Errorf("tool loop exceeded max iterations (%d)", maxToolIterations)
}

func executeToolCall(ctx context.Context, tc toolCall, repo *repo_model.Repository, doer *user_model.User, canWrite bool) (string, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
		return "", fmt.Errorf("invalid arguments for %s: %w", tc.Function.Name, err)
	}

	switch tc.Function.Name {
	case "read_file":
		return toolReadFile(ctx, repo, args)
	case "list_directory":
		return toolListDirectory(ctx, repo, args)
	case "search_code":
		return toolSearchCode(ctx, repo, args)
	case "get_diff":
		return toolGetDiff(ctx, repo, args)
	case "write_file":
		if !canWrite {
			return "Permission denied: you do not have write access to this repository. Only read-only operations are allowed.", nil
		}
		return toolWriteFile(ctx, repo, doer, args)
	case "create_branch":
		if !canWrite {
			return "Permission denied: you do not have write access to this repository.", nil
		}
		return toolCreateBranch(ctx, repo, doer, args)
	case "create_pull_request":
		if !canWrite {
			return "Permission denied: you do not have write access to this repository.", nil
		}
		return toolCreatePullRequest(ctx, repo, doer, args)
	default:
		return "", fmt.Errorf("unknown tool: %s", tc.Function.Name)
	}
}

func toolReadFile(ctx context.Context, repo *repo_model.Repository, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}

	gitRepo, err := gitrepo.OpenRepository(ctx, repo)
	if err != nil {
		return "", fmt.Errorf("open repository: %w", err)
	}
	defer gitRepo.Close()

	commit, err := gitRepo.GetBranchCommit(repo.DefaultBranch)
	if err != nil {
		return "", fmt.Errorf("get default branch commit: %w", err)
	}

	entry, err := commit.Tree.GetTreeEntryByPath(path)
	if err != nil {
		return "", fmt.Errorf("file not found: %s: %w", path, err)
	}

	if entry.IsDir() {
		return "", fmt.Errorf("%s is a directory, use list_directory instead", path)
	}

	content, err := entry.Blob().GetBlobContent(50000)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	result := fmt.Sprintf("File: %s\n\n%s", path, content)
	return result, nil
}

func toolListDirectory(ctx context.Context, repo *repo_model.Repository, args map[string]any) (string, error) {
	path, _ := args["path"].(string)

	gitRepo, err := gitrepo.OpenRepository(ctx, repo)
	if err != nil {
		return "", fmt.Errorf("open repository: %w", err)
	}
	defer gitRepo.Close()

	commit, err := gitRepo.GetBranchCommit(repo.DefaultBranch)
	if err != nil {
		return "", fmt.Errorf("get default branch commit: %w", err)
	}

	var entries git.Entries
	if path == "" {
		entries, err = commit.Tree.ListEntries()
	} else {
		subTree, err := commit.Tree.SubTree(path)
		if err != nil {
			return "", fmt.Errorf("directory not found: %s", path)
		}
		entries, err = subTree.ListEntries()
	}
	if err != nil {
		return "", fmt.Errorf("list entries: %w", err)
	}

	var b strings.Builder
	dirStr := "root"
	if path != "" {
		dirStr = path
	}
	b.WriteString(fmt.Sprintf("Contents of '%s':\n", dirStr))
	for _, entry := range entries {
		entryType := "file"
		if entry.IsDir() {
			entryType = "dir"
		}
		b.WriteString(fmt.Sprintf("  [%s] %s\n", entryType, entry.Name()))
	}
	return b.String(), nil
}

func toolSearchCode(ctx context.Context, repo *repo_model.Repository, args map[string]any) (string, error) {
	query, _ := args["query"].(string)
	if query == "" {
		return "", fmt.Errorf("query is required")
	}

	gitRepo, err := gitrepo.OpenRepository(ctx, repo)
	if err != nil {
		return "", fmt.Errorf("open repository: %w", err)
	}
	defer gitRepo.Close()

	commit, err := gitRepo.GetBranchCommit(repo.DefaultBranch)
	if err != nil {
		return "", fmt.Errorf("get default branch commit: %w", err)
	}

	entries, err := commit.Tree.ListEntriesRecursiveFast()
	if err != nil {
		return "", fmt.Errorf("list entries: %w", err)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Search results for '%s':\n", query))
	count := 0
	for _, entry := range entries {
		if count >= 30 {
			b.WriteString(fmt.Sprintf("\n...(showing first 30 results)"))
			break
		}
		if entry.IsDir() {
			continue
		}
		reader, err := entry.Blob().DataAsync()
		if err != nil {
			continue
		}
		contentBytes, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			continue
		}
		content := string(contentBytes)
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if strings.Contains(line, query) {
				b.WriteString(fmt.Sprintf("%s:%d: %s\n", entry.Name(), i+1, strings.TrimSpace(line)))
				count++
				if count >= 30 {
					break
				}
			}
		}
	}

	if count == 0 {
		b.WriteString("No matches found.")
	}
	return b.String(), nil
}

func toolGetDiff(ctx context.Context, repo *repo_model.Repository, args map[string]any) (string, error) {
	base, _ := args["base"].(string)
	head, _ := args["head"].(string)
	if base == "" || head == "" {
		return "", fmt.Errorf("both base and head are required")
	}

	gitRepo, err := gitrepo.OpenRepository(ctx, repo)
	if err != nil {
		return "", fmt.Errorf("open repository: %w", err)
	}
	defer gitRepo.Close()

	compareArg := base + "..." + head
	var diffBuf bytes.Buffer
	if err := gitRepo.GetDiff(compareArg, &diffBuf); err != nil {
		return "", fmt.Errorf("get diff: %w", err)
	}

	diffStr := diffBuf.String()
	if len(diffStr) > 16000 {
		diffStr = diffStr[:16000] + "\n...(diff truncated)"
	}
	return fmt.Sprintf("Diff between %s...%s:\n\n```diff\n%s\n```", base, head, diffStr), nil
}

func toolWriteFile(ctx context.Context, repo *repo_model.Repository, doer *user_model.User, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	content, _ := args["content"].(string)
	branch, _ := args["branch"].(string)
	message, _ := args["message"].(string)

	if path == "" || content == "" || branch == "" || message == "" {
		return "", fmt.Errorf("path, content, branch, and message are all required")
	}

	if branch == repo.DefaultBranch {
		return "", fmt.Errorf("cannot write directly to default branch (%s). Create a new branch first.", repo.DefaultBranch)
	}

	// Create the branch if it doesn't exist
	if err := repo_service.CreateNewBranch(ctx, doer, repo, repo.DefaultBranch, branch); err != nil {
		log.Debug("ccopilot: creating branch %s (may already exist): %v", branch, err)
	}

	_, err := files_service.ChangeRepoFiles(ctx, repo, doer, &files_service.ChangeRepoFilesOptions{
		OldBranch: branch,
		NewBranch: branch,
		Message:   message,
		Files: []*files_service.ChangeRepoFile{
			{
				Operation:     "create",
				TreePath:      path,
				ContentReader: strings.NewReader(content),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return fmt.Sprintf("Successfully wrote file '%s' on branch '%s'.\nCommit message: %s", path, branch, message), nil
}

func toolCreateBranch(ctx context.Context, repo *repo_model.Repository, doer *user_model.User, args map[string]any) (string, error) {
	name, _ := args["name"].(string)
	if name == "" {
		return "", fmt.Errorf("branch name is required")
	}

	if err := repo_service.CreateNewBranch(ctx, doer, repo, repo.DefaultBranch, name); err != nil {
		return "", fmt.Errorf("create branch: %w", err)
	}

	return fmt.Sprintf("Branch '%s' created successfully from '%s'", name, repo.DefaultBranch), nil
}

func toolCreatePullRequest(ctx context.Context, repo *repo_model.Repository, doer *user_model.User, args map[string]any) (string, error) {
	title, _ := args["title"].(string)
	body, _ := args["body"].(string)
	head, _ := args["head"].(string)

	if title == "" || head == "" {
		return "", fmt.Errorf("title and head branch are required")
	}

	pr := &issues_model.PullRequest{
		HeadBranch: head,
		BaseBranch: repo.DefaultBranch,
	}
	prIssue := &issues_model.Issue{
		RepoID:   repo.ID,
		Title:    title,
		Content:  body,
		PosterID: doer.ID,
		IsPull:   true,
	}

	err := pull.NewPullRequest(ctx, &pull.NewPullRequestOptions{
		Repo:        repo,
		Issue:       prIssue,
		PullRequest: pr,
	})
	if err != nil {
		return "", fmt.Errorf("create pull request: %w", err)
	}

	prURL := fmt.Sprintf("%s/%s/pulls/%d", setting.AppURL, repo.FullName(), prIssue.Index)
	return fmt.Sprintf("Pull request created: [%s](%s)", prURL, prURL), nil
}
