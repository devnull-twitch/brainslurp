package server

import (
	"net/http"
	"strconv"

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
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))
				tagList, err := tag.List(db, uint64(projectNo))
				if err != nil {
					logrus.WithError(err).Warn("error loading tag list")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Internal server errorr").Render(r.Context(), w)
					return
				}

				if r.Header.Get("HX-Request") != "" {
					if r.Header.Get("HX-Target") == "flow-list" {
						pages.TagListItems(tagList).Render(r.Context(), w)
					} else {
						pages.TagListingBody(uint64(projectNo), tagList).Render(r.Context(), w)
					}
				} else {
					pages.TagListing(uint64(projectNo), tagList).Render(r.Context(), w)
				}
			},
		)
	}
}
