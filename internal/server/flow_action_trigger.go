package server

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/internal/server/components/shared"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleFlowActionTrigger(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkIssueNumber,
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))
				flowNo, _ := strconv.Atoi(r.PathValue("flowNumber"))
				issueNo, _ := strconv.Atoi(r.PathValue("issueNo"))

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

				issueObj, issueFlows, _ := issues.Get(db, uint64(projectNo), uint64(issueNo))
				for _, connectedFlow := range issueFlows {
					if connectedFlow.GetNumber() == uint64(flowNo) {
						flowAction := connectedFlow.GetActions()[actionIndex]

						processFlow(flowAction, issueObj)

						modIssue, issueFlows, err := issues.Update(db, uint64(projectNo), issueObj)
						if err != nil {
							logrus.WithError(err).Warn("error updating issue")
							w.WriteHeader(http.StatusInternalServerError)
							pages.Error("Error updating issue").Render(r.Context(), w)
							return
						}

						tagMap := make(map[uint64]*pb_tag.Tag)
						tagList, err := tag.GetMany(db, uint64(projectNo), issueObj.TagNumbers)
						if err != nil {
							pages.Error("Error loading issues").Render(r.Context(), w)
							return
						}
						for _, tagObj := range tagList {
							tagMap[tagObj.GetNumber()] = tagObj
						}

						w.WriteHeader(http.StatusOK)
						if r.Header.Get("HX-Request") != "" {
							shared.IssueRow(uint64(projectNo), modIssue, issueFlows, tagMap).Render(r.Context(), w)
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

	issue.TagNumbers = newTags
}

func checkIssueNumber(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
	projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))

	issueNo, err := strconv.Atoi(r.PathValue("issueNo"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "issueNo",
		}).Warn("URL parsing error")
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

	if _, _, err := issues.Get(db, uint64(projectNo), uint64(issueNo)); err != nil {
		logrus.WithError(err).Warn("error loading issue")
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

	next()
}
