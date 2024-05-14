package server

import (
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	"github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func checkFlowNumber(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
	projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))

	flowNo, err := strconv.Atoi(r.PathValue("flowNumber"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "flowNumber",
		}).Warn("URL parsing error")
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

	if _, err := flows.Get(db, uint64(projectNo), uint64(flowNo)); err != nil {
		logrus.WithError(err).Warn("error loading flow for edit")
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

	next()
}

func HandleFlowEdit(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkFlowNumber,
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))
				flowNo, _ := strconv.Atoi(r.PathValue("flowNumber"))

				if r.Method == "GET" {
					renderFlowEditForm(db, uint64(projectNo), uint64(flowNo), w, r)
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

	tagMap := make(map[uint64]*pb_tag.Tag)
	for _, tagObj := range projectTags {
		tagMap[tagObj.GetNumber()] = tagObj
	}

	if r.Header.Get("HX-Request") != "" {
		pages.FlowEditBody(projectNo, flow, tagMap).Render(r.Context(), w)
	} else {
		pages.FlowEdit(projectNo, flow, tagMap).Render(r.Context(), w)
	}

}

// /project/{projectNo}/flow/{flowNumber}/requirement/{reqIndex}/change
func HandleFlowRequirementChange(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkFlowNumber,
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))
				flowNo, _ := strconv.Atoi(r.PathValue("flowNumber"))

				flowObj, err := flows.Get(db, uint64(projectNo), uint64(flowNo))
				if err != nil {
					logrus.WithError(err).Warn("error loading flow for edit")
					pages.Error("Error parsing URL").Render(r.Context(), w)
					return
				}

				reqIndex, err := strconv.Atoi(r.PathValue("reqIndex"))
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error":      err,
						"path_value": "reqIndex",
					}).Warn("URL parsing error")
					pages.Error("Error parsing URL").Render(r.Context(), w)
					return
				}

				projectTags, err := tag.List(db, uint64(projectNo))
				if err != nil {
					logrus.WithError(err).Warn("error loading tags")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Server error").Render(r.Context(), w)
					return
				}

				tagMap := make(map[uint64]*pb_tag.Tag)
				for _, tagObj := range projectTags {
					tagMap[tagObj.GetNumber()] = tagObj
				}

				query := r.URL.Query()
				switch query.Get("action") {
				case "change-category":
					r.ParseForm()
					categoryrIdStr := r.Form.Get("req_category")
					categoryrId, err := strconv.Atoi(categoryrIdStr)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"error":      err,
							"path_value": "tpl_required_tag_id",
						}).Warn("URL parsing error")
						pages.Error("Error parsing URL").Render(r.Context(), w)
						return
					}
					req := flowObj.GetRequirements()[reqIndex]

					req.InCategory = pb_issue.IssueCategory(categoryrId)

					if err := flows.Update(db, uint64(projectNo), flowObj); err != nil {
						logrus.WithError(err).Warn("error updating gflow")
						w.WriteHeader(http.StatusInternalServerError)
						pages.Error("Server error").Render(r.Context(), w)
						return
					}

					pages.FlowRequirementEditEntry(
						reqIndex, uint64(projectNo),
						req,
						flowObj,
						tagMap,
					).Render(r.Context(), w)
				case "add-required-tag":
					r.ParseForm()
					tagNumberToAddInReqStr := r.Form.Get("tpl_required_tag_id")
					tagNumberToAddInReq, err := strconv.Atoi(tagNumberToAddInReqStr)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"error":      err,
							"path_value": "tpl_required_tag_id",
						}).Warn("URL parsing error")
						pages.Error("Error parsing URL").Render(r.Context(), w)
						return
					}

					req := flowObj.GetRequirements()[reqIndex]
					req.RequiredTagIds = append(req.RequiredTagIds, uint64(tagNumberToAddInReq))

					if err := flows.Update(db, uint64(projectNo), flowObj); err != nil {
						logrus.WithError(err).Warn("error updating gflow")
						w.WriteHeader(http.StatusInternalServerError)
						pages.Error("Server error").Render(r.Context(), w)
						return
					}

					pages.FlowRequirementEditEntry(
						reqIndex, uint64(projectNo),
						req,
						flowObj,
						tagMap,
					).Render(r.Context(), w)
				case "add-prohibited-tag":
					r.ParseForm()
					tagNumberToAddInReqStr := r.Form.Get("tpl_prohibited_tag_id")
					tagNumberToAddInReq, err := strconv.Atoi(tagNumberToAddInReqStr)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"error":      err,
							"path_value": "tpl_required_tag_id",
						}).Warn("URL parsing error")
						pages.Error("Error parsing URL").Render(r.Context(), w)
						return
					}

					req := flowObj.GetRequirements()[reqIndex]
					req.ProhibitedTagIds = append(req.ProhibitedTagIds, uint64(tagNumberToAddInReq))

					if err := flows.Update(db, uint64(projectNo), flowObj); err != nil {
						logrus.WithError(err).Warn("error updating gflow")
						w.WriteHeader(http.StatusInternalServerError)
						pages.Error("Server error").Render(r.Context(), w)
						return
					}

					pages.FlowRequirementEditEntry(
						reqIndex, uint64(projectNo),
						req,
						flowObj,
						tagMap,
					).Render(r.Context(), w)
				case "delete-required-tag":
					tagNumberToRemove := query.Get("tag-number")
					logrus.WithFields(logrus.Fields{
						"project_no": projectNo,
						"flow_no":    flowNo,
						"tag_no":     tagNumberToRemove,
						"req_index":  reqIndex,
					}).Info("remove required tag")
				case "delete-prohibited-tag":
					tagNumberToRemove := query.Get("tag-number")
					logrus.WithFields(logrus.Fields{
						"project_no": projectNo,
						"flow_no":    flowNo,
						"tag_no":     tagNumberToRemove,
						"req_index":  reqIndex,
					}).Info("remove required tag")
				case "new":
					nextReqIndex := len(flowObj.GetRequirements())

					flowObj.Requirements = append(flowObj.Requirements, &flow.FlowRequirement{})
					if err := flows.Update(db, uint64(projectNo), flowObj); err != nil {
						logrus.WithError(err).Warn("error updating gflow")
						w.WriteHeader(http.StatusInternalServerError)
						pages.Error("Server error").Render(r.Context(), w)
						return
					}

					pages.FlowRequirementEditEntry(
						nextReqIndex, uint64(projectNo),
						&flow.FlowRequirement{},
						flowObj,
						tagMap,
					).Render(r.Context(), w)
				}
			},
		)
	}
}
