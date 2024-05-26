package server

import (
	"context"
	"fmt"
	"net/http"

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
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				if r.Method == "GET" {
					renderTagCreateForm(ctx, w, r)
				}
				if r.Method == "POST" {
					handleNewTagSubmit(ctx, w, r)
				}
			},
		)
	}
}

func renderTagCreateForm(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	projectObj := getProjectFromContext(ctx)

	if r.Header.Get("HX-Request") != "" {
		pages.TagCreateBody(projectObj.GetNumber()).Render(r.Context(), w)
	} else {
		pages.TagCreate(projectObj.GetNumber()).Render(r.Context(), w)
	}
}

func handleNewTagSubmit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db := getDbFromContext(ctx)
	projectObj := getProjectFromContext(ctx)

	r.ParseForm()
	tagTitle := r.Form.Get("title")
	if tagTitle == "" {
		handleInputErr("title", fmt.Errorf("no tag title"), w)
		return
	}

	if _, err := tag.Create(db, projectObj.GetNumber(), &pb_tag.Tag{
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
