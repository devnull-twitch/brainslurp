package view

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/devnull-twitch/brainslurp/lib/database"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func SetIssue(db *badger.DB, projectNo uint64, issueNo uint64, viewNos []uint64) error {
	viewIssuePrefix := database.Keygen(database.ViewIssuesPrefix, projectNo)
	insertKey := database.Keygen(database.IssuePrefix, projectNo, issueNo)

	return db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		for it.Seek(viewIssuePrefix); it.ValidForPrefix(viewIssuePrefix); it.Next() {
			item := it.Item()

			viewListKey := item.Key()
			viewNoByte := make([]byte, binary.MaxVarintLen64)
			copy(viewNoByte, viewListKey[1:1+binary.MaxVarintLen64])
			viewNo, err := binary.ReadUvarint(bytes.NewReader(viewNoByte))
			if err != nil {
				return fmt.Errorf("")
			}

			viewListData, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("error loading view list content: %w", err)
			}

			hasChanged := false
			if slices.Contains(viewNos, viewNo) {
				hasChanged, err = tryAdd(viewListData, insertKey)
			} else {
				hasChanged, err = tryRemove(viewListData, insertKey)
			}
			if err != nil {
				return fmt.Errorf("error adding or removing issue: %w", err)
			}

			if hasChanged {
				err = txn.Set(viewListKey, viewListData)
				logrus.WithField("viewlist_key", fmt.Sprintf("%x", viewListKey)).Info("updated view list")
				if err != nil {
					return fmt.Errorf("error saving updated issue list: %w", err)
				}
			}
		}

		return nil
	})
}

func tryAdd(viewListData []byte, insertKey []byte) (bool, error) {
	keyLength := 2*binary.MaxVarintLen64 + 1

	r := bytes.NewReader(viewListData)
	keyBuf := make([]byte, keyLength)
	for r.Len() > 0 {
		_, err := r.Read(keyBuf)
		if err != nil {
			return false, fmt.Errorf("error reading issue ID from view: %w", err)
		}

		if slices.Equal(keyBuf, insertKey) {
			logrus.WithFields(logrus.Fields{
				"issue_key": insertKey,
			}).Debug("issue was already in view")
			return false, nil
		}
	}

	slices.Grow(viewListData, keyLength)
	startIndex := len(viewListData) - keyLength
	copy(viewListData[startIndex:keyLength], insertKey)

	return true, nil
}

func tryRemove(viewListData []byte, insertKey []byte) (bool, error) {
	keyLength := 2*binary.MaxVarintLen64 + 1

	r := bytes.NewReader(viewListData)
	keyBuf := make([]byte, keyLength)
	for r.Len() > 0 {
		_, err := r.Read(keyBuf)
		if err != nil {
			return false, fmt.Errorf("error reading issue ID from view: %w", err)
		}

		if slices.Equal(keyBuf, insertKey) {
			start := int(r.Size()) - r.Len()
			slices.Delete(viewListData, start, start+keyLength)
			return true, nil
		}
	}

	return false, nil
}
