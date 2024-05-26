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
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
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

	userNo, err := user.Create(db, "testing", &pb_user.User{
		Name:        "devnull@brainslurp",
		CreatedAt:   time.Now().Unix(),
		Memberships: []*pb_user.Membership{},
	})
	if err != nil {
		panic(err)
	}

	secondUserNo, err := user.Create(db, "testing", &pb_user.User{
		Name:        "second@brainslurp",
		CreatedAt:   time.Now().Unix(),
		Memberships: []*pb_user.Membership{},
	})
	if err != nil {
		panic(err)
	}

	projectNo, err := project.Create(db, project.CreateOptions{
		Title:  "Brainslurp Dev",
		Public: false,
	})
	if err != nil {
		panic(err)
	}

	err = user.AddProject(db, userNo, projectNo)
	if err != nil {
		panic(err)
	}
	err = user.AddProject(db, secondUserNo, projectNo)
	if err != nil {
		panic(err)
	}
	project.AddUser(db, userNo, projectNo)
	if err != nil {
		panic(err)
	}
	project.AddUser(db, secondUserNo, projectNo)
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
		CreatedBy:  userNo,
		AssignedTo: []uint64{userNo},
		Category:   pb_issue.IssueCategory_Bug,
		TagNumbers: []uint64{testTagNo},
	}); err != nil {
		panic(err)
	}
	if err := issues.Create(db, projectNo, &pb_issue.Issue{
		Title:      "Test no2",
		CreatedBy:  userNo,
		Body:       "We might *want* this feature but we dont *need* it.",
		Category:   pb_issue.IssueCategory_Feature,
		TagNumbers: []uint64{testTagNo},
		AssignedTo: []uint64{secondUserNo},
	}); err != nil {
		panic(err)
	}
	if err := issues.Create(db, projectNo, &pb_issue.Issue{
		Title:      "Test no3",
		CreatedBy:  userNo,
		CreatedAt:  time.Now().Add(-time.Hour * 50).Unix(),
		Category:   pb_issue.IssueCategory_Feature,
		TagNumbers: []uint64{},
		AssignedTo: []uint64{userNo, secondUserNo},
		Body:       "We **really** need this top notch feature.",
		Views: []*pb_issue.ViewStatus{
			{Number: viewNo, SetAt: time.Now().Unix()},
		},
		Activities: []*pb_issue.IssueActivity{
			{
				CreatedAt:   time.Now().Unix(),
				TriggerUser: userNo,
				Body:        "This is a message",
			},
			{
				CreatedAt:   time.Now().Add(-time.Hour * 12).Unix(),
				TriggerUser: userNo,
				Body:        "This is a message",
			},
			{
				CreatedAt:   time.Now().Add(-time.Hour * 48).Unix(),
				TriggerUser: userNo,
				Body:        "This is a message",
			},
		},
	}); err != nil {
		panic(err)
	}
	if err := issues.Create(db, projectNo, &pb_issue.Issue{
		Title:      "Test no4",
		CreatedBy:  userNo,
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
