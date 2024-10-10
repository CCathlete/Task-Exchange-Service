package infrastructure

import (
	"database/sql"
	"fmt"
	"net/http"
)

type resources struct {
	db *sql.DB
}

type HandlersGroup struct { // A container object for resources and methods for handling http routes.
	resources resources
}

func NewHandlersGroup(db *sql.DB) *HandlersGroup {
	return &HandlersGroup{resources: resources{db: db}}
}

func (h *HandlersGroup) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			fmt.Fprintln(w, "Placeholder for fetching tasks logic.")
			fmt.Println("Checking if case can get multi lines without braces.")
		}
	case "POST":
		fmt.Fprintln(w, "Placeholder for creating a new task.")
	default:
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}

func (h *HandlersGroup) AccountingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Placeholder for accouting logic.")
}
