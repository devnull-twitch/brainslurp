package server

import (
	"context"
	"net/http"

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
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				db := getDbFromContext(ctx)
				projectObj := getProjectFromContext(ctx)

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

				if r.Header.Get("HX-Request") != "" {
					pages.BodyLogoOOB(projectObj.GetNumber()).Render(r.Context(), w)
					if r.Header.Get("HX-Target") == "flow-list" {
						pages.FlowListItems(projectObj.GetNumber(), flows).Render(r.Context(), w)
					} else {
						pages.FlowListBody(projectObj.GetNumber(), flows).Render(r.Context(), w)
					}
				} else {
					pages.FlowList(projectObj.GetNumber(), flows).Render(r.Context(), w)
				}
			},
		)
	}
}
