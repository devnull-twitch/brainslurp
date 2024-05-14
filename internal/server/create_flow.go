package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleFlowCreate(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				// projectNo is already validated by authUserWithProjectNo so we dont need to error check again
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))

				if r.Method == "GET" {
					renderFlowCreateForm(db, uint64(projectNo), w, r)
				}
				if r.Method == "POST" {
					handleNewFlowSubmit(db, uint64(projectNo), w, r)
				}
			},
		)
	}
}

func renderFlowCreateForm(db *badger.DB, projectNo uint64, w http.ResponseWriter, r *http.Request) {
	tagList, err := tag.List(db, projectNo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "projectNo",
		}).Warn("error inserting new flow")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Error inserting flow").Render(r.Context(), w)
	}

	if r.Header.Get("HX-Request") != "" {
		pages.FlowCreateBody(projectNo, tagList).Render(r.Context(), w)
	} else {
		pages.FlowCreate(projectNo, tagList).Render(r.Context(), w)
	}
}

func handleNewFlowSubmit(
	db *badger.DB,
	projectNo uint64,
	w http.ResponseWriter, r *http.Request,
) {
	r.ParseForm()
	flowTitle := r.Form.Get("title")
	if flowTitle == "" {
		handleInputErr("title", fmt.Errorf("no flow title"), w)
		return
	}

	reqs := make([]*pb_flow.FlowRequirement, 0)
	actions := make([]*pb_flow.FlowActions, 0)

	newFlowNo, err := flows.Create(db, flows.CreateOptions{
		ProjectNo:    projectNo,
		Title:        flowTitle,
		Requirements: reqs,
		Actions:      actions,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "projectNo",
		}).Warn("error inserting new flow")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Error inserting flow").Render(r.Context(), w)
	}

	w.Header().Add("HX-Location", fmt.Sprintf("/project/%d/flow/%d/edit", projectNo, newFlowNo))
	w.WriteHeader(http.StatusCreated)
}

func handleInputErr(
	fieldName string, err error,
	w http.ResponseWriter,
) {
	logrus.WithFields(logrus.Fields{
		"error":      err,
		"form_field": fieldName,
	}).Warn("Form Input Error")
	w.WriteHeader(http.StatusBadRequest)
}
