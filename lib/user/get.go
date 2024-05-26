package user

import (
	"encoding/binary"
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

func Get(db *badger.DB, no uint64) (*pb_user.User, error) {
	user := &pb_user.User{}
	if err := db.View(func(txn *badger.Txn) error {
		userKey := make([]byte, binary.MaxVarintLen64+1)
		userKey[0] = database.UserPrefix
		binary.PutUvarint(userKey[1:], no)

		userItem, err := txn.Get(userKey)
		if err != nil {
			return fmt.Errorf("unable to get user: %w", err)
		}

		userData, err := userItem.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("error copying user data: %w", err)
		}

		if err := proto.Unmarshal(userData, user); err != nil {
			return fmt.Errorf("error unmarshaling user: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return user, nil
}
