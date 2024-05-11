package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
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
					renderFlowCreateForm(uint64(projectNo), w, r)
				}
				if r.Method == "POST" {
					handleNewFlowSubmit(db, uint64(projectNo), w, r)
				}
			},
		)
	}
}

func renderFlowCreateForm(projectNo uint64, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "" {
		pages.FlowCreateBody(projectNo).Render(r.Context(), w)
	} else {
		pages.FlowCreate(projectNo).Render(r.Context(), w)
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
	for reqIndex, tagReq := range r.Form["req_tags"] {
		rawTags := strings.Split(tagReq, ",")
		tags := make([]uint64, 0, len(rawTags))
		for _, tagStr := range rawTags {
			tagNo, err := strconv.Atoi(strings.TrimSpace(tagStr))
			if err != nil {
				logrus.WithError(err).Warn("error converting number")
				continue
			}
			tags = append(tags, uint64(tagNo))
		}

		rawNoTags := strings.Split(r.Form["req_no_tags"][reqIndex], ",")
		noTags := make([]uint64, 0, len(rawNoTags))
		for _, noTagStr := range rawNoTags {
			tagNo, err := strconv.Atoi(strings.TrimSpace(noTagStr))
			if err != nil {
				logrus.WithError(err).Warn("error converting number")
				continue
			}
			noTags = append(noTags, uint64(tagNo))
		}

		catID, err := strconv.Atoi(r.Form["req_category"][reqIndex])
		if err != nil {
			logrus.WithError(err).Warn("unable to parse category ID in create flow form")
		}

		reqs = append(reqs, &pb_flow.FlowRequirement{
			CheckTagIds:   tags,
			CheckNoTagIds: noTags,
			InCategory:    pb_issue.IssueCategory(int32(catID)),
		})
	}

	actions := make([]*pb_flow.FlowActions, 0)
	for actionIndex, actionName := range r.Form["action_name"] {
		rawTagAdds := strings.Split(r.Form["action_adds"][actionIndex], ",")
		addTags := make([]uint64, 0, len(rawTagAdds))
		for _, addTagStr := range rawTagAdds {
			tagNo, err := strconv.Atoi(strings.TrimSpace(addTagStr))
			if err != nil {
				logrus.WithError(err).Warn("error converting number")
				continue
			}
			addTags = append(addTags, uint64(tagNo))
		}

		rawTagRemoval := strings.Split(r.Form["action_removes"][actionIndex], ",")
		removeTags := make([]uint64, 0, len(rawTagRemoval))
		for _, removeTagStr := range rawTagRemoval {
			tagNo, err := strconv.Atoi(strings.TrimSpace(removeTagStr))
			if err != nil {
				logrus.WithError(err).Warn("error converting number")
				continue
			}
			removeTags = append(removeTags, uint64(tagNo))
		}

		actions = append(actions, &pb_flow.FlowActions{
			Title:        actionName,
			AddTagIds:    addTags,
			RemoveTagIds: removeTags,
		})
	}

	if err := flows.Create(db, flows.CreateOptions{
		ProjectNo:    projectNo,
		Title:        flowTitle,
		Requirements: reqs,
		Actions:      actions,
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "projectNo",
		}).Warn("error inserting new flow")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Error inserting flow").Render(r.Context(), w)
	}

	HandleFlowList(db)(w, r)
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
