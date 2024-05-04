package user

import (
	"encoding/binary"
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

func AddProject(db *badger.DB, userNo uint64, projectNo uint64) error {
	userObj, err := getProtoUser(db, userNo)
	if err != nil {
		return fmt.Errorf("unable to get user for update: %w", err)
	}

	userObj.Memberships = append(userObj.Memberships, &pb_user.Membership{
		ProjectNo: projectNo,
	})

	userKey := make([]byte, binary.MaxVarintLen64+1)
	userKey[0] = database.UserPrefix
	binary.PutUvarint(userKey[1:], userNo)

	userVal, err := proto.Marshal(userObj)
	if err != nil {
		return fmt.Errorf("unable to marshal user: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(userKey, userVal)
	}); err != nil {
		return fmt.Errorf("error updaing user: %w", err)
	}

	return nil
}
