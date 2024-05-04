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

// HandleNewIssueTag handles requests for /project/{projectNo}/issue/{issueNo}/tags
func HandleNewIssueTag(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			doAddTag(db, w, r)
		} else {
			// TODO render form page for user input
			// also ... Watch out for htmx requests versus normal ones.
		}
	}
}

func doAddTag(db *badger.DB, w http.ResponseWriter, r *http.Request) {
	projectNo, err := strconv.Atoi(r.PathValue("projectNo"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "projectNo",
		}).Warn("URL parsing error")
		w.WriteHeader(http.StatusBadRequest)
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

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
}

// HandleIssueTag handles requests for DELETE /project/{projectNo}/issue/{issueNo}/tag/{tagName}
func HandleIssueTag(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		projectNo, err := strconv.Atoi(r.PathValue("projectNo"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":      err,
				"path_value": "projectNo",
			}).Warn("URL parsing error")
			w.WriteHeader(http.StatusBadRequest)
			pages.Error("Error parsing URL").Render(r.Context(), w)
			return
		}

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

		if r.Header.Get("HX-Request") != "" {
			shared.IssueRow(uint64(projectNo), modIssue, issueFlows).Render(r.Context(), w)
		}
	}
}
