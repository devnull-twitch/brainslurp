package project

import (
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	"github.com/devnull-twitch/brainslurp/lib/proto/project"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

func AddUser(db *badger.DB, userNo uint64, projectNo uint64) error {
	projectObj, err := Get(db, projectNo)
	if err != nil {
		return fmt.Errorf("unable to get user for update: %w", err)
	}

	projectObj.Members = append(projectObj.Members, &project.Member{
		UserNo: userNo,
	})

	projectKey := database.Keygen(database.ProjectPrefix, projectNo)
	projectVal, err := proto.Marshal(projectObj)
	if err != nil {
		return fmt.Errorf("unable to marshal project: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(projectKey, projectVal)
	}); err != nil {
		return fmt.Errorf("error updaing user: %w", err)
	}

	return nil
}
