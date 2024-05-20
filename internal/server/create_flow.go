package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	"github.com/devnull-twitch/brainslurp/lib/proto/issue"
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
		pages.FlowFormBody(projectNo, nil, tagList).Render(r.Context(), w)
	} else {
		pages.FlowForm(projectNo, nil, tagList).Render(r.Context(), w)
	}
}

func handleNewFlowSubmit(
	db *badger.DB,
	projectNo uint64,
	w http.ResponseWriter, r *http.Request,
) {
	newFlow, err := parseFormFlow(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("missing flow title").Render(r.Context(), w)
		return
	}

	_, err = flows.Create(db, flows.CreateOptions{
		ProjectNo:    projectNo,
		Title:        newFlow.GetTitle(),
		Requirements: newFlow.GetRequirements(),
		Actions:      newFlow.GetActions(),
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "projectNo",
		}).Warn("error inserting new flow")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Error inserting flow").Render(r.Context(), w)
	}

	flows, err := flows.List(db, uint64(projectNo))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "projectNo",
		}).Warn("URL parsing error")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Error loading flows").Render(r.Context(), w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	pages.FlowListBody(projectNo, flows).Render(r.Context(), w)
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

func StrToIntSlice[O uint | uint32 | uint64 | int | int32 | int64](strIn []string) []O {
	out := make([]O, 0, len(strIn))
	for _, strID := range strIn {
		id, err := strconv.Atoi(strID)
		if err == nil {
			out = append(out, O(id))
		}
	}

	return out
}

func parseFormFlow(r *http.Request) (*pb_flow.Flow, error) {
	r.ParseForm()
	flowTitle := r.Form.Get("title")
	if flowTitle == "" {
		return nil, fmt.Errorf("missing title")
	}

	reqs := make([]*pb_flow.FlowRequirement, 0)
	actions := make([]*pb_flow.FlowActions, 0)

	hasReqs := true
	reqIndex := 0
	for hasReqs {
		reqTagsRequired, hasRequired := r.Form[fmt.Sprintf("req_%d_required", reqIndex)]
		reqTagsProhibited, hasProhibited := r.Form[fmt.Sprintf("req_%d_prohibited", reqIndex)]

		if !hasRequired && !hasProhibited {
			hasReqs = false
			continue
		}

		reqCategory := issue.IssueCategory_Unknown
		reqCategoryStr := r.Form.Get(fmt.Sprintf("req_%d_category", reqIndex))
		if reqCategoryStr != "" {
			reqCategoryInt, err := strconv.Atoi(reqCategoryStr)
			if err != nil {
				logrus.WithError(err).Warn("int conversion error for category")
			} else {
				reqCategory = issue.IssueCategory(reqCategoryInt)
			}
		}

		reqs = append(reqs, &pb_flow.FlowRequirement{
			InCategory:       reqCategory,
			RequiredTagIds:   StrToIntSlice[uint64](reqTagsRequired),
			ProhibitedTagIds: StrToIntSlice[uint64](reqTagsProhibited),
		})

		reqIndex += 1
	}

	hasActs := true
	actIndex := 0
	for hasActs {
		actTagsAdd, hasAddTags := r.Form[fmt.Sprintf("act_%d_add", actIndex)]
		actTagsRemove, hasRemoveTags := r.Form[fmt.Sprintf("act_%d_remove", actIndex)]

		if !hasAddTags && !hasRemoveTags {
			hasActs = false
			continue
		}

		actions = append(actions, &pb_flow.FlowActions{
			Title:        r.Form.Get(fmt.Sprintf("act_%d_title", actIndex)),
			AddTagIds:    StrToIntSlice[uint64](actTagsAdd),
			RemoveTagIds: StrToIntSlice[uint64](actTagsRemove),
		})

		actIndex += 1
	}

	return &pb_flow.Flow{
		Title:        flowTitle,
		Requirements: reqs,
		Actions:      actions,
	}, nil
}
