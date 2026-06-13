package ccopilot

import (
	"context"
	"fmt"
	"sync"

	repo_model "creepercoding.dev/models/repo"
	user_model "creepercoding.dev/models/user"
	"creepercoding.dev/modules/log"
	"creepercoding.dev/modules/setting"
)

type agentSession struct {
	messages []chatMessage
}

var (
	agentSessionsMu sync.Mutex
	agentSessions   = map[int64]*agentSession{}
)

func getOrCreateSession(repoID int64) *agentSession {
	agentSessionsMu.Lock()
	defer agentSessionsMu.Unlock()
	s, ok := agentSessions[repoID]
	if !ok {
		s = &agentSession{}
		agentSessions[repoID] = s
	}
	return s
}

func AgentChat(ctx context.Context, repo *repo_model.Repository, doer *user_model.User, message string) (string, error) {
	session := getOrCreateSession(repo.ID)

	repoInfo := fmt.Sprintf("Repository: %s\nDefault branch: %s\nDescription: %s",
		repo.FullName(), repo.DefaultBranch, repo.Description)

	systemPrompt := fmt.Sprintf(`You are ccopilot, an AI coding assistant for CreeperCoding (%s).

You are working on a CreeperCoding repository. You have the ability to read and edit files, create branches, and open pull requests using the tools provided to you.

This is the CreeperCoding Agent tab — a chat interface on a repository page. Every message you send will be displayed directly to the user in this chat.

When the user asks you to make changes:
1. Use the tools to explore and understand the codebase first
2. Make changes by writing files to a new branch
3. When ready, create a pull request with your changes
4. After you finish making changes, ALWAYS include a text reply summarizing what you did

Current context:
%s

Instance: %s`, setting.AppURL, repoInfo, setting.AppURL)

	session.messages = append(session.messages, chatMessage{Role: "user", Content: message})

	reply, allMessages, err := agentChat(ctx, repo, doer, session.messages, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("agent chat failed: %w", err)
	}

	// Update session: keep only the user+assistant messages (not system or tool results)
	session.messages = filterSessionMessages(allMessages)
	log.Debug("ccopilot: agent session %d now has %d messages", repo.ID, len(session.messages))

	return reply, nil
}

func filterSessionMessages(messages []chatMessage) []chatMessage {
	var result []chatMessage
	for _, m := range messages {
		if m.Role == "system" || m.Role == "tool" {
			continue
		}
		result = append(result, m)
	}
	return result
}
