package server

import (
	"net/http"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/internal/server/components/shared"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/dgraph-io/badger/v4"
)

func renderIssueRow(db *badger.DB, projectObj *pb_project.Project, issueObj *pb_issue.Issue, issueFlows []*pb_flow.Flow, r *http.Request, w http.ResponseWriter) {
	tagMap := make(map[uint64]*pb_tag.Tag)
	tagList, err := tag.List(db, projectObj.GetNumber())
	if err != nil {
		pages.Error("Error loading tags").Render(r.Context(), w)
		return
	}
	for _, tagObj := range tagList {
		tagMap[tagObj.GetNumber()] = tagObj
	}

	userNos := make([]uint64, len(projectObj.GetMembers()))
	for i, membership := range projectObj.GetMembers() {
		userNos[i] = membership.GetUserNo()
	}
	projectUsers, err := user.List(db, userNos)
	if err != nil {
		pages.Error("Error loading users").Render(r.Context(), w)
		return
	}
	userMap := make(map[uint64]*pb_user.User)
	for _, memberUserObj := range projectUsers {
		userMap[memberUserObj.GetNumber()] = memberUserObj
	}

	shared.IssueRow(projectObj.GetNumber(), issueObj, issueFlows, tagMap, userMap).Render(r.Context(), w)
}
