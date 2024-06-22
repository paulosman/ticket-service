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
	if event.Id == 0 {
		event.Id = d.nextEventId
		d.nextEventId++
	}
	d.Events[event.Id] = event
	d.mutex.Unlock()
	return true
}

func (d *Database) GetEvents() []*Event {
	d.mutex.Lock()
	events := make([]*Event, 0)
	for _, event := range d.Events {
		events = append(events, event)
	}
	d.mutex.Unlock()
	return events
}

func (d *Database) GetEvent(id int32) *Event {
	var event *Event
	d.mutex.Lock()
	event = d.Events[id]
	d.mutex.Unlock()
	return event
}

func (d *Database) DeleteEvent(id int32) {
	d.Events[id] = nil
}
