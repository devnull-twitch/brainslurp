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

type CreateOptions struct {
	Title      string
	Body       string
	Category   pb_issue.IssueCategory
	ViewStatus []uint64
	ProjectNo  uint64
	Tags       []*pb_issue.Tag
}

func Create(db *badger.DB, opts CreateOptions) error {
	issueNo, err := database.NextNumber(db, database.IssueSequenceKey)
	if err != nil {
		return fmt.Errorf("error getting new issue number: %w", err)
	}

	if opts.ProjectNo <= 0 {
		return fmt.Errorf("no project number")
	}

	issueKeyLength := (2 * binary.MaxVarintLen64) + 1
	issueKey := make([]byte, issueKeyLength)
	issueKey[0] = database.IssuePrefix
	binary.PutUvarint(issueKey[1:], opts.ProjectNo)
	binary.PutUvarint(issueKey[binary.MaxVarintLen64+1:], issueNo)

	unixNow := time.Now().Unix()

	viewStatus := make([]*pb_issue.ViewStatus, len(opts.ViewStatus))
	viewIssueListKeys := make([][]byte, 0)
	for i, vsID := range opts.ViewStatus {
		viewStatus[i] = &pb_issue.ViewStatus{Number: vsID, SetAt: unixNow}

		viewIssueKey := make([]byte, issueKeyLength)
		viewIssueKey[0] = database.ViewIssuesPrefix
		binary.PutUvarint(viewIssueKey[1:], opts.ProjectNo)
		binary.PutUvarint(viewIssueKey[binary.MaxVarintLen64+1:], vsID)

		viewIssueListKeys = append(viewIssueListKeys, viewIssueKey)
	}

	issueVal, err := proto.Marshal(&pb_issue.Issue{
		Number:    issueNo,
		CreatedAt: unixNow,
		Title:     opts.Title,
		Body:      opts.Body,
		Category:  opts.Category,
		Views:     viewStatus,
		Tags:      opts.Tags,
	})
	if err != nil {
		return fmt.Errorf("unable to marshal issue: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(issueKey, issueVal)
	}); err != nil {
		return fmt.Errorf("unable to insert issue: %w", err)
	}

	logrus.WithField("key", fmt.Sprintf("%x", issueKey)).Info("issue inserted")

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
