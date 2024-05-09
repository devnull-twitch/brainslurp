package server

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/internal/server/components/shared"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	"github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

// HandleNewIssueTag handles requests for POST /project/{projectNo}/issue/{issueNo}/tags
func HandleNewIssueTag(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))

				issueNo, err := strconv.Atoi(r.PathValue("issueNo"))
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error":      err,
						"path_value": "issueNo",
					}).Warn("URL parsing error")
					w.WriteHeader(http.StatusBadRequest)
					pages.Error("Error parsing URL").Render(r.Context(), w)
					return
				}

				issueObj, _, err := issues.Get(db, uint64(projectNo), uint64(issueNo))
				if err != nil {
					logrus.WithError(err).Warn("error getting issue")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error updating issue").Render(r.Context(), w)
					return
				}

				newTagTitle := r.FormValue("title")
				newTagColor := r.FormValue("color")

				issueObj.Tags = append(issueObj.Tags, &issue.Tag{Label: newTagTitle, ColorCode: newTagColor})
				modIssue, issueFlows, err := issues.Update(db, uint64(projectNo), issueObj)
				if err != nil {
					logrus.WithError(err).Warn("error updating issue")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error updating issue").Render(r.Context(), w)
					return
				}

				if r.Header.Get("HX-Request") != "" {
					shared.IssueRow(uint64(projectNo), modIssue, issueFlows).Render(r.Context(), w)
				}
			},
		)
	}
}

// HandleIssueTag handles requests for DELETE /project/{projectNo}/issue/{issueNo}/tag/{tagName}
func HandleIssueTag(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))

				issueNo, err := strconv.Atoi(r.PathValue("issueNo"))
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error":      err,
						"path_value": "issueNo",
					}).Warn("URL parsing error")
					w.WriteHeader(http.StatusBadRequest)
					pages.Error("Error parsing URL").Render(r.Context(), w)
					return
				}

				issueObj, _, err := issues.Get(db, uint64(projectNo), uint64(issueNo))
				if err != nil {
					logrus.WithError(err).Warn("error getting issue")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error updating issue").Render(r.Context(), w)
					return
				}

				removeIndex := -1
				for tagIndex, tag := range issueObj.Tags {
					if tag.GetLabel() == r.PathValue("tagName") {
						removeIndex = tagIndex
					}
				}
				if removeIndex > -1 {
					issueObj.Tags = slices.Delete(issueObj.Tags, removeIndex, removeIndex+1)
				}

				modIssue, issueFlows, err := issues.Update(db, uint64(projectNo), issueObj)
				if err != nil {
					logrus.WithError(err).Warn("error updating issue")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error updating issue").Render(r.Context(), w)
					return
				}

				w.WriteHeader(http.StatusCreated)
				if r.Header.Get("HX-Request") != "" {
					shared.IssueRow(uint64(projectNo), modIssue, issueFlows).Render(r.Context(), w)
				}
			},
		)
	}
}
