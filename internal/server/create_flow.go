package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	"github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleFlowCreate(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				if r.Method == "GET" {
					renderFlowCreateForm(ctx, w, r)
				}
				if r.Method == "POST" {
					handleNewFlowSubmit(ctx, w, r)
				}
			},
		)
	}
}

func renderFlowCreateForm(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db := getDbFromContext(ctx)
	projectObj := getProjectFromContext(ctx)
	tagList, err := tag.List(db, projectObj.GetNumber())
	if err != nil {
		logrus.WithError(err).Warn("error loading tags")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Internal System Error").Render(r.Context(), w)
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
		pages.FlowFormBody(projectObj.GetNumber(), nil, tagList, projectUsers).Render(r.Context(), w)
	} else {
		pages.FlowForm(projectObj.GetNumber(), nil, tagList, projectUsers).Render(r.Context(), w)
	}
}

func handleNewFlowSubmit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db := getDbFromContext(ctx)
	projectObj := getProjectFromContext(ctx)

	newFlow, err := parseFormFlow(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("missing flow title").Render(r.Context(), w)
		return
	}

	_, err = flows.Create(db, flows.CreateOptions{
		ProjectNo:    projectObj.GetNumber(),
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

	flows, err := flows.List(db, projectObj.GetNumber())
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
	pages.FlowListBody(projectObj.GetNumber(), flows).Render(r.Context(), w)
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

		removeAllAssignees := false
		if r.Form.Get(fmt.Sprintf("act_%d_remove_assignees", actIndex)) != "" {
			removeAllAssignees = true
		}

		assignUserInput := r.Form.Get(fmt.Sprintf("act_%d_assign_user", actIndex))
		assignUserNumbers := make([]uint64, 0)
		if assignUserInput != "" && assignUserInput != "0" {
			userNo, err := strconv.Atoi(assignUserInput)
			if err != nil {
				logrus.WithError(err).Warn("unable to convert assignee user number to int")
			} else {
				assignUserNumbers = append(assignUserNumbers, uint64(userNo))
			}
		}

		if !hasAddTags && !hasRemoveTags && len(assignUserNumbers) == 0 && !removeAllAssignees {
			hasActs = false
			continue
		}

		actions = append(actions, &pb_flow.FlowActions{
			Title:              r.Form.Get(fmt.Sprintf("act_%d_title", actIndex)),
			AddTagIds:          StrToIntSlice[uint64](actTagsAdd),
			RemoveTagIds:       StrToIntSlice[uint64](actTagsRemove),
			RemoveAllAssignees: removeAllAssignees,
			AssignUser:         assignUserNumbers,
		})

		actIndex += 1
	}

	return &pb_flow.Flow{
		Title:        flowTitle,
		Requirements: reqs,
		Actions:      actions,
	}, nil
}
