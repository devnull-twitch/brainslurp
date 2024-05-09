package flows

import (
	"encoding/binary"
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

func List(db *badger.DB, projectNo uint64) ([]*pb_flow.Flow, error) {
	if projectNo <= 0 {
		return nil, fmt.Errorf("no project number")
	}

	keyLength := binary.MaxVarintLen64 + 1
	flowPrefix := make([]byte, keyLength)
	flowPrefix[0] = database.FlowPrefix
	binary.PutUvarint(flowPrefix[1:], projectNo)

	issueFlows := make([]*pb_flow.Flow, 0)

	if err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(flowPrefix); it.ValidForPrefix(flowPrefix); it.Next() {
			item := it.Item()

			flowData, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("error copying flow data: %w", err)
			}

			flowObj := &pb_flow.Flow{}
			if err := proto.Unmarshal(flowData, flowObj); err != nil {
				return fmt.Errorf("error unmarshaling flow: %w", err)
			}

			issueFlows = append(issueFlows, flowObj)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return issueFlows, nil
}
