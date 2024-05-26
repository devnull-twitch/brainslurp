package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	"github.com/devnull-twitch/brainslurp/lib/tag"
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
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))
				flowNo, _ := strconv.Atoi(r.PathValue("flowNumber"))

				if r.Method == "GET" {
					renderFlowEditForm(db, uint64(projectNo), uint64(flowNo), w, r)
				}
				if r.Method == "POST" {
					handleEditFlowSubmit(db, uint64(projectNo), uint64(flowNo), w, r)
				}
			},
		)
	}
}

func renderFlowEditForm(db *badger.DB, projectNo uint64, flowNo uint64, w http.ResponseWriter, r *http.Request) {
	flow, _ := flows.Get(db, projectNo, flowNo)

	projectTags, err := tag.List(db, projectNo)
	if err != nil {
		logrus.WithError(err).Warn("error loading tags")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Server error").Render(r.Context(), w)
	}

	if r.Header.Get("HX-Request") != "" {
		pages.FlowFormBody(projectNo, flow, projectTags).Render(r.Context(), w)
	} else {
		pages.FlowForm(projectNo, flow, projectTags).Render(r.Context(), w)
	}
}

func handleEditFlowSubmit(db *badger.DB, projectNo uint64, flowNo uint64, w http.ResponseWriter, r *http.Request) {
	flowEditObj, err := parseFormFlow(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("missing flow title").Render(r.Context(), w)
		return
	}

	flow, _ := flows.Get(db, projectNo, flowNo)
	flow.Title = flowEditObj.GetTitle()
	flow.Requirements = flowEditObj.GetRequirements()
	flow.Actions = flowEditObj.GetActions()

	if err := flows.Update(db, projectNo, flow); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("error saving flow").Render(r.Context(), w)
		return
	}
}
