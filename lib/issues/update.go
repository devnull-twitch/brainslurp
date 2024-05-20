package issues

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/devnull-twitch/brainslurp/lib/database"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func Update(db *badger.DB, projectNo uint64, issue *pb_issue.Issue) (*pb_issue.Issue, []*pb_flow.Flow, error) {

	issueKeyLength := (2 * binary.MaxVarintLen64) + 1
	issueKey := make([]byte, issueKeyLength)
	issueKey[0] = database.IssuePrefix
	binary.PutUvarint(issueKey[1:], projectNo)
	binary.PutUvarint(issueKey[binary.MaxVarintLen64+1:], issue.GetNumber())

	var modIssue *pb_issue.Issue
	var issueFlows []*pb_flow.Flow
	if err := db.Update(func(txn *badger.Txn) error {
		modIssue = proto.Clone(issue).(*pb_issue.Issue)
		var err error
		issueFlows, err = loadAndCheckFlows(txn, projectNo, modIssue)
		if err != nil {
			return err
		}

		issueVal, err := proto.Marshal(modIssue)
		if err != nil {
			return fmt.Errorf("unable to marshal issue: %w", err)
		}

		return txn.Set(issueKey, issueVal)
	}); err != nil {
		return nil, nil, fmt.Errorf("unable to update issue: %w", err)
	}

	logrus.WithField("key", fmt.Sprintf("%x", issueKey)).Info("issue updated")

	return modIssue, issueFlows, nil
}

func loadAndCheckFlows(txn *badger.Txn, projectNo uint64, issueObj *pb_issue.Issue) ([]*pb_flow.Flow, error) {
	keyLength := binary.MaxVarintLen64 + 1
	flowPrefix := make([]byte, keyLength)
	flowPrefix[0] = database.FlowPrefix
	binary.PutUvarint(flowPrefix[1:], projectNo)

	issueFlows := make([]*pb_flow.Flow, 0)

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	for it.Seek(flowPrefix); it.ValidForPrefix(flowPrefix); it.Next() {
		item := it.Item()

		flowData, err := item.ValueCopy(nil)
		if err != nil {
			return nil, fmt.Errorf("error copying flow data: %w", err)
		}

		flowObj := &pb_flow.Flow{}
		if err := proto.Unmarshal(flowData, flowObj); err != nil {
			return nil, fmt.Errorf("error unmarshaling flow: %w", err)
		}

		currentState := hasFlow(issueObj, flowObj.GetNumber())
		reqOk := doesIssueMatchRequirement(issueObj, flowObj.GetRequirements())

		if !reqOk && currentState {
			issueObj.Flows = removeFlow(issueObj, flowObj.GetNumber())
			logrus.WithField("flow", flowObj.GetNumber()).Info("removed flow from issue")
		}
		if reqOk && !currentState {
			issueObj.Flows = append(issueObj.Flows, &pb_issue.FlowStatus{Number: flowObj.GetNumber(), SetAt: time.Now().Unix()})
			logrus.WithField("flow", flowObj.GetNumber()).Info("added flow to issue")
		}
		if reqOk {
			issueFlows = append(issueFlows, flowObj)
		}
	}

	return issueFlows, nil
}
