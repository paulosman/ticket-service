package db

type Ticket struct {
	EventId int32  `json:"event_id"`
	Name    string `json:"name"`
}

func NewTicket(eventId int32, name string) *Ticket {
	return &Ticket{
		EventId: eventId,
		Name:    name,
	}
}
