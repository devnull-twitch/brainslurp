package tag

import (
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	"github.com/dgraph-io/badger/v4"
)

func Get(db *badger.DB, projectNo uint64, tagNumber uint64) (*pb_tag.Tag, error) {
	if projectNo <= 0 {
		return nil, fmt.Errorf("no project number")
	}

	tagKey := database.Keygen(database.TagPrefix, projectNo, tagNumber)
	tagObj := &pb_tag.Tag{}
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(tagKey)
		if err != nil {
			return err
		}

		return database.Read(item, tagObj)
	}); err != nil {
		return nil, err
	}

	return tagObj, nil
}

func GetMany(db *badger.DB, projectNo uint64, tagNumbers []uint64) ([]*pb_tag.Tag, error) {
	if projectNo <= 0 {
		return nil, fmt.Errorf("no project number")
	}

	tags := make([]*pb_tag.Tag, 0)
	for _, tagNum := range tagNumbers {
		tagKey := database.Keygen(database.TagPrefix, projectNo, tagNum)
		tagObj := &pb_tag.Tag{}
		if err := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get(tagKey)
			if err != nil {
				return err
			}

			if err := database.Read(item, tagObj); err != nil {
				return err
			}

			tags = append(tags, tagObj)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return tags, nil
}
