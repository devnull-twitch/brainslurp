package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	"github.com/devnull-twitch/brainslurp/lib/project"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/devnull-twitch/brainslurp/lib/view"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func userIsInProject(userObj user.User, projectNo uint64) bool {
	for _, userProjects := range userObj.Projects {
		if userProjects == projectNo {
			return true
		}
	}

	return false
}

// HandleIssueList handles requests for:
// * /project/{projectNo}/issues/view/{viewNo}
// * /project/{projectNo}/issues
func HandleIssueList(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		projectNo, err := strconv.Atoi(r.PathValue("projectNo"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":      err,
				"path_value": "projectNo",
			}).Warn("URL parsing error")
			pages.Error("Error parsing URL").Render(r.Context(), w)
			return
		}

		projectObj, err := project.Get(db, uint64(projectNo))
		if err != nil {
			logrus.WithError(err).Error("project could not be loaded")
			pages.Error("Project missing?!").Render(r.Context(), w)
			return
		}

		if !projectObj.Public {
			userObj, err := parseAuthAndLoadUser(db, r)
			if err != nil {
				if errors.Is(err, ErrNoCookie) {
					// maybe if user is not logged in redirect to start page?
					w.Header().Set("Location", "/")
					w.WriteHeader(http.StatusTemporaryRedirect)
					return
				}

				logrus.WithError(err).Error("error authenticating user")
				pages.Error("Halp! Something went wrong").Render(r.Context(), w)
				return
			}

			if !userIsInProject(userObj, uint64(projectNo)) {
				pages.Error("You may not see this!").Render(r.Context(), w)
				return
			}
		}

		var issueList []*pb_issue.Issue
		var mapFlowsToIssues map[uint64][]*pb_flow.Flow
		if r.PathValue("viewNo") == "" {
			var err error
			issueList, mapFlowsToIssues, err = issues.ListAll(db, uint64(projectNo))
			if err != nil {
				pages.Error("Error loading issues").Render(r.Context(), w)
				return
			}
		} else {
			viewNo, err := strconv.Atoi(r.PathValue("viewNo"))
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":      err,
					"path_value": "viewNo",
				}).Warn("URL parsing error")
				pages.Error("Error parsing URL").Render(r.Context(), w)
				return
			}

			issueList, mapFlowsToIssues, err = issues.ListFromView(db, uint64(projectNo), uint64(viewNo))
			if err != nil {
				logrus.WithError(err).Error("error loading issues")
				pages.Error("Error loading issues").Render(r.Context(), w)
				return
			}
		}

		viewList, err := view.List(db, uint64(projectNo))
		if err != nil {
			pages.Error("Error loading issues").Render(r.Context(), w)
			return
		}

		tagMap := make(map[uint64]*pb_tag.Tag)
		tagList, err := tag.List(db, uint64(projectNo))
		if err != nil {
			pages.Error("Error loading issues").Render(r.Context(), w)
			return
		}
		for _, tagObj := range tagList {
			tagMap[tagObj.GetNumber()] = tagObj
		}

		if r.Header.Get("HX-Request") != "" || w.Header().Get("HX-Retarget") != "" {
			if r.Header.Get("HX-Target") == "issue-list" {
				pages.ListItems(uint64(projectNo), issueList, mapFlowsToIssues, tagMap).Render(r.Context(), w)
			} else {
				pages.ListBody(uint64(projectNo), issueList, mapFlowsToIssues, viewList, tagMap).Render(r.Context(), w)
			}
		} else {
			pages.List(uint64(projectNo), issueList, mapFlowsToIssues, viewList, tagMap).Render(r.Context(), w)
		}
	}
}
