package issues

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

func ListFromView(db *badger.DB, projectNo uint64, viewNo uint64) ([]*pb_issue.Issue, map[uint64][]*pb_flow.Flow, error) {
	if projectNo <= 0 {
		return nil, nil, fmt.Errorf("no project number")
	}

	viewIsssueKeyLength := (2 * binary.MaxVarintLen64) + 1
	viewIssueKey := make([]byte, viewIsssueKeyLength)
	viewIssueKey[0] = database.ViewIssuesPrefix
	binary.PutUvarint(viewIssueKey[1:], projectNo)
	binary.PutUvarint(viewIssueKey[binary.MaxVarintLen64+1:], viewNo)

	issueList := make([]*pb_issue.Issue, 0)
	flowToIssueMap := make(map[uint64][]*pb_flow.Flow)
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(viewIssueKey)
		if err != nil {
			return fmt.Errorf("error getting view list: %w", err)
		}

		idList, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("error copying view list data: %w", err)
		}

		r := bytes.NewReader(idList)
		// key is uint64 + 1 byte as a type
		keyBuf := make([]byte, 2*binary.MaxVarintLen64+1)
		for r.Len() > 0 {
			_, err := r.Read(keyBuf)
			if err != nil {
				return fmt.Errorf("error reading issue ID: %w", err)
			}

			issueItem, err := txn.Get(keyBuf)
			if err != nil {
				return fmt.Errorf("error getting issue %x: %w", keyBuf, err)
			}

			issueData, err := issueItem.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("error copying issue data: %w", err)
			}

			issue := &pb_issue.Issue{}
			if err := proto.Unmarshal(issueData, issue); err != nil {
				return fmt.Errorf("error unmarshaling issue: %w", err)
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
