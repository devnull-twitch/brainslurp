package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/project"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleProjectSettings(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectObj := getProjectFromContext(ctx)
				userObj := getUserFromContext(ctx)

				if r.Method == "GET" {
					renderProjectSettingsForm(db, projectObj, w, r)
				}
				if r.Method == "POST" {
					handleProjectSettingsSubmit(db, projectObj, userObj, w, r)
				}
			},
		)
	}
}

func renderProjectSettingsForm(db *badger.DB, projectObj *pb_project.Project, w http.ResponseWriter, r *http.Request) {
	userNos := make([]uint64, len(projectObj.GetMembers()))
	for i, membership := range projectObj.GetMembers() {
		userNos[i] = membership.GetUserNo()
	}
	projectUsers, err := user.ListAll(db)
	if err != nil {
		logrus.WithError(err).Error("error loading users")
		pages.Error("Internal server error").Render(r.Context(), w)
		return
	}

	if r.Header.Get("HX-Request") != "" {
		pages.ProjectSettingsBody(projectObj, projectUsers).Render(r.Context(), w)
	} else {
		pages.ProjectSettings(projectObj, projectUsers).Render(r.Context(), w)
	}
}

func handleProjectSettingsSubmit(db *badger.DB, projectObj *pb_project.Project, userObj *pb_user.User, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Form.Get("public") != "" {
		projectObj.Public = true
	} else {
		projectObj.Public = false
	}

	projectObj.Members = make([]*pb_project.Member, 0)
	for _, userNoStr := range r.Form["project_member"] {
		//TODO: validate if its a valid user ID
		userNo, err := strconv.Atoi(userNoStr)
		if err != nil {
			logrus.WithError(err).Warn("unable to convert user number to int")
			return
		} else {
			projectObj.Members = append(projectObj.Members, &pb_project.Member{
				UserNo: uint64(userNo),
			})
		}
	}

	if err := project.Update(db, projectObj); err != nil {
		logrus.WithError(err).Error("error updating project")
		pages.Error("Internal server error").Render(r.Context(), w)
		return
	}

	projects := make([]*pb_project.Project, len(userObj.GetMemberships()))
	var err error
	for i, userMembership := range userObj.GetMemberships() {
		projects[i], err = project.Get(db, userMembership.GetProjectNo())
		if err != nil {
			logrus.WithError(err).Error("unable to load user project")
			continue
		}
	}

	w.Header().Set("HX-Push-Url", "/projects")

	if r.Header.Get("HX-Request") != "" {
		pages.ProjectListingBody(projects).Render(r.Context(), w)
	} else {
		pages.ProjectListing(projects).Render(r.Context(), w)
	}
}
