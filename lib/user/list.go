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
