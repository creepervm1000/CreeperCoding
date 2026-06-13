// Copyright 2025 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package ccopilot

import (
	"context"
	"errors"
	"strings"

	issues_model "creepercoding.dev/models/issues"
	repo_model "creepercoding.dev/models/repo"
	user_model "creepercoding.dev/models/user"
	"creepercoding.dev/modules/graceful"
	"creepercoding.dev/modules/log"
	"creepercoding.dev/modules/queue"
	"creepercoding.dev/modules/setting"
	notify_service "creepercoding.dev/services/notify"
)

type ccopilotNotifier struct {
	notify_service.NullNotifier
}

func init() {
	notify_service.RegisterNotifier(&ccopilotNotifier{})
}

var _ notify_service.Notifier = &ccopilotNotifier{}

type mentionTask struct {
	RepoID    int64
	IssueID   int64
	CommentID int64
	DoerID    int64
}

var mentionQueue *queue.WorkerPoolQueue[*mentionTask]

func Init(ctx context.Context) error {
	_, err := user_model.GetUserByName(ctx, user_model.CcopilotUserName)
	if err != nil {
		if err := user_model.CreateUser(ctx, user_model.NewCcopilotUser(), nil); err != nil {
			return err
		}
	}

	mentionQueue = queue.CreateUniqueQueue[*mentionTask](
		graceful.GetManager().ShutdownContext(),
		"ccopilot_mention",
		handler,
	)
	if mentionQueue == nil {
		return errors.New("unable to create ccopilot_mention queue")
	}
	go graceful.GetManager().RunWithCancel(mentionQueue)

	return nil
}

func (n *ccopilotNotifier) CreateIssueComment(ctx context.Context, doer *user_model.User, repo *repo_model.Repository,
	issue *issues_model.Issue, comment *issues_model.Comment, _ []*user_model.User,
) {
	if doer.Name == user_model.CcopilotUserName {
		return
	}

	if !strings.Contains(comment.Content, "@"+user_model.CcopilotUserName) {
		return
	}

	if !setting.Config().Ccopilot.Enabled.Value(ctx) {
		return
	}

	if repo.CcopilotDisabled {
		return
	}

	ccfg := setting.Config().Ccopilot
	if ccfg.Endpoint.Value(ctx) == "" || ccfg.APIKey.Value(ctx) == "" || ccfg.ModelName.Value(ctx) == "" {
		return
	}

	log.Debug("ccopilot: queuing mention task for issue #%d in repo %s", issue.ID, repo.FullName())

	_ = mentionQueue.Push(&mentionTask{
		RepoID:    repo.ID,
		IssueID:   issue.ID,
		CommentID: comment.ID,
		DoerID:    doer.ID,
	})
}

func handler(items ...*mentionTask) []*mentionTask {
	if len(items) == 0 {
		return nil
	}
	for _, task := range items {
		if err := processMention(task); err != nil {
			log.Error("ccopilot: failed to process mention task: %v", err)
		}
	}
	return nil
}
