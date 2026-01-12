package db

import "testing"

func TestEventCapacityLeft(t *testing.T) {
	event := NewEvent("show", 2)
	if got := event.CapacityLeft(); got != 2 {
		t.Fatalf("expected capacity left 2, got %d", got)
	}

	event.Tickets = append(event.Tickets, Ticket{EventId: 1, Name: "Alex"})
	if got := event.CapacityLeft(); got != 1 {
		t.Fatalf("expected capacity left 1, got %d", got)
	}
}
