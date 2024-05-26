package server

import (
	"context"
	"net/http"
	"slices"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/internal/server/components/shared"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

// HandleNewIssueTag handles requests for POST /project/{projectNo}/issue/{issueNo}/tags
func HandleNewIssueTag(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkIssueNumber,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectObj := getProjectFromContext(ctx)
				issueObj := getIssueFromContext(ctx)

				newTagIdStr := r.FormValue("tag_id")
				newTagId, err := strconv.Atoi(newTagIdStr)
				if err != nil {
					logrus.WithError(err).Warn("error getting issue")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error updating issue").Render(r.Context(), w)
					return
				}

				issueObj.TagNumbers = append(issueObj.GetTagNumbers(), uint64(newTagId))
				modIssue, issueFlows, err := issues.Update(db, projectObj.GetNumber(), issueObj)
				if err != nil {
					logrus.WithError(err).Warn("error updating issue")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error updating issue").Render(r.Context(), w)
					return
				}

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

				if r.Header.Get("HX-Request") != "" {
					shared.IssueRow(projectObj.GetNumber(), modIssue, issueFlows, tagMap, userMap).Render(r.Context(), w)
				}
			},
		)
	}
}

// HandleIssueTag handles requests for DELETE /project/{projectNo}/issue/{issueNo}/tag/{tagNumber}
func HandleIssueTag(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkIssueNumber,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectObj := getProjectFromContext(ctx)
				issueObj := getIssueFromContext(ctx)

				removeTagNoStr := r.PathValue("tagNumber")
				removeTagNumber, err := strconv.Atoi(removeTagNoStr)
				if err != nil {
					logrus.WithError(err).Warn("error converting tag number from string")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error parsing input").Render(r.Context(), w)
					return
				}

				removeIndex := -1
				for tagIndex, tagNo := range issueObj.GetTagNumbers() {
					if tagNo == uint64(removeTagNumber) {
						removeIndex = tagIndex
					}
				}
				if removeIndex > -1 {
					issueObj.TagNumbers = slices.Delete(issueObj.TagNumbers, removeIndex, removeIndex+1)
				}

				modIssue, issueFlows, err := issues.Update(db, projectObj.GetNumber(), issueObj)
				if err != nil {
					logrus.WithError(err).Warn("error updating issue")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error updating issue").Render(r.Context(), w)
					return
				}

				tagMap := make(map[uint64]*pb_tag.Tag)
				tagList, err := tag.GetMany(db, projectObj.GetNumber(), issueObj.TagNumbers)
				if err != nil {
					pages.Error("Error loading issues").Render(r.Context(), w)
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

				w.WriteHeader(http.StatusCreated)
				if r.Header.Get("HX-Request") != "" {
					shared.IssueRow(projectObj.GetNumber(), modIssue, issueFlows, tagMap, userMap).Render(r.Context(), w)
				}
			},
		)
	}
}
