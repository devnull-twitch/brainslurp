package server

import (
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleFlowList(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			func(db *badger.DB, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectNo, _ := strconv.Atoi(r.PathValue("projectNo"))
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

				if r.Header.Get("HX-Request") != "" {
					pages.BodyLogoOOB(uint64(projectNo)).Render(r.Context(), w)
					if r.Header.Get("HX-Target") == "flow-list" {
						pages.FlowListItems(uint64(projectNo), flows).Render(r.Context(), w)
					} else {
						pages.FlowListBody(uint64(projectNo), flows).Render(r.Context(), w)
					}
				} else {
					pages.FlowList(uint64(projectNo), flows).Render(r.Context(), w)
				}
			},
		)
	}
}
