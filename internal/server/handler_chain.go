package server

import (
	"net/http"

	"github.com/dgraph-io/badger/v4"
)

type (
	nextCall     func()
	chainHandler func(
		db *badger.DB,
		w http.ResponseWriter,
		r *http.Request,
		next nextCall,
	)
)

func walkChain(db *badger.DB, w http.ResponseWriter, r *http.Request, stack ...chainHandler) {
	for _, currentHandler := range stack {
		end := true
		currentHandler(db, w, r, func() {
			end = false
		})
		if end {
			return
		}
	}
}
