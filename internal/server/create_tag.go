package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleTagCreate(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				// projectNo is already validated by authUserWithProjectNo so we dont need to error check again
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))

				if r.Method == "GET" {
					renderTagCreateForm(uint64(projectNo), w, r)
				}
				if r.Method == "POST" {
					handleNewTagSubmit(db, uint64(projectNo), w, r)
				}
			},
		)
	}
}

func renderTagCreateForm(projectNo uint64, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "" {
		pages.TagCreateBody(projectNo).Render(r.Context(), w)
	} else {
		pages.TagCreate(projectNo).Render(r.Context(), w)
	}
}

func handleNewTagSubmit(
	db *badger.DB,
	projectNo uint64,
	w http.ResponseWriter, r *http.Request,
) {
	r.ParseForm()
	tagTitle := r.Form.Get("title")
	if tagTitle == "" {
		handleInputErr("title", fmt.Errorf("no tag title"), w)
		return
	}

	if _, err := tag.Create(db, projectNo, &pb_tag.Tag{
		Title:    tagTitle,
		HexColor: r.Form.Get("color"),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "projectNo",
		}).Warn("error inserting new flow")
		w.WriteHeader(http.StatusInternalServerError)
		pages.Error("Error inserting flow").Render(r.Context(), w)
	}

	HandleTagsListing(db)(w, r)
}
