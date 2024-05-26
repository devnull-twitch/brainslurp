package project

import (
	"encoding/binary"
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

func Get(db *badger.DB, no uint64) (*pb_project.Project, error) {
	if no <= 0 {
		return nil, fmt.Errorf("invalid project number")
	}

	projectKey := make([]byte, binary.MaxVarintLen64+1)
	projectKey[0] = database.ProjectPrefix
	binary.PutUvarint(projectKey[1:], no)

	project := &pb_project.Project{}
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(projectKey)
		if err != nil {
			return err
		}

		projectData, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("error copying view data: %w", err)
		}

		if err := proto.Unmarshal(projectData, project); err != nil {
			return fmt.Errorf("error unmarshaling view: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error getting project: %w", err)
	}

	return project, nil
}
