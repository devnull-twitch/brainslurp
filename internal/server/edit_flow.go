package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func checkFlowNumber(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
	db := getDbFromContext(ctx)
	projectObj := getProjectFromContext(ctx)

	flowNo, err := strconv.Atoi(r.PathValue("flowNumber"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "flowNumber",
		}).Warn("URL parsing error")
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

	flowObj, err := flows.Get(db, projectObj.GetNumber(), uint64(flowNo))
	if err != nil {
		logrus.WithError(err).Warn("error loading flow for edit")
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

	next(setFlowOnContext(ctx, flowObj))
}

func HandleFlowEdit(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkFlowNumber,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectObj := getProjectFromContext(ctx)
				flowObj := getFlowFromContext(ctx)

				if r.Method == "GET" {
					renderFlowEditForm(db, projectObj, flowObj, w, r)
				}
				if r.Method == "POST" {
					handleEditFlowSubmit(db, projectObj, flowObj, w, r)
				}
			},
		)
	}
}

func renderFlowEditForm(db *badger.DB, projectObj *pb_project.Project, flowObj *pb_flow.Flow, w http.ResponseWriter, r *http.Request) {
	projectTags, err := tag.List(db, projectObj.GetNumber())
	if err != nil {
		logrus.WithError(err).Warn("error loading tags")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Server error").Render(r.Context(), w)
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

	if r.Header.Get("HX-Request") != "" {
		pages.FlowFormBody(projectObj.GetNumber(), flowObj, projectTags, projectUsers).Render(r.Context(), w)
	} else {
		pages.FlowForm(projectObj.GetNumber(), flowObj, projectTags, projectUsers).Render(r.Context(), w)
	}
}

func handleEditFlowSubmit(db *badger.DB, projectObj *pb_project.Project, flowObj *pb_flow.Flow, w http.ResponseWriter, r *http.Request) {
	flowEditObj, err := parseFormFlow(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("missing flow title").Render(r.Context(), w)
		return
	}

	flowObj.Title = flowEditObj.GetTitle()
	flowObj.Requirements = flowEditObj.GetRequirements()
	flowObj.Actions = flowEditObj.GetActions()

	if err := flows.Update(db, projectObj.GetNumber(), flowObj); err != nil {
		logrus.WithError(err).Error("unable to update flow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	flows, err := flows.List(db, projectObj.GetNumber())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "projectNo",
		}).Warn("URL parsing error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Push-Url", fmt.Sprintf("/project/%d/flows", projectObj.GetNumber()))

	if r.Header.Get("HX-Request") != "" {
		pages.FlowListBody(projectObj.GetNumber(), flows).Render(r.Context(), w)
	} else {
		pages.FlowList(projectObj.GetNumber(), flows).Render(r.Context(), w)
	}
}
