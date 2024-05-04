package issues

import (
	"encoding/binary"
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

func Get(db *badger.DB, projectNo uint64, issueNo uint64) (*pb_issue.Issue, []*pb_flow.Flow, error) {
	if projectNo <= 0 {
		return nil, nil, fmt.Errorf("no project number")
	}

	if issueNo <= 0 {
		return nil, nil, fmt.Errorf("no issue number")
	}

	keyLength := (2 * binary.MaxVarintLen64) + 1
	issueKey := make([]byte, keyLength)
	issueKey[0] = database.IssuePrefix
	binary.PutUvarint(issueKey[1:], projectNo)
	binary.PutUvarint(issueKey[binary.MaxVarintLen64+1:], issueNo)

	issueObj := &pb_issue.Issue{}
	var flows []*pb_flow.Flow
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(issueKey)
		if err != nil {
			return fmt.Errorf("error getting issue entry: %w", err)
		}

		issueData, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("error copying issue data: %w", err)
		}

		if err := proto.Unmarshal(issueData, issueObj); err != nil {
			return fmt.Errorf("error unmarshaling issue: %w", err)
		}

		flows, err = fetchFlows(txn, projectNo, issueObj)
		if err != nil {
			return fmt.Errorf("error loading issue flows: %w", err)
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return issueObj, flows, nil
}
