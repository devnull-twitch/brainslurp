package server

import (
	"net/http"

	"github.com/dgraph-io/badger/v4"
)

func Run(db *badger.DB) error {
	fileServer := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.HandleFunc("GET /{$}", LoginHandler(db))

	http.HandleFunc("GET /projects", HandleProjectListing(db))

	http.HandleFunc("GET /project/{projectNo}/issues", HandleIssueList(db))
	http.HandleFunc("GET /project/{projectNo}/issues/view/{viewNo}", HandleIssueList(db))
	http.HandleFunc("/project/{projectNo}/issues/new", HandleIssueCreate(db))

	http.HandleFunc("POST /project/{projectNo}/issue/{issueNo}/tags", HandleNewIssueTag(db))
	http.HandleFunc("DELETE /project/{projectNo}/issue/{issueNo}/tag/{tagNumber}", HandleIssueTag(db))

	http.HandleFunc("GET /project/{projectNo}/flows", HandleFlowList(db))
	http.HandleFunc("/project/{projectNo}/flows/new", HandleFlowCreate(db))

	http.HandleFunc("GET /project/{projectNo}/tags", HandleTagsListing(db))
	http.HandleFunc("/project/{projectNo}/tags/new", HandleTagCreate(db))

	return http.ListenAndServe(":3000", nil)
}
