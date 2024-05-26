package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/database"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/devnull-twitch/brainslurp/lib/view"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func userIsInProject(userObj *pb_user.User, projectNo uint64) bool {
	for _, userMembership := range userObj.GetMemberships() {
		if userMembership.GetProjectNo() == projectNo {
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
		walkChain(
			db, w, r,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				db := getDbFromContext(ctx)

				projectNo, err := strconv.Atoi(r.PathValue("projectNo"))
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error":      err,
						"path_value": "projectNo",
					}).Warn("URL parsing error")
					pages.Error("Error parsing URL").Render(r.Context(), w)
					return
				}

				projectObj := &pb_project.Project{}
				if err := db.View(func(txn *badger.Txn) error {
					projectKey := database.Keygen(database.ProjectPrefix, uint64(projectNo))
					projectItem, err := txn.Get(projectKey)
					if err != nil {
						return fmt.Errorf("error loading project %d: %w", projectNo, err)
					}

					if err := database.Read(projectItem, projectObj); err != nil {
						return fmt.Errorf("error reading project: %w", err)
					}

					return nil
				}); err != nil {
					logrus.WithError(err).Error("error authenticating user")
					pages.Error("Halp! Something went wrong").Render(r.Context(), w)
					return
				}

				if !projectObj.GetPublic() {
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

					if !userIsInProject(userObj, projectObj.GetNumber()) {
						pages.Error("You may not see this!").Render(r.Context(), w)
						return
					}
				}

				userNos := make([]uint64, len(projectObj.GetMembers()))
				for i, membership := range projectObj.GetMembers() {
					userNos[i] = membership.GetUserNo()
				}
				projectUsers, err := user.List(db, userNos)
				if err != nil {
					logrus.WithError(err).Error("error loading users")
					pages.Error("Internal server error").Render(r.Context(), w)
					return
				}
				userMap := make(map[uint64]*pb_user.User)
				for _, memberUserObj := range projectUsers {
					userMap[memberUserObj.GetNumber()] = memberUserObj
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
					pages.Error("Error loading tags").Render(r.Context(), w)
					return
				}
				for _, tagObj := range tagList {
					tagMap[tagObj.GetNumber()] = tagObj
				}

				if r.Header.Get("HX-Request") != "" || w.Header().Get("HX-Retarget") != "" {
					pages.BodyLogoOOB(uint64(projectNo)).Render(r.Context(), w)
					if r.Header.Get("HX-Target") == "issue-list" {
						pages.ListItems(uint64(projectNo), issueList, mapFlowsToIssues, tagMap, userMap).Render(r.Context(), w)
					} else {
						pages.ListBody(uint64(projectNo), issueList, mapFlowsToIssues, viewList, tagMap, userMap).Render(r.Context(), w)
					}
				} else {
					pages.List(uint64(projectNo), issueList, mapFlowsToIssues, viewList, tagMap, userMap).Render(r.Context(), w)
				}
			},
		)
	}
}
