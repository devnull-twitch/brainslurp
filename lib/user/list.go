package user

import (
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/dgraph-io/badger/v4"
)

func List(db *badger.DB, userNumbers []uint64) ([]*pb_user.User, error) {
	users := make([]*pb_user.User, 0)
	if err := db.View(func(txn *badger.Txn) error {
		for _, userNo := range userNumbers {
			key := database.Keygen(database.UserPrefix, userNo)

			userItem, err := txn.Get(key)
			if err != nil {
				return fmt.Errorf("unable to load user: %w", err)
			}

			userObj := &pb_user.User{}
			if err := database.Read(userItem, userObj); err != nil {
				return fmt.Errorf("unable to read user: %w", err)
			}

			users = append(users, userObj)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return users, nil
}

func ListAll(db *badger.DB) ([]*pb_user.User, error) {
	users := make([]*pb_user.User, 0)
	if err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte{database.UserPrefix}

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			userObj := &pb_user.User{}
			if err := database.Read(item, userObj); err != nil {
				return err
			}

			users = append(users, userObj)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return users, nil
}
