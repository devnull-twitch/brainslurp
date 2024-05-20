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

func fetchFlows(txn *badger.Txn, projectNo uint64, issueObj *pb_issue.Issue) ([]*pb_flow.Flow, error) {
	listFlows := make([]*pb_flow.Flow, 0)
	for _, flowState := range issueObj.GetFlows() {
		keyLength := (2 * binary.MaxVarintLen64) + 1
		flowKey := make([]byte, keyLength)
		flowKey[0] = database.FlowPrefix
		binary.PutUvarint(flowKey[1:], projectNo)
		binary.PutUvarint(flowKey[binary.MaxVarintLen64+1:], flowState.GetNumber())

		flowItem, err := txn.Get(flowKey)
		if err != nil {
			return nil, fmt.Errorf("error getting flow entry: %w", err)
		}

		flowData, err := flowItem.ValueCopy(nil)
		if err != nil {
			return nil, fmt.Errorf("error copying flow data: %w", err)
		}

		flowObj := &pb_flow.Flow{}
		if err := proto.Unmarshal(flowData, flowObj); err != nil {
			return nil, fmt.Errorf("error unmarshaling flow: %w", err)
		}

		listFlows = append(listFlows, flowObj)
	}

	return listFlows, nil
}
