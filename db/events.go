package db

type Event struct {
	Id       int32    `json:"id"`
	Name     string   `json:"name"`
	Capacity int      `json:"capacity"`
	Tickets  []Ticket `json:"tickets"`
}

func NewEvent(name string, capacity int) *Event {
	return &Event{
		Name:     name,
		Capacity: capacity,
		Tickets:  make([]Ticket, 0),
	}
}

func (e *Event) CapacityLeft() int {
	return len(e.Tickets) - e.Capacity
}
