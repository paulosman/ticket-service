package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/paulosman/ticket-service/db"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	database := db.NewDatabase()

	router.HandleFunc("/status", StatusHandler())
	router.HandleFunc("/events", AddEventHandler(database)).Methods("POST")
	router.HandleFunc("/events", GetEventsHandler(database))
	router.HandleFunc("/events/{id:[0-9]+}", GetEventHandler(database))
	router.HandleFunc("/events/{id:[0-9]+}", DeleteEventHandler(database)).Methods("DELETE")
	router.HandleFunc("/events/{id:[0-9]+}/ticket", CreateTicketHandler(database)).Methods("POST")

	router.Use(loggingMiddleware)
	return router
}

func StatusHandler() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func DeleteEventHandler(database *db.Database) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			error(w, http.StatusBadRequest, r.Method, "/events/:id", "Invalid id provided: "+err.Error())
			return
		}
		database.DeleteEvent(int32(id))
		ok(w, r.Method, "/events/:id", "Deleted")
	}
}

func AddEventHandler(database *db.Database) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event db.Event
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&event); err != nil {
			error(w, http.StatusBadRequest, r.Method, r.URL.Path, "Error decoding payload: "+err.Error())
			return
		}
		defer r.Body.Close()
		database.AddEvent(&event)
		accepted(w, r.Method, r.URL.Path, event)
	}
}

func GetEventsHandler(database *db.Database) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok(w, r.Method, r.URL.Path, database.GetEvents())
	}
}

func GetEventHandler(database *db.Database) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			error(w, http.StatusBadRequest, r.Method, "/events/:id", "Invalid id provided: "+err.Error())
			return
		}
		event := database.GetEvent(int32(id))
		if event == nil {
			error(w, http.StatusNotFound, r.Method, "/events/:id", "no event found")
			return
		}
		ok(w, r.Method, "/events/:id", event)
	}
}

func CreateTicketHandler(database *db.Database) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := "/events/:id/ticket"
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			error(w, http.StatusBadRequest, r.Method, path, "Invalid id provided: "+err.Error())
			return
		}
		decoder := json.NewDecoder(r.Body)
		var ticket db.Ticket
		if err := decoder.Decode(&ticket); err != nil {
			error(w, http.StatusBadRequest, r.Method, path, "Error decoding payload: "+err.Error())
			return
		}
		event, added := database.AddTicketToEvent(int32(id), ticket)
		if event == nil {
			error(w, http.StatusNotFound, r.Method, path, "no event found")
			return
		}
		if !added {
			error(w, http.StatusForbidden, r.Method, path, "Insufficient capacity for event")
			return
		}

		ok(w, r.Method, path, event)
	}
}

func error(w http.ResponseWriter, code int, method string, path string, message string) {
	respond(w, code, method, path, map[string]string{"error": message})
}

func accepted(w http.ResponseWriter, method string, path string, body interface{}) {
	respond(w, http.StatusAccepted, method, path, body)
}

func ok(w http.ResponseWriter, method string, path string, body interface{}) {
	respond(w, http.StatusOK, method, path, body)
}

func respond(w http.ResponseWriter, code int, method string, path string, body interface{}) {
	response, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
