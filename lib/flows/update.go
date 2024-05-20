package flows

import (
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

func Update(db *badger.DB, projectNo uint64, flowObj *pb_flow.Flow) error {
	flowKey := database.Keygen(database.FlowPrefix, projectNo, flowObj.GetNumber())
	if err := db.Update(func(txn *badger.Txn) error {
		flowBytes, err := proto.Marshal(flowObj)
		if err != nil {
			return fmt.Errorf("error marshaling flow: %w", err)
		}

		return txn.Set(flowKey, flowBytes)
	}); err != nil {
		return fmt.Errorf("error updating flow: %w", err)
	}

	go issues.UpdateIssues(db, projectNo, flowObj)

	return nil
}
