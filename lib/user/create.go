package user

import (
	"fmt"

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

func Create(db *badger.DB, password string, newUserObj *pb_user.User) (uint64, error) {
	userNo, err := database.NextNumber(db, database.UserSequenceKey)
	if err != nil {
		return 0, fmt.Errorf("error getting new user number: %w", err)
	}
	newUserObj.Number = userNo

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, fmt.Errorf("error hasing password: %w", err)
	}
	newUserObj.PasswordHash = hash

	userKey := database.Keygen(database.UserPrefix, userNo)
	userVal, err := proto.Marshal(newUserObj)
	if err != nil {
		return 0, fmt.Errorf("unable to marshal user: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		usernameBytes := []byte(newUserObj.Name)
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
