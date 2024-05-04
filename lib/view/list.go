package view

import (
	"encoding/binary"
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_view "github.com/devnull-twitch/brainslurp/lib/proto/view"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

func List(db *badger.DB, projectNo uint64) ([]*pb_view.View, error) {
	if projectNo <= 0 {
		return nil, fmt.Errorf("no project number")
	}

	viewList := make([]*pb_view.View, 0)
	if err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := make([]byte, binary.MaxVarintLen64+1)
		prefix[0] = database.ViewPrefix
		binary.PutUvarint(prefix[1:], projectNo)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			viewData, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("error copying view data: %w", err)
			}

			view := &pb_view.View{}
			if err := proto.Unmarshal(viewData, view); err != nil {
				return fmt.Errorf("error unmarshaling view: %w", err)
			}

			viewList = append(viewList, view)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return viewList, nil
}
