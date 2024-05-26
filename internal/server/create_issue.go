package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleIssueCreate(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				if r.Method == "GET" {
					renderIssueCreateForm(ctx, w, r)
				}
				if r.Method == "POST" {
					handleNewIssueSubmit(ctx, w, r)
				}
			},
		)
	}
}

func renderIssueCreateForm(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db := getDbFromContext(ctx)
	projectObj := getProjectFromContext(ctx)

	projectTags, err := tag.List(db, projectObj.GetNumber())
	if err != nil {
		logrus.WithError(err).Warn("error loading tags")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Server error").Render(r.Context(), w)
	}

	if r.Header.Get("HX-Request") != "" {
		pages.IssueCreateBody(projectObj.GetNumber(), projectTags).Render(r.Context(), w)
	} else {
		pages.IssueCreate(projectObj.GetNumber(), projectTags).Render(r.Context(), w)
	}
}

func handleNewIssueSubmit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db := getDbFromContext(ctx)
	projectObj := getProjectFromContext(ctx)

	r.ParseForm()
	issueTitle := r.Form.Get("title")
	if issueTitle == "" {
		handleInputErr("title", fmt.Errorf("no issue title"), w)
		return
	}

	catStr := r.Form.Get("category")
	if catStr == "" {
		handleInputErr("category", fmt.Errorf("no issue category"), w)
		return
	}
	catInt, err := strconv.Atoi(catStr)
	if err != nil {
		logrus.WithError(err).Warn("error converting category to int")
		handleInputErr("category", fmt.Errorf("no issue category"), w)
		return
	}

	tagStrNos := r.Form["tag"]
	tagNos := StrToIntSlice[uint64](tagStrNos)

	if err := issues.Create(db, projectObj.GetNumber(), &pb_issue.Issue{
		Title:      issueTitle,
		Body:       r.Form.Get("body"),
		Category:   pb_issue.IssueCategory(catInt),
		TagNumbers: tagNos,
	}); err != nil {
		logrus.WithError(err).Warn("error saving new issue")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Error saving neew issue").Render(r.Context(), w)
	}

	header := w.Header()
	header.Set("HX-Push-Url", fmt.Sprintf("/project/%d/issues", projectObj.GetNumber()))
	header.Set("HX-Retarget", "body")
	HandleIssueList(db)(w, r)
}
