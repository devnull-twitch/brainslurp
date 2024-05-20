package server

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/internal/server/components/shared"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	"github.com/devnull-twitch/brainslurp/lib/tag"
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

				newTagIdStr := r.FormValue("tag_id")
				newTagId, err := strconv.Atoi(newTagIdStr)
				if err != nil {
					logrus.WithError(err).Warn("error getting issue")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error updating issue").Render(r.Context(), w)
					return
				}

				issueObj.TagNumbers = append(issueObj.GetTagNumbers(), uint64(newTagId))
				modIssue, issueFlows, err := issues.Update(db, uint64(projectNo), issueObj)
				if err != nil {
					logrus.WithError(err).Warn("error updating issue")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error updating issue").Render(r.Context(), w)
					return
				}

				tagMap := make(map[uint64]*pb_tag.Tag)
				tagList, err := tag.List(db, uint64(projectNo))
				if err != nil {
					pages.Error("Error loading tags").Render(r.Context(), w)
					return
				}
				for _, tagObj := range tagList {
					tagMap[tagObj.GetNumber()] = tagObj
				}

				if r.Header.Get("HX-Request") != "" {
					shared.IssueRow(uint64(projectNo), modIssue, issueFlows, tagMap).Render(r.Context(), w)
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

				w.WriteHeader(http.StatusCreated)
				if r.Header.Get("HX-Request") != "" {
					shared.IssueRow(uint64(projectNo), modIssue, issueFlows, tagMap).Render(r.Context(), w)
				}
			},
		)
	}
}
