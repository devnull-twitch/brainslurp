package server

import (
	"context"
	"net/http"
	"slices"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/internal/server/components/shared"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleFlowActionTrigger(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkIssueNumber,
			checkFlowNumber,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				db := getDbFromContext(ctx)
				projectObj := getProjectFromContext(ctx)
				issueObj := getIssueFromContext(ctx)
				issueFlows := getIssueFlowsFromContext(ctx)
				flowObj := getFlowFromContext(ctx)

				actionIndex, err := strconv.Atoi(r.PathValue("actionIndex"))
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error":      err,
						"path_value": "actionIndex",
					}).Warn("URL parsing error")
					w.WriteHeader(http.StatusBadRequest)
					pages.Error("Error parsing URL").Render(r.Context(), w)
					return
				}

				for _, connectedFlow := range issueFlows {
					if connectedFlow.GetNumber() == flowObj.GetNumber() {
						flowAction := connectedFlow.GetActions()[actionIndex]

						processFlow(flowAction, issueObj)

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

						w.WriteHeader(http.StatusOK)
						if r.Header.Get("HX-Request") != "" {
							shared.IssueRow(projectObj.GetNumber(), modIssue, issueFlows, tagMap, userMap).Render(r.Context(), w)
						}
						return
					}
				}

				w.WriteHeader(http.StatusNotFound)
			},
		)
	}
}

func processFlow(flowAction *pb_flow.FlowActions, issue *pb_issue.Issue) {
	newTags := make([]uint64, len(issue.GetTagNumbers()))
	copy(newTags, issue.GetTagNumbers())
	for _, addTag := range flowAction.GetAddTagIds() {
		if !slices.Contains(newTags, addTag) {
			newTags = append(newTags, addTag)
		}
	}

	for _, removeTag := range flowAction.GetRemoveTagIds() {
		removeIndex := slices.Index(newTags, removeTag)
		if removeIndex > -1 {
			newTags = slices.Delete(newTags, removeIndex, removeIndex+1)
		}
	}

	if flowAction.GetRemoveAllAssignees() {
		issue.AssignedTo = []uint64{}
	}

	if len(flowAction.GetAssignUser()) > 0 {
		issue.AssignedTo = append(issue.AssignedTo, flowAction.GetAssignUser()...)
	}

	issue.TagNumbers = newTags
}

func checkIssueNumber(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
	db := getDbFromContext(ctx)
	projectObj := getProjectFromContext(ctx)

	issueNo, err := strconv.Atoi(r.PathValue("issueNo"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "issueNo",
		}).Warn("URL parsing error")
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

	issueObj, issueFlows, err := issues.Get(db, projectObj.GetNumber(), uint64(issueNo))
	if err != nil {
		logrus.WithError(err).Warn("error loading issue")
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

	next(setIssueFlowsOnContext(setIssueOnContext(ctx, issueObj), issueFlows))
}
