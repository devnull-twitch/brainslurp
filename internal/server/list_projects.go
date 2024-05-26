package server

import (
	"context"
	"net/http"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/project"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleProjectListing(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUser,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				userObj := getUserFromContext(ctx)

				projects := make([]*pb_project.Project, len(userObj.GetMemberships()))
				var err error
				for i, userMembership := range userObj.GetMemberships() {
					projects[i], err = project.Get(db, userMembership.GetProjectNo())
					if err != nil {
						logrus.WithError(err).Error("unable to load user project")
						continue
					}
				}

				if r.Header.Get("HX-Request") != "" {
					pages.BodyLogoOOB(0).Render(r.Context(), w)
					pages.ProjectListingBody(projects).Render(r.Context(), w)
				} else {
					pages.ProjectListing(projects).Render(r.Context(), w)
				}
			},
		)
	}
}
