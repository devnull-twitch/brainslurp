package server

import (
	"errors"
	"net/http"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/project"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleProjectListing(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userObj, err := parseAuthAndLoadUser(db, r)
		if err != nil {
			if errors.Is(err, ErrNoCookie) {
				// maybe if user is not logged in redirect to start page?
				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}

			logrus.WithError(err).Error("error authenticating user")
			pages.Error("Halp! Something went wrong").Render(r.Context(), w)
			return
		}

		projects := make([]project.ListItem, len(userObj.Projects))
		for i, projectNo := range userObj.Projects {
			projects[i], err = project.Get(db, projectNo)
			if err != nil {
				logrus.WithError(err).Error("unable to load user project")
				continue
			}
		}

		pages.ProjectListing(projects).Render(r.Context(), w)
	}
}
