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

func ListAll(db *badger.DB, projectNo uint64) ([]*pb_issue.Issue, map[uint64][]*pb_flow.Flow, error) {
	if projectNo <= 0 {
		return nil, nil, fmt.Errorf("no project number")
	}

	issueList := make([]*pb_issue.Issue, 0)
	flowToIssueMap := make(map[uint64][]*pb_flow.Flow)
	if err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := make([]byte, binary.MaxVarintLen64+1)
		prefix[0] = database.IssuePrefix
		binary.PutUvarint(prefix[1:], projectNo)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			issueData, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("error copying view data: %w", err)
			}

			issue := &pb_issue.Issue{}
			if err := proto.Unmarshal(issueData, issue); err != nil {
				return fmt.Errorf("error unmarshaling view: %w", err)
			}

			flowList, err := fetchFlows(txn, projectNo, issue)
			if err != nil {
				return fmt.Errorf("error fetching connected flows: %w", err)
			}
			flowToIssueMap[issue.GetNumber()] = flowList

			issueList = append(issueList, issue)
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return issueList, flowToIssueMap, nil
}
