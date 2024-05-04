package server

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/dgraph-io/badger/v4"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoCookie = errors.New("no cookie")
)

func parseAuthAndLoadUser(db *badger.DB, r *http.Request) (user.User, error) {
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		return user.User{}, ErrNoCookie
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(jwtCookie.Value, &claims, func(t *jwt.Token) (interface{}, error) {
		return SigningSecret, nil
	})
	if err != nil {
		return user.User{}, fmt.Errorf("error parsing jwt: %w", err)
	}

	noBuf, err := hex.DecodeString(claims["no"].(string))
	if err != nil {
		return user.User{}, fmt.Errorf("error jwt doenst contain user no in hex format: %w", err)
	}
	userNo, err := binary.ReadUvarint(bytes.NewReader(noBuf))
	if err != nil {
		return user.User{}, fmt.Errorf("error jwt no isnt usigned 64 int: %w", err)
	}

	userObj, err := user.Get(db, userNo)
	if err != nil {
		return user.User{}, fmt.Errorf("error user could not be loaded: %w", err)
	}

	return userObj, nil
}
