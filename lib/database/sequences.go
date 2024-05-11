package database

import (
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v4"
)

const (
	IssueSequenceKey   byte = 0x1
	ViewSequenceKey    byte = 0x2
	ProjectSequenceKey byte = 0x3
	UserSequenceKey    byte = 0x4
	FlowSequenceKey    byte = 0x5
	TagSequenceKey     byte = 0x6
)

var sequenceGenLock = sync.Mutex{}
var sequenceInstances = make(map[byte]*badger.Sequence)

func NextNumber(db *badger.DB, sequenceKey byte) (uint64, error) {
	sequenceGenLock.Lock()

	sequence, ok := sequenceInstances[sequenceKey]
	if !ok {
		var err error
		sequence, err = db.GetSequence([]byte{SequencePrefix, sequenceKey}, 10)
		if err != nil {
			sequenceGenLock.Unlock()
			return 0, fmt.Errorf("unable to create issue sequence: %w", err)
		}

		sequenceInstances[sequenceKey] = sequence
	}
	sequenceGenLock.Unlock()

	nextIssue, err := sequence.Next()
	if err != nil {
		return 0, fmt.Errorf("unabel to get next issue number from sequence: %w", err)
	}

	// we dontr really like 0 as a number because that we want to know real zero values
	if nextIssue <= 0 {
		nextIssue, err = sequence.Next()
		if err != nil {
			return 0, fmt.Errorf("unabel to get next issue number from sequence: %w", err)
		}
	}

	return nextIssue, nil
}
