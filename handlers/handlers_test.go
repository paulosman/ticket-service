package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/paulosman/ticket-service/db"
)

func TestCreateTicketHandlerNotFound(t *testing.T) {
	database := db.NewDatabase()
	handler := CreateTicketHandler(database)

	request := httptest.NewRequest(http.MethodPost, "/events/1/ticket", bytes.NewBufferString(`{"name":"Alex"}`))
	request = mux.SetURLVars(request, map[string]string{"id": "1"})
	response := httptest.NewRecorder()

	handler(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestCreateTicketHandlerFull(t *testing.T) {
	database := db.NewDatabase()
	event := db.NewEvent("show", 0)
	database.AddEvent(event)
	handler := CreateTicketHandler(database)

	request := httptest.NewRequest(http.MethodPost, "/events/1/ticket", bytes.NewBufferString(`{"name":"Alex"}`))
	request = mux.SetURLVars(request, map[string]string{"id": "1"})
	response := httptest.NewRecorder()

	handler(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, response.Code)
	}
}

func TestCreateTicketHandlerOK(t *testing.T) {
	database := db.NewDatabase()
	event := db.NewEvent("show", 1)
	database.AddEvent(event)
	handler := CreateTicketHandler(database)

	request := httptest.NewRequest(http.MethodPost, "/events/1/ticket", bytes.NewBufferString(`{"name":"Alex"}`))
	request = mux.SetURLVars(request, map[string]string{"id": "1"})
	response := httptest.NewRecorder()

	handler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if got := database.GetEvent(event.Id); got == nil || len(got.Tickets) != 1 {
		t.Fatalf("expected ticket to be created, got %#v", got)
	}
}
