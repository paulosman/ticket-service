package db

type Database struct {
	Events map[int32]*Event

	nextEventId int32
}

func NewDatabase() *Database {
	return &Database{
		Events:      make(map[int32]*Event),
		nextEventId: 1,
	}
}

func (d *Database) AddEvent(event *Event) bool {
	if event.Id == 0 {
		event.Id = d.nextEventId
		d.nextEventId++
	}
	d.Events[event.Id] = event
	return true
}

func (d *Database) GetEvents() []*Event {
	events := make([]*Event, 0)
	for _, event := range d.Events {
		events = append(events, event)
	}
	return events
}

func (d *Database) GetEvent(id int32) *Event {
	return d.Events[id]
}
