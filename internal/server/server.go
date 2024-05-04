package server

import (
	"net/http"

	"github.com/dgraph-io/badger/v4"
)

func Run(db *badger.DB) error {
	fileServer := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.HandleFunc("/{$}", LoginHandler(db))

	http.HandleFunc("/projects", HandleProjectListing(db))

	http.HandleFunc("/project/{projectNo}", HandleIssueList(db))
	http.HandleFunc("/project/{projectNo}/issues", HandleIssueList(db))
	http.HandleFunc("/project/{projectNo}/issues/view/{viewNo}", HandleIssueList(db))

	http.HandleFunc("/project/{projectNo}/issue/{issueNo}/tags", HandleNewIssueTag(db))
	http.HandleFunc("DELETE /project/{projectNo}/issue/{issueNo}/tag/{tagName}", HandleIssueTag(db))

	return http.ListenAndServe(":3000", nil)
}
