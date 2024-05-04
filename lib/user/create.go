package user

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

type CreateOptions struct {
	Username string
	Password string
}

func Create(db *badger.DB, opts CreateOptions) (uint64, error) {
	userNo, err := database.NextNumber(db, database.UserSequenceKey)
	if err != nil {
		return 0, fmt.Errorf("error getting new user number: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(opts.Password), 12)
	if err != nil {
		return 0, fmt.Errorf("error hasing password: %w", err)
	}

	userKey := make([]byte, binary.MaxVarintLen64+1)
	userKey[0] = database.UserPrefix
	binary.PutUvarint(userKey[1:], userNo)

	userVal, err := proto.Marshal(&pb_user.User{
		Number:       userNo,
		CreatedAt:    time.Now().Unix(),
		Name:         opts.Username,
		PasswordHash: hash,
	})
	if err != nil {
		return 0, fmt.Errorf("unable to marshal user: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		usernameBytes := []byte(opts.Username)
		usernameKey := make([]byte, len(usernameBytes)+1)
		usernameKey[0] = database.UsernameLookupPrefix
		copy(usernameKey[1:], usernameBytes)
		txn.Set(usernameKey, userKey)

		logrus.WithField("key", fmt.Sprintf("%x", usernameKey)).Info("username lookup inserted")

		return txn.Set(userKey, userVal)
	}); err != nil {
		return 0, fmt.Errorf("unable to insert user: %w", err)
	}

	logrus.WithField("key", fmt.Sprintf("%x", userKey)).Info("user inserted")

	return userNo, nil
}
