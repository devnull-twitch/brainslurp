package flows

import (
	"fmt"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	"github.com/dgraph-io/badger/v4"
)

func Get(db *badger.DB, projectNo uint64, flowNo uint64) (*pb_flow.Flow, error) {
	flowKey := database.Keygen(database.FlowPrefix, projectNo, flowNo)
	flowObj := &pb_flow.Flow{}

	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(flowKey)
		if err != nil {
			return err
		}

		if err := database.Read(item, flowObj); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error loading flow: %w", err)
	}

	return flowObj, nil
}
