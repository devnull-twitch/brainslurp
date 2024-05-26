package server

import (
	"context"
	"net/http"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/tag"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleTagsListing(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectObj := getProjectFromContext(ctx)

				tagList, err := tag.List(db, projectObj.GetNumber())
				if err != nil {
					logrus.WithError(err).Warn("error loading tag list")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Internal server errorr").Render(r.Context(), w)
					return
				}

				if r.Header.Get("HX-Request") != "" {
					pages.BodyLogoOOB(projectObj.GetNumber()).Render(r.Context(), w)
					if r.Header.Get("HX-Target") == "flow-list" {
						pages.TagListItems(tagList).Render(r.Context(), w)
					} else {
						pages.TagListingBody(projectObj.GetNumber(), tagList).Render(r.Context(), w)
					}
				} else {
					pages.TagListing(projectObj.GetNumber(), tagList).Render(r.Context(), w)
				}
			},
		)
	}
}
