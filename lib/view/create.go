package view

import (
	"fmt"
	"time"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_view "github.com/devnull-twitch/brainslurp/lib/proto/view"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type CreateOptions struct {
	ProjectNo uint64
	Title     string
}

func Create(db *badger.DB, opts CreateOptions) (uint64, error) {
	viewNo, err := database.NextNumber(db, database.ViewSequenceKey)
	if err != nil {
		return 0, fmt.Errorf("error getting new view number: %w", err)
	}

	if opts.ProjectNo <= 0 {
		return 0, fmt.Errorf("no project ID given")
	}

	viewKey := database.Keygen(database.ViewPrefix, opts.ProjectNo, viewNo)
	viewIssueKey := database.Keygen(database.ViewIssuesPrefix, opts.ProjectNo, viewNo)

	viewVal, err := proto.Marshal(&pb_view.View{
		Number:    viewNo,
		CreatedAt: time.Now().Unix(),
		Title:     opts.Title,
	})
	if err != nil {
		return 0, fmt.Errorf("unable to marshal issue: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(viewIssueKey, []byte{}); err != nil {
			return fmt.Errorf("error creating issue list entry: %w", err)
		}
		return txn.Set(viewKey, viewVal)
	}); err != nil {
		return 0, fmt.Errorf("error persisting view: %w", err)
	}

	logrus.WithField("key", fmt.Sprintf("%x", viewKey)).Info("view inserted")
	logrus.WithField("key", fmt.Sprintf("%x", viewIssueKey)).Info("view issue list started")

	return viewNo, nil
}
