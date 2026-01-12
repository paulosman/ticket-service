package db

import "sync"

type Database struct {
	Events      map[int32]*Event
	mutex       sync.Mutex
	nextEventId int32
}

func NewDatabase() *Database {
	return &Database{
		Events:      make(map[int32]*Event),
		nextEventId: 1,
	}
}

func (d *Database) AddEvent(event *Event) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if event.Id == 0 {
		event.Id = d.nextEventId
		d.nextEventId++
	}
	d.Events[event.Id] = event
	return true
}

func (d *Database) GetEvents() []*Event {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	events := make([]*Event, 0)
	for _, event := range d.Events {
		if event != nil {
			events = append(events, event)
		}
	}
	return events
}

func (d *Database) GetEvent(id int32) *Event {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.Events[id]
}

func (d *Database) DeleteEvent(id int32) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.Events, id)
}

func (d *Database) AddTicketToEvent(id int32, ticket Ticket) (*Event, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	event := d.Events[id]
	if event == nil {
		return nil, false
	}
	if event.CapacityLeft() == 0 {
		return event, false
	}
	ticket.EventId = id
	event.Tickets = append(event.Tickets, ticket)
	return event, true
}
