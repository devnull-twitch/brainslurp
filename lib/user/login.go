package user

import (
	"errors"
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/dgraph-io/badger/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

var ErrUserNotFound = errors.New("User not found")

type User struct {
	Number   uint64
	Username string
	Projects []uint64
}

func Login(db *badger.DB, username string, password string) (User, error) {
	user := &pb_user.User{}
	if err := db.View(func(txn *badger.Txn) error {
		usernameBytes := []byte(username)
		usernameKey := make([]byte, len(usernameBytes)+1)
		usernameKey[0] = database.UsernameLookupPrefix
		copy(usernameKey[1:], usernameBytes)

		keyItem, err := txn.Get(usernameKey)
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return ErrUserNotFound
			}
			return fmt.Errorf("unable to get username lookup: %w", err)
		}

		userKey, err := keyItem.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("error copying user key: %w", err)
		}

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
		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword(user.GetPasswordHash(), []byte(password)); err != nil {
		return User{}, fmt.Errorf("password mismatch: %w", err)
	}

	projectNos := make([]uint64, len(user.GetMemberships()))
	for i, projectEntry := range user.GetMemberships() {
		projectNos[i] = projectEntry.GetProjectNo()
	}

	return User{
		Number:   user.GetNumber(),
		Username: user.GetName(),
		Projects: projectNos,
	}, nil
}
