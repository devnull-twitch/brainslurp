package server

import (
	"context"
	"net/http"
	"slices"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleIssueUnassign(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkIssueNumber,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectObj := getProjectFromContext(ctx)
				issueObj := getIssueFromContext(ctx)

				userNo, err := strconv.Atoi(r.PathValue("userNo"))
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error":      err,
						"path_value": "userNo",
					}).Warn("URL parsing error")
					pages.Error("Error parsing URL").Render(r.Context(), w)
					return
				}

				if !isUserAssigned(issueObj, uint64(userNo)) {
					logrus.WithFields(logrus.Fields{
						"issue_no": issueObj.GetNumber(),
						"user_no":  userNo,
					}).Warn("user isnt assigned to issue")
					pages.Error("User error").Render(r.Context(), w)
					return
				}

				removeIndex := slices.Index(issueObj.GetAssignedTo(), uint64(userNo))
				issueObj.AssignedTo = slices.Delete(issueObj.GetAssignedTo(), removeIndex, removeIndex+1)

				_, _, err = issues.Update(db, projectObj.GetNumber(), issueObj)
				if err != nil {
					logrus.WithError(err).Warn("error saving issue with new comment")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error saving issue with comment").Render(r.Context(), w)
				}

				w.WriteHeader(http.StatusOK)
			},
		)
	}
}

func HandleIssueAssign(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkIssueNumber,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectObj := getProjectFromContext(ctx)
				issueObj := getIssueFromContext(ctx)

				r.ParseForm()
				userNo, err := strconv.Atoi(r.Form.Get("assign_user"))
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error":      err,
						"path_value": "userNo",
					}).Warn("URL parsing error")
					pages.Error("Error parsing URL").Render(r.Context(), w)
					return
				}

				if isUserAssigned(issueObj, uint64(userNo)) {
					logrus.WithFields(logrus.Fields{
						"issue_no": issueObj.GetNumber(),
						"user_no":  userNo,
					}).Warn("user is already assigned to issue")
					pages.Error("User error").Render(r.Context(), w)
					return
				}

				issueObj.AssignedTo = append(issueObj.AssignedTo, uint64(userNo))

				modIssue, issueFlows, err := issues.Update(db, projectObj.GetNumber(), issueObj)
				if err != nil {
					logrus.WithError(err).Warn("error saving issue with new comment")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error saving issue with comment").Render(r.Context(), w)
				}

				renderIssueRow(db, projectObj, modIssue, issueFlows, r, w)
			},
		)
	}
}

func isUserAssigned(issueObj *pb_issue.Issue, cheeckUserNo uint64) bool {
	for _, issueAssigneeNo := range issueObj.GetAssignedTo() {
		if issueAssigneeNo == cheeckUserNo {
			return true
		}
	}

	return false
}
