package tag

import (
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func Create(db *badger.DB, projectNo uint64, tagObj *pb_tag.Tag) (uint64, error) {
	tagNo, err := database.NextNumber(db, database.TagSequenceKey)
	if err != nil {
		return 0, fmt.Errorf("error getting new view number: %w", err)
	}

	if projectNo <= 0 {
		return 0, fmt.Errorf("no project ID given")
	}

	tagObj.Number = tagNo
	tagVal, err := proto.Marshal(tagObj)
	if err != nil {
		return 0, fmt.Errorf("unable to marshal tag: %w", err)
	}

	newKey := database.Keygen(database.TagPrefix, projectNo, tagNo)

	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(newKey, tagVal)
	}); err != nil {
		return 0, fmt.Errorf("error inserting new tag: %w", err)
	}

	logrus.WithField("key", fmt.Sprintf("%x", newKey)).Info("tag created")

	return tagNo, nil
}
