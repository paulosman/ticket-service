package db

import "testing"

func TestAddTicketToEventNotFound(t *testing.T) {
	database := NewDatabase()
	event, ok := database.AddTicketToEvent(42, Ticket{Name: "Avery"})
	if event != nil || ok {
		t.Fatalf("expected not found to return nil,false, got %#v,%v", event, ok)
	}
}

func TestAddTicketToEventFull(t *testing.T) {
	database := NewDatabase()
	event := NewEvent("show", 1)
	database.AddEvent(event)

	if _, ok := database.AddTicketToEvent(event.Id, Ticket{Name: "Avery"}); !ok {
		t.Fatal("expected first ticket to be accepted")
	}
	if _, ok := database.AddTicketToEvent(event.Id, Ticket{Name: "Jordan"}); ok {
		t.Fatal("expected second ticket to be rejected when at capacity")
	}
}

func TestDeleteEventRemovesKey(t *testing.T) {
	database := NewDatabase()
	event := NewEvent("show", 1)
	database.AddEvent(event)
	database.DeleteEvent(event.Id)
	if got := database.GetEvent(event.Id); got != nil {
		t.Fatalf("expected event to be deleted, got %#v", got)
	}
}
