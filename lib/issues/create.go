package issues

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func Create(db *badger.DB, projectNo uint64, newIssue *pb_issue.Issue) error {
	issueNo, err := database.NextNumber(db, database.IssueSequenceKey)
	if err != nil {
		return fmt.Errorf("error getting new issue number: %w", err)
	}

	if projectNo <= 0 {
		return fmt.Errorf("no project number")
	}

	issueKeyLength := (2 * binary.MaxVarintLen64) + 1
	issueKey := make([]byte, issueKeyLength)
	issueKey[0] = database.IssuePrefix
	binary.PutUvarint(issueKey[1:], projectNo)
	binary.PutUvarint(issueKey[binary.MaxVarintLen64+1:], issueNo)

	newIssue.Number = issueNo
	newIssue.CreatedAt = time.Now().Unix()

	issueVal, err := proto.Marshal(newIssue)
	if err != nil {
		return fmt.Errorf("unable to marshal issue: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(issueKey, issueVal)
	}); err != nil {
		return fmt.Errorf("unable to insert issue: %w", err)
	}

	logrus.WithField("key", fmt.Sprintf("%x", issueKey)).Info("issue inserted")

	viewIssueListKeys := make([][]byte, 0)
	for _, viewEntry := range newIssue.GetViews() {
		viewIsssueListeKey := database.Keygen(database.ViewIssuesPrefix, projectNo, viewEntry.GetNumber())
		viewIssueListKeys = append(viewIssueListKeys, viewIsssueListeKey)
	}
	for _, viewIssueKey := range viewIssueListKeys {
		merger := db.GetMergeOperator(viewIssueKey, func(existingVal, newVal []byte) []byte {
			return append(existingVal, newVal...)
		}, time.Second)
		merger.Add(issueKey)

		// Get applies merges and saves in db
		_, err := merger.Get()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"errors":        err,
				"issueview_key": viewIssueKey,
			}).Warn("unable to add new issue ID to view list")
		}
		merger.Stop()
	}

	return nil
}
