package server

import (
	"context"
	"net/http"

	pb_flow "github.com/devnull-twitch/brainslurp/lib/proto/flow"
	pb_issue "github.com/devnull-twitch/brainslurp/lib/proto/issue"
	pb_project "github.com/devnull-twitch/brainslurp/lib/proto/project"
	pb_user "github.com/devnull-twitch/brainslurp/lib/proto/user"
	"github.com/dgraph-io/badger/v4"
)

type (
	nextCall     func(context.Context)
	chainHandler func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		next nextCall,
	)
)

func walkChain(db *badger.DB, w http.ResponseWriter, r *http.Request, stack ...chainHandler) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextDbKey, db)
	for _, currentHandler := range stack {
		end := true
		currentHandler(ctx, w, r, func(newCtx context.Context) {
			ctx = newCtx
			end = false
		})
		if end {
			return
		}
	}
}

const contextDbKey string = "db"

func getDbFromContext(ctx context.Context) *badger.DB {
	return ctx.Value(contextDbKey).(*badger.DB)
}

const contextProjectKey string = "project"

func setProjectOnContext(ctx context.Context, projectObj *pb_project.Project) context.Context {
	return context.WithValue(ctx, contextProjectKey, projectObj)
}
func getProjectFromContext(ctx context.Context) *pb_project.Project {
	return ctx.Value(contextProjectKey).(*pb_project.Project)
}

const contextUserKey string = "user"

func setUserOnContext(ctx context.Context, projectObj *pb_user.User) context.Context {
	return context.WithValue(ctx, contextUserKey, projectObj)
}
func getUserFromContext(ctx context.Context) *pb_user.User {
	return ctx.Value(contextUserKey).(*pb_user.User)
}

const contextFlowKey string = "flow"

func setFlowOnContext(ctx context.Context, flowObj *pb_flow.Flow) context.Context {
	return context.WithValue(ctx, contextFlowKey, flowObj)
}
func getFlowFromContext(ctx context.Context) *pb_flow.Flow {
	return ctx.Value(contextFlowKey).(*pb_flow.Flow)
}

const contextIssueKey string = "issue"

func setIssueOnContext(ctx context.Context, issueObj *pb_issue.Issue) context.Context {
	return context.WithValue(ctx, contextIssueKey, issueObj)
}
func getIssueFromContext(ctx context.Context) *pb_issue.Issue {
	return ctx.Value(contextIssueKey).(*pb_issue.Issue)
}

const contextIssueFlowsKey string = "issue_flows"

func setIssueFlowsOnContext(ctx context.Context, issueFlows []*pb_flow.Flow) context.Context {
	return context.WithValue(ctx, contextIssueFlowsKey, issueFlows)
}
func getIssueFlowsFromContext(ctx context.Context) []*pb_flow.Flow {
	return ctx.Value(contextIssueFlowsKey).([]*pb_flow.Flow)
}
