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

func UpdateIssues(db *badger.DB, projectNo uint64, flow *pb_flow.Flow) {
	db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := make([]byte, binary.MaxVarintLen64+1)
		prefix[0] = database.IssuePrefix
		binary.PutUvarint(prefix[1:], projectNo)

		issueObj := &pb_issue.Issue{}
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			issueData, err := item.ValueCopy(nil)
			if err != nil {
				logrus.WithError(err).Error("unable to copy data from issue")
				continue
			}

			if err := proto.Unmarshal(issueData, issueObj); err != nil {
				logrus.WithError(err).Error("unable to unmarshal data from issue")
				continue
			}

			reqOk := doesIssueMatchRequirement(issueObj, flow.GetRequirements())
			currentState := hasFlow(issueObj, flow.GetNumber())
			needsUpdate := false
			if reqOk && !currentState {
				issueObj.Flows = append(issueObj.Flows, &pb_issue.FlowStatus{Number: flow.GetNumber(), SetAt: time.Now().Unix()})
				needsUpdate = true
			}
			if !reqOk && currentState {
				issueObj.Flows = removeFlow(issueObj, flow.GetNumber())
				needsUpdate = true
			}

			if needsUpdate {
				issueVal, err := proto.Marshal(issueObj)
				if err != nil {
					logrus.WithError(err).Error("unable to marshal data from issue")
					continue
				}
				if err := txn.Set(item.Key(), issueVal); err != nil {
					logrus.WithError(err).WithField("key", fmt.Sprintf("%x", item.Key())).Error("error updating issue with flow")
				} else {
					logrus.WithField("key", fmt.Sprintf("%x", item.Key())).Info("issue updated with flow")
				}
			}
		}

		return nil
	})
}

func hasFlow(issueObj *pb_issue.Issue, flowNo uint64) bool {
	for _, flowState := range issueObj.GetFlows() {
		if flowState.GetNumber() == flowNo {
			return true
		}
	}

	return false
}

func removeFlow(issueObj *pb_issue.Issue, flowNo uint64) []*pb_issue.FlowStatus {
	states := make([]*pb_issue.FlowStatus, 0, len(issueObj.GetFlows()))
	for _, current := range issueObj.GetFlows() {
		if current.GetNumber() != flowNo {
			states = append(states, current)
		}
	}

	return states
}

func doesIssueMatchRequirement(issueObj *pb_issue.Issue, requirements []*pb_flow.FlowRequirement) bool {
	for _, req := range requirements {
		// make sure issue has all tags
		if len(req.GetHasTags()) > 0 {
			for _, tagReq := range req.GetHasTags() {
				if !hasTag(issueObj, tagReq) {
					return false
				}
			}
		}

		// make sure tags does not have tags it should not have
		if len(req.GetNotTags()) > 0 {
			for _, noTagReq := range req.GetNotTags() {
				if hasTag(issueObj, noTagReq) {
					return false
				}
			}
		}

		if req.GetInCategory() != pb_issue.IssueCategory_Unknown {
			if req.GetInCategory() != issueObj.GetCategory() {
				return false
			}
		}
	}

	return true
}

func hasTag(issueObj *pb_issue.Issue, tag *pb_issue.Tag) bool {
	for _, issueTag := range issueObj.GetTags() {
		if issueTag.Label == tag.Label {
			return true
		}
	}

	return false
}
