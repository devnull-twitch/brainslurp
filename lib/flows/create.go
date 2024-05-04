package flows

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/devnull-twitch/brainslurp/lib/database"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type CreateOptions struct {
	ProjectNo    uint64
	Title        string
	Requirements []*pb_flow.FlowRequirement
	Actions      []*pb_flow.FlowActions
}

func Create(db *badger.DB, opts CreateOptions) error {
	flowNo, err := database.NextNumber(db, database.FlowSequenceKey)
	if err != nil {
		return fmt.Errorf("error getting new flow number: %w", err)
	}

	if opts.ProjectNo <= 0 {
		return fmt.Errorf("no project number")
	}

	keyLength := (2 * binary.MaxVarintLen64) + 1
	flowKey := make([]byte, keyLength)
	flowKey[0] = database.FlowPrefix
	binary.PutUvarint(flowKey[1:], opts.ProjectNo)
	binary.PutUvarint(flowKey[binary.MaxVarintLen64+1:], flowNo)

	flowObj := &pb_flow.Flow{
		Number:       flowNo,
		Title:        opts.Title,
		CreatedAt:    time.Now().Unix(),
		Requirements: opts.Requirements,
		Actions:      opts.Actions,
	}
	flowVal, err := proto.Marshal(flowObj)
	if err != nil {
		return fmt.Errorf("unable to marshal issue: %w", err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(flowKey, flowVal)
	}); err != nil {
		return fmt.Errorf("unable to insert issue: %w", err)
	}

	logrus.WithField("key", fmt.Sprintf("%x", flowKey)).Info("flow inserted")

	go issues.UpdateIssues(db, opts.ProjectNo, flowObj)

	return nil
}
