package tag

import (
	"encoding/binary"
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	"github.com/dgraph-io/badger/v4"
)

func List(db *badger.DB, projectNo uint64) ([]*pb_tag.Tag, error) {
	if projectNo <= 0 {
		return nil, fmt.Errorf("no project number")
	}

	tagList := make([]*pb_tag.Tag, 0)
	if err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := make([]byte, binary.MaxVarintLen64+1)
		prefix[0] = database.TagPrefix
		binary.PutUvarint(prefix[1:], projectNo)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			tagObj := &pb_tag.Tag{}
			if err := database.Read(item, tagObj); err != nil {
				return err
			}

			tagList = append(tagList, tagObj)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return tagList, nil
}
