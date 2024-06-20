package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/paulosman/ticket-service/db"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	database := db.NewDatabase()

	router.HandleFunc("/status", StatusHandler())
	router.HandleFunc("/events", AddEventHandler(database)).Methods("POST")
	router.HandleFunc("/events", GetEventsHandler(database))
	router.HandleFunc("/events/{id:[0-9]+}", GetEventHandler(database))
	router.HandleFunc("/events/{id:[0-9]+}/ticket", CreateTicketHandler(database)).Methods("POST")

	return router
}

func StatusHandler() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func AddEventHandler(database *db.Database) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event db.Event
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&event); err != nil {
			error(w, http.StatusBadRequest, "Error decoding payload: "+err.Error())
			return
		}
		defer r.Body.Close()
		database.AddEvent(&event)
		accepted(w, event)
	}
}

func GetEventsHandler(database *db.Database) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok(w, database.GetEvents())
	}
}

func GetEventHandler(database *db.Database) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			error(w, http.StatusBadRequest, "Invalid id provided: "+err.Error())
			return
		}
		event := database.GetEvent(int32(id))
		ok(w, event)
	}
}

func CreateTicketHandler(database *db.Database) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			error(w, http.StatusBadRequest, "Invalid id provided: "+err.Error())
			return
		}
		decoder := json.NewDecoder(r.Body)
		var ticket db.Ticket
		if err := decoder.Decode(&ticket); err != nil {
			error(w, http.StatusBadRequest, "Error decoding payload: "+err.Error())
			return
		}
		event := database.GetEvent(int32(id))
		if event.CapacityLeft() == 0 {
			error(w, http.StatusForbidden, "Insufficient capacity for event")
			return
		}
		ticket.EventId = int32(id)
		event.Tickets = append(event.Tickets, ticket)
		database.AddEvent(event)

		ok(w, event)
	}
}

func error(w http.ResponseWriter, code int, message string) {
	respond(w, code, map[string]string{"error": message})
}

func accepted(w http.ResponseWriter, body interface{}) {
	respond(w, http.StatusAccepted, body)
}

func ok(w http.ResponseWriter, body interface{}) {
	respond(w, http.StatusOK, body)
}

func respond(w http.ResponseWriter, code int, body interface{}) {
	response, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}
