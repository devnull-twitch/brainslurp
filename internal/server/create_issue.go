package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleIssueCreate(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				// projectNo is already validated by authUserWithProjectNo so we dont need to error check again
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))

				if r.Method == "GET" {
					renderIssueCreateForm(uint64(projectNo), w, r)
				}
				if r.Method == "POST" {
					handleNewIssueSubmit(db, uint64(projectNo), w, r)
				}
			},
		)
	}
}

func renderIssueCreateForm(projectNo uint64, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "" {
		pages.IssueCreateBody(projectNo).Render(r.Context(), w)
	} else {
		pages.IssueCreate(projectNo).Render(r.Context(), w)
	}
}

func handleNewIssueSubmit(db *badger.DB, projectNo uint64, w http.ResponseWriter, r *http.Request) {
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

	if err := issues.Create(db, projectNo, &pb_issue.Issue{
		Title:    issueTitle,
		Body:     r.Form.Get("body"),
		Category: pb_issue.IssueCategory(catInt),
	}); err != nil {
		logrus.WithError(err).Warn("error saving new issue")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Error saving neew issue").Render(r.Context(), w)
	}

	header := w.Header()
	header.Set("HX-Push-Url", fmt.Sprintf("/project/%d/issues", projectNo))
	header.Set("HX-Retarget", "body")
	HandleIssueList(db)(w, r)
}
