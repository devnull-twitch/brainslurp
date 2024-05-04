package server

import (
	"encoding/binary"
	"encoding/hex"
	"net/http"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/project"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/dgraph-io/badger/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var SigningSecret = []byte("ThisIsVerySecret")

func LoginHandler(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			userObj, err := user.Login(db, r.FormValue("username"), r.FormValue("password"))
			if err != nil {
				// at least handle username not found in a nicer way
				// also. maybe check against timing attacks?? although low prio. because this might never see a prod server so ...
				logrus.WithError(err).Error("error performing login")
				pages.Error("Error").Render(r.Context(), w)
				return
			}

			noBuf := make([]byte, binary.MaxVarintLen64)
			binary.PutUvarint(noBuf, userObj.Number)
			noAsStr := hex.EncodeToString(noBuf)

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"no":       noAsStr,
				"username": userObj.Username,
			})
			tokenStr, err := token.SignedString(SigningSecret)
			if err != nil {
				logrus.WithError(err).Error("error signing jwt")
				pages.Error("JWT signing error").Render(r.Context(), w)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "jwt",
				Value: tokenStr,
			})

			projects := make([]project.ListItem, len(userObj.Projects))
			for i, projectNo := range userObj.Projects {
				projects[i], err = project.Get(db, projectNo)
				if err != nil {
					logrus.WithError(err).Error("unable to load user project")
					continue
				}
			}

			w.Header().Add("HX-Push-Url", "/projects")
			pages.ProjectListingBody(projects).Render(r.Context(), w)
		} else {
			pages.Login().Render(r.Context(), w)
		}
	}
}
