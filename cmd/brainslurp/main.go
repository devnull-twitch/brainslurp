package main

import (
	"time"

	"github.com/devnull-twitch/brainslurp/internal/server"
	"github.com/devnull-twitch/brainslurp/lib/flows"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	"github.com/devnull-twitch/brainslurp/lib/project"
	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	pb_tag "github.com/devnull-twitch/brainslurp/lib/proto/tag"
	"github.com/devnull-twitch/brainslurp/lib/tag"
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

	testTagNo, err := tag.Create(db, projectNo, &pb_tag.Tag{
		Title:    "Test",
		HexColor: "#F000F0",
	})
	if err != nil {
		panic(err)
	}

	replaceTagNo, err := tag.Create(db, projectNo, &pb_tag.Tag{
		Title:    "Done",
		HexColor: "#00FFF0",
	})
	if err != nil {
		panic(err)
	}

	_, err = tag.Create(db, projectNo, &pb_tag.Tag{
		Title:    "Critical",
		HexColor: "#FC1616",
	})
	if err != nil {
		panic(err)
	}

	_, err = tag.Create(db, projectNo, &pb_tag.Tag{
		Title:    "Blocked",
		HexColor: "#AD0104",
	})
	if err != nil {
		panic(err)
	}

	_, err = tag.Create(db, projectNo, &pb_tag.Tag{
		Title:    "Feedback",
		HexColor: "#FC1689",
	})
	if err != nil {
		panic(err)
	}

	if err := issues.Create(db, projectNo, &pb_issue.Issue{
		Title:      "Test no1",
		Category:   pb_issue.IssueCategory_Bug,
		TagNumbers: []uint64{testTagNo},
	}); err != nil {
		panic(err)
	}
	if err := issues.Create(db, projectNo, &pb_issue.Issue{
		Title:      "Test no2",
		Category:   pb_issue.IssueCategory_Feature,
		TagNumbers: []uint64{testTagNo},
	}); err != nil {
		panic(err)
	}
	if err := issues.Create(db, projectNo, &pb_issue.Issue{
		Title:      "Test no3",
		Category:   pb_issue.IssueCategory_Feature,
		TagNumbers: []uint64{},
		Views: []*pb_issue.ViewStatus{
			{Number: viewNo, SetAt: time.Now().Unix()},
		},
	}); err != nil {
		panic(err)
	}
	if err := issues.Create(db, projectNo, &pb_issue.Issue{
		Title:      "Test no4",
		Category:   pb_issue.IssueCategory_Question,
		TagNumbers: []uint64{replaceTagNo},
	}); err != nil {
		panic(err)
	}

	flows.Create(db, flows.CreateOptions{
		ProjectNo: projectNo,
		Title:     "Setup Test",
		Requirements: []*pb_flow.FlowRequirement{
			{
				InCategory:     pb_issue.IssueCategory_Feature,
				RequiredTagIds: []uint64{testTagNo},
			},
		},
		Actions: []*pb_flow.FlowActions{
			{
				Title:        "Tescht",
				RemoveTagIds: []uint64{testTagNo},
				AddTagIds:    []uint64{replaceTagNo},
			},
		},
	})

	server.Run(db)
}
