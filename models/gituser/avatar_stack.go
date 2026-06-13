// Copyright 2026 The CreeperCoding Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package gituser

import (
	"context"

	"creepercoding.dev/models/user"
	"creepercoding.dev/modules/git"
	"creepercoding.dev/modules/log"
)

// AvatarStackData is the view-model for the AvatarStack render helpers. Participants[0] is
// the primary participant (commit author), painted on top; the rest follow.
type AvatarStackData struct {
	Participants      []*CommitParticipant
	SearchByEmailLink string
}

func BuildAvatarStackData(ctx context.Context, allParticipants []*git.CommitIdentity, emailUserMap *user.EmailUserMap) *AvatarStackData {
	if emailUserMap == nil {
		emails := make([]string, len(allParticipants))
		for i, sig := range allParticipants {
			emails[i] = sig.Email
		}
		var err error
		emailUserMap, err = user.GetUsersByEmails(ctx, emails)
		if err != nil {
			log.Error("GetUsersByEmails failed: %v", err)
		}
	}
	ret := &AvatarStackData{
		Participants: make([]*CommitParticipant, 0, len(allParticipants)),
	}
	for _, p := range allParticipants {
		var giteaUser *user.User
		if emailUserMap != nil {
			giteaUser = emailUserMap.GetByEmail(p.Email)
		}
		ret.Participants = append(ret.Participants, &CommitParticipant{GiteaUser: giteaUser, GitIdentity: p})
	}
	return ret
}
