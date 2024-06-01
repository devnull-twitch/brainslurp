package server

import (
	"context"
	"net/http"
	"time"

	"github.com/devnull-twitch/brainslurp/internal/server/components/pages"
	"github.com/devnull-twitch/brainslurp/internal/server/components/shared"
	"github.com/devnull-twitch/brainslurp/lib/issues"
	"github.com/devnull-twitch/brainslurp/lib/proto/issue"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

func HandleNewIssueComment(db *badger.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		walkChain(
			db, w, r,
			authUserWithProjectNo,
			checkIssueNumber,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request, next nextCall) {
				projectObj := getProjectFromContext(ctx)
				issueObj := getIssueFromContext(ctx)
				userObj := getUserFromContext(ctx)

				r.ParseForm()
				comment := r.Form.Get("comment")

				newCommentActivity := &issue.IssueActivity{
					CreatedAt:   time.Now().Unix(),
					TriggerUser: userObj.GetNumber(),
					Body:        comment,
				}
				issueObj.Activities = append(issueObj.Activities, newCommentActivity)

				if _, _, err := issues.Update(db, projectObj.GetNumber(), issueObj); err != nil {
					logrus.WithError(err).Warn("error saving issue with new comment")
					w.WriteHeader(http.StatusInternalServerError)
					pages.Error("Error saving issue with comment").Render(r.Context(), w)
				}

				shared.CommentBox(newCommentActivity, map[uint64]*pb_user.User{userObj.GetNumber(): userObj}).Render(r.Context(), w)
			},
		)
	}
}
