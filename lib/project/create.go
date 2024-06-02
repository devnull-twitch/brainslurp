package project

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type CreateOptions struct {
	Title  string
	Public bool
}

func Create(db *badger.DB, opts CreateOptions) (uint64, error) {
	projectNo, err := database.NextNumber(db, database.ProjectSequenceKey)
	if err != nil {
		return 0, fmt.Errorf("error getting new project number: %w", err)
	}

	projectKey := make([]byte, binary.MaxVarintLen64+1)
	projectKey[0] = database.ProjectPrefix
	binary.PutUvarint(projectKey[1:], projectNo)

	nowUnix := time.Now().Unix()

	projectVal, err := proto.Marshal(&pb_project.Project{
		Number:    projectNo,
		CreatedAt: nowUnix,
		UpdatedAt: nowUnix,
		Name:      opts.Title,
		Public:    opts.Public,
		Members:   []*pb_project.Member{},
	})
	if err != nil {
		return 0, fmt.Errorf("unable to marshal project: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(projectKey, projectVal)
	}); err != nil {
		return 0, fmt.Errorf("unable to insert project: %w", err)
	}

	logrus.WithField("key", fmt.Sprintf("%x", projectKey)).Info("project inserted")

	return projectNo, nil
}
