package project

import (
	"fmt"
	"time"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func Update(db *badger.DB, projectObj *pb_project.Project) error {
	projectKey := database.Keygen(database.ProjectPrefix, projectObj.GetNumber())

	projectObj.UpdatedAt = time.Now().Unix()

	projectVal, err := proto.Marshal(projectObj)
	if err != nil {
		return fmt.Errorf("unable to marshal project: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(projectKey, projectVal)
	}); err != nil {
		return fmt.Errorf("unable to insert project: %w", err)
	}

	logrus.WithField("key", fmt.Sprintf("%x", projectKey)).Info("project updated")

	return nil
}
