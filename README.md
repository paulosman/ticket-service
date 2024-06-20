# ticket-service

This is an example service that exposes endpoints for managing events and attendees. It is meant to be a toy used in demos.

### Endpoints

```
GET  /events           # Get all events
POST /events           # Create a new event
GET  /events/1         # Get the event with the id 1
POST /events/1/ticket  # Create a ticket for the event with the id 1
```