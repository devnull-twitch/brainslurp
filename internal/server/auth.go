package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func authUserWithProjectNo(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
	db := getDbFromContext(ctx)
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

	projectNo, err := strconv.Atoi(r.PathValue("projectNo"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"path_value": "projectNo",
		}).Warn("URL parsing error")
		pages.Error("Error parsing URL").Render(r.Context(), w)
		return
	}

	projectObj := &pb_project.Project{}
	if err := db.View(func(txn *badger.Txn) error {
		projectKey := database.Keygen(database.ProjectPrefix, uint64(projectNo))
		projectItem, err := txn.Get(projectKey)
		if err != nil {
			return fmt.Errorf("error loading project %d: %w", projectNo, err)
		}

		if err := database.Read(projectItem, projectObj); err != nil {
			return fmt.Errorf("error reading project: %w", err)
		}

		for _, mem := range projectObj.GetMembers() {
			if mem.GetUserNo() == userObj.GetNumber() {
				return nil
			}
		}

		return fmt.Errorf("unauthorized")
	}); err != nil {
		logrus.WithError(err).Error("error authenticating user")
		pages.Error("Halp! Something went wrong").Render(r.Context(), w)
		return
	}

	next(setProjectOnContext(setUserOnContext(ctx, userObj), projectObj))
}

func authUser(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
	db := getDbFromContext(ctx)
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

	next(setUserOnContext(ctx, userObj))
}
