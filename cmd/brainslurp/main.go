package main

import (
	"github.com/devnull-twitch/brainslurp/internal/server"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	"github.com/devnull-twitch/brainslurp/lib/project"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	"github.com/devnull-twitch/brainslurp/lib/user"
	"github.com/devnull-twitch/brainslurp/lib/view"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		logrus.WithError(err).Error("unable to open BadgerDB database")
		return
	}
	defer db.Close()

	userNo, err := user.Create(db, user.CreateOptions{
		Username: "devnull@brainslurp",
		Password: "testing",
	})
	if err != nil {
		panic(err)
	}

	projectNo, err := project.Create(db, project.CreateOptions{
		Title:         "Brainslurp Dev",
		Public:        false,
		CreatorUserNo: userNo,
	})
	if err != nil {
		panic(err)
	}

	err = user.AddProject(db, userNo, projectNo)
	if err != nil {
		panic(err)
	}

	viewNo, err := view.Create(db, view.CreateOptions{
		Title:     "Test List",
		ProjectNo: projectNo,
	})
	if err != nil {
		panic(err)
	}

	if err := issues.Create(db, issues.CreateOptions{
		ProjectNo:  projectNo,
		Title:      "Test no1",
		ViewStatus: []uint64{viewNo},
		Category:   pb_issue.IssueCategory_Bug,
	}); err != nil {
		panic(err)
	}
	if err := issues.Create(db, issues.CreateOptions{
		ProjectNo:  projectNo,
		Title:      "Test no2",
		ViewStatus: []uint64{},
		Category:   pb_issue.IssueCategory_Feature,
		Tags: []*pb_issue.Tag{
			{Label: "ToDo", ColorCode: "#FF0000"},
		},
	}); err != nil {
		panic(err)
	}
	if err := issues.Create(db, issues.CreateOptions{
		ProjectNo:  projectNo,
		Title:      "Test no3",
		ViewStatus: []uint64{viewNo},
		Category:   pb_issue.IssueCategory_Operations,
	}); err != nil {
		panic(err)
	}
	if err := issues.Create(db, issues.CreateOptions{
		ProjectNo:  projectNo,
		Title:      "Test no4",
		ViewStatus: []uint64{},
		Category:   pb_issue.IssueCategory_Question,
	}); err != nil {
		panic(err)
	}

	flows.Create(db, flows.CreateOptions{
		ProjectNo: projectNo,
		Title:     "Complete Feature",
		Requirements: []*pb_flow.FlowRequirement{
			{
				InCategory: pb_issue.IssueCategory_Feature,
				HasTags:    []*pb_issue.Tag{{Label: "ToDo"}},
			},
		},
		Actions: []*pb_flow.FlowActions{
			{
				Title:      "Finish",
				RemoveTags: []*pb_issue.Tag{{Label: "ToDo"}},
				AddTags:    []*pb_issue.Tag{{Label: "Finished", ColorCode: "#0F0"}},
			},
		},
	})

	server.Run(db)
}
