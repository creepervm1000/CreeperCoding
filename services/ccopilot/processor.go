package ccopilot

import (
	"context"
	"fmt"
	"strings"
	"time"

	issues_model "creepercoding.dev/models/issues"
	"creepercoding.dev/models/perm/access"
	perm_model "creepercoding.dev/models/perm"
	"creepercoding.dev/models/unit"
	repo_model "creepercoding.dev/models/repo"
	user_model "creepercoding.dev/models/user"
	"creepercoding.dev/modules/gitrepo"
	"creepercoding.dev/modules/log"
	"creepercoding.dev/modules/setting"
	issue_service "creepercoding.dev/services/issue"
)

func processMention(task *mentionTask) error {
	ctx := context.Background()

	repo, err := repo_model.GetRepositoryByID(ctx, task.RepoID)
	if err != nil {
		return fmt.Errorf("get repo: %w", err)
	}

	issue, err := issues_model.GetIssueByID(ctx, task.IssueID)
	if err != nil {
		return fmt.Errorf("get issue: %w", err)
	}
	if err := issue.LoadRepo(ctx); err != nil {
		return fmt.Errorf("load repo: %w", err)
	}
	if err := issue.LoadPoster(ctx); err != nil {
		return fmt.Errorf("load poster: %w", err)
	}

	doer, err := user_model.GetUserByID(ctx, task.DoerID)
	if err != nil {
		return fmt.Errorf("get doer: %w", err)
	}

	ccopilotUser := getCcopilotUser(ctx)

	triggerComment, _ := issues_model.GetCommentByID(ctx, task.CommentID)
	userContent := ""
	if triggerComment != nil {
		userContent = triggerComment.Content
	}

	if wantsCommitMessage(userContent) {
		suggestCommitMessage(ctx, doer, repo, issue, ccopilotUser, userContent)
		return nil
	}

	if issue.IsPull {
		reviewPullRequest(ctx, doer, repo, issue, ccopilotUser, userContent)
		return nil
	}

	perm, _ := access.GetDoerRepoPermission(ctx, repo, doer)
	canWrite := perm.CanAccess(perm_model.AccessModeWrite, unit.TypeCode)

	systemPrompt := fmt.Sprintf(`You are CreeperCoding Copilot (@ccopilot), an AI assistant integrated into the CreeperCoding Git service running on the %q instance (%s). You help users with code reviews, debugging, and implementing changes.

You are currently responding to a user on an issue in this CreeperCoding instance. ANY text you generate after your tool calls will be posted as a reply to this issue. The user will see your response directly in the issue thread.

You have access to tools to read files, list directories, search code, and get diffs. When you have write permission, you can also write files, create branches, and open pull requests.

IMPORTANT RULES:
- Use the available tools to actually perform tasks, not just describe what you would do
- When asked to make code changes, write the files and create a pull request
- Always reference specific files and line numbers when discussing code
- If you cannot make changes due to permissions, explain clearly what the user needs to do
- Be concise but thorough
- Format code blocks with proper language tags
- After you finish making tool calls, you MUST include a text reply summarizing what you did. The user will see your final message as a comment on the issue.
- Do NOT tell the user to "open the issue in your GitHub UI" or "add a comment yourself" — this IS a CreeperCoding issue, and your reply will automatically be posted as a comment.

Repository: %s
Default branch: %s
Instance: %s`, setting.AppName, setting.AppURL, repo.FullName(), repo.DefaultBranch, setting.AppURL)

	// Build issue context: opener, participants, recent comments
	var issueCtx strings.Builder

	issueCtx.WriteString(fmt.Sprintf("## Issue #%d: %s\n", issue.Index, issue.Title))
	issueCtx.WriteString(fmt.Sprintf("**Opened by:** %s (@%s)\n\n", issue.Poster.DisplayName(), issue.Poster.Name))

	participantIDs, _ := issues_model.GetParticipantsIDsByIssueID(ctx, issue.ID)
	if len(participantIDs) > 0 {
		var participantNames []string
		for _, pid := range participantIDs {
			if pu, err := user_model.GetUserByID(ctx, pid); err == nil {
				participantNames = append(participantNames, fmt.Sprintf("%s (@%s)", pu.DisplayName(), pu.Name))
			}
		}
		if len(participantNames) > 0 {
			issueCtx.WriteString(fmt.Sprintf("**Participants:** %s\n\n", strings.Join(participantNames, ", ")))
		}
	}

	if issue.Content != "" {
		issueCtx.WriteString(fmt.Sprintf("**Description:**\n%s\n\n", truncateString(issue.Content, 2000)))
	}

	comments, _ := issues_model.FindComments(ctx, &issues_model.FindCommentsOptions{
		IssueID: issue.ID,
		Type:    issues_model.CommentTypeComment,
	})
	if len(comments) > 0 {
		issueCtx.WriteString("**Comments:**\n")
		start := 0
		if len(comments) > 20 {
			start = len(comments) - 20
		}
		for i := start; i < len(comments); i++ {
			c := comments[i]
			if c.Poster == nil {
				_ = c.LoadPoster(ctx)
			}
			posterName := "unknown"
			if c.Poster != nil {
				posterName = fmt.Sprintf("%s (@%s)", c.Poster.DisplayName(), c.Poster.Name)
			}
			commentPreview := truncateString(c.Content, 500)
			issueCtx.WriteString(fmt.Sprintf("- %s: %s\n", posterName, commentPreview))
		}
	}

	userPrompt := fmt.Sprintf("%s\n\nThe user just said: %s", issueCtx.String(), userContent)

	if canWrite {
		userPrompt += "\n\nYou have permission to make code changes. If changes are requested, create a branch, make the changes, and open a pull request."
	}

	branchName := fmt.Sprintf("ccopilot-%s-%d", time.Now().Format("20060102"), issue.Index)
	userPrompt += fmt.Sprintf("\n\nWhen creating branches for code changes, use: %s", branchName)

	messages := []chatMessage{{Role: "user", Content: userPrompt}}
	reply, _, err := agentChat(ctx, repo, doer, messages, systemPrompt)
	if err != nil {
		log.Error("ccopilot: AI query failed: %v", err)
		postReply(ctx, ccopilotUser, repo, issue, fmt.Sprintf("I encountered an error contacting the AI: %v", err))
		return err
	}

	postReply(ctx, ccopilotUser, repo, issue, reply)
	return nil
}

func postReply(ctx context.Context, poster *user_model.User, repo *repo_model.Repository, issue *issues_model.Issue, content string) {
	if len(content) > 10000 {
		content = content[:10000] + "\n\n*(response truncated)*"
	}

	comment, err := issue_service.CreateIssueComment(ctx, poster, repo, issue, content, nil)
	if err != nil {
		log.Error("ccopilot: failed to post reply: %v", err)
	}
	_ = comment
}

func getCcopilotUser(ctx context.Context) *user_model.User {
	u, err := user_model.GetUserByName(ctx, user_model.CcopilotUserName)
	if err != nil {
		return user_model.NewCcopilotUser()
	}
	return u
}

func wantsCommitMessage(content string) bool {
	lower := strings.ToLower(content)
	for _, phrase := range []string{
		"commit message",
		"suggest commit",
		"suggest a commit message",
		"suggestcommitmessage",
		"generate commit message",
		"write a commit message",
		"propose a commit message",
	} {
		if strings.Contains(lower, phrase) {
			return true
		}
	}
	return false
}

func loadPRDiff(ctx context.Context, issue *issues_model.Issue) (pr *issues_model.PullRequest, diffStr string, filesChanged []string, err error) {
	if !issue.IsPull {
		return nil, "", nil, fmt.Errorf("issue is not a pull request")
	}
	if err := issue.LoadPullRequest(ctx); err != nil {
		return nil, "", nil, fmt.Errorf("load pull request: %w", err)
	}
	pr = issue.PullRequest
	if pr == nil {
		return nil, "", nil, fmt.Errorf("pull request not found")
	}
	if err := pr.LoadBaseRepo(ctx); err != nil {
		return nil, "", nil, fmt.Errorf("load base repo: %w", err)
	}
	gitRepo, closer, err := gitrepo.RepositoryFromContextOrOpen(ctx, pr.BaseRepo)
	if err != nil {
		return nil, "", nil, fmt.Errorf("open repository: %w", err)
	}
	defer closer.Close()
	if pr.MergeBase == "" {
		return nil, "", nil, fmt.Errorf("merge base not computed yet")
	}
	compareArg := pr.MergeBase + "..." + pr.GetGitHeadRefName()
	filesChanged, err = gitRepo.GetFilesChangedBetween(pr.MergeBase, pr.GetGitHeadRefName())
	if err != nil {
		log.Error("ccopilot: GetFilesChangedBetween: %v", err)
		filesChanged = []string{"<error listing files>"}
	}
	var diffBuf strings.Builder
	if err := gitRepo.GetDiff(compareArg, &diffBuf); err != nil {
		return nil, "", nil, fmt.Errorf("get diff: %w", err)
	}
	diffStr = diffBuf.String()
	if len(diffStr) > 8000 {
		diffStr = diffStr[:8000] + "\n...(diff truncated)"
	}
	return pr, diffStr, filesChanged, nil
}

func suggestCommitMessage(ctx context.Context, doer *user_model.User, repo *repo_model.Repository, issue *issues_model.Issue, ccopilotUser *user_model.User, userContent string) {
	pr, diffStr, filesChanged, err := loadPRDiff(ctx, issue)
	if err != nil {
		postReply(ctx, ccopilotUser, repo, issue, fmt.Sprintf("Cannot suggest a commit message: %v", err))
		return
	}

	var prompt strings.Builder
	prompt.WriteString(fmt.Sprintf("## Repository: %s\n", repo.FullName()))
	prompt.WriteString(fmt.Sprintf("Branch: `%s` → `%s`\n", pr.HeadBranch, pr.BaseBranch))
	prompt.WriteString(fmt.Sprintf("Files changed (%d):\n", len(filesChanged)))
	for _, f := range filesChanged {
		prompt.WriteString(fmt.Sprintf("  - %s\n", f))
	}
	prompt.WriteString(fmt.Sprintf("\n## Diff:\n%s\n", diffStr))
	prompt.WriteString("\nBased on the above changes, suggest a concise commit message using Conventional Commits format (e.g. `feat(scope): message`). Explain your reasoning briefly.")

	systemPrompt := `You are CreeperCoding Copilot (@ccopilot). Your task is to analyze a git diff and suggest a good commit message.

Rules:
- Use Conventional Commits format: type(scope): description
- Types: feat, fix, refactor, test, docs, style, chore, perf, ci, build, revert
- Keep the subject line under 72 characters
- Provide a brief body explaining the motivation if needed
- Be concise and specific`

	response, err := queryAI(ctx, systemPrompt, prompt.String())
	if err != nil {
		log.Error("ccopilot: AI query failed for commit message: %v", err)
		postReply(ctx, ccopilotUser, repo, issue, fmt.Sprintf("I encountered an error contacting the AI: %v", err))
		return
	}

	postReply(ctx, ccopilotUser, repo, issue, response)
}

func reviewPullRequest(ctx context.Context, doer *user_model.User, repo *repo_model.Repository, issue *issues_model.Issue, ccopilotUser *user_model.User, userContent string) {
	pr, diffStr, filesChanged, err := loadPRDiff(ctx, issue)
	if err != nil {
		postReply(ctx, ccopilotUser, repo, issue, fmt.Sprintf("Cannot review pull request: %v", err))
		return
	}

	var prompt strings.Builder
	prompt.WriteString(fmt.Sprintf("## Repository: %s\n", repo.FullName()))
	prompt.WriteString(fmt.Sprintf("Reviewed by: %s (@%s)\n", doer.DisplayName(), doer.Name))
	prompt.WriteString(fmt.Sprintf("Branch: `%s` → `%s`\n", pr.HeadBranch, pr.BaseBranch))
	prompt.WriteString(fmt.Sprintf("PR title: %s\n", issue.Title))
	prompt.WriteString(fmt.Sprintf("PR description: %s\n", issue.Content))
	prompt.WriteString(fmt.Sprintf("Files changed (%d):\n", len(filesChanged)))
	for _, f := range filesChanged {
		prompt.WriteString(fmt.Sprintf("  - %s\n", f))
	}
	prompt.WriteString(fmt.Sprintf("\n## Diff:\n%s\n", diffStr))
	if userContent != "" {
		userMsg := strings.TrimSpace(strings.Replace(userContent, "@"+user_model.CcopilotUserName, "", 1))
		if userMsg != "" {
			prompt.WriteString(fmt.Sprintf("\n## User comment:\n%s\n", userMsg))
		}
	}

	systemPrompt := `You are CreeperCoding Copilot (@ccopilot), a code review assistant. Review the following pull request and provide:

1. **Rating**: Score out of 10
2. **Merge recommendation**: Yes / No / With changes
3. **Summary**: 2-3 sentence overview of what the PR does
4. **Strengths**: What the PR does well
5. **Issues**: Specific problems, bugs, or concerns with file/line references
6. **Suggestions**: Concrete improvements with code examples where appropriate

Be constructive, specific, and thorough. Reference exact file paths and line numbers from the diff.`

	response, err := queryAI(ctx, systemPrompt, prompt.String())
	if err != nil {
		log.Error("ccopilot: AI query failed for PR review: %v", err)
		postReply(ctx, ccopilotUser, repo, issue, fmt.Sprintf("I encountered an error reviewing the PR: %v", err))
		return
	}

	postReply(ctx, ccopilotUser, repo, issue, response)
}
