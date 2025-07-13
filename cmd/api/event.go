package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rizky-ardiansah/event-api/internal/database"
)

// CreateEvent creates a new event
//
//	@Summary		Create a new event
//	@Description	Create a new event with authentication required
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		database.Event	true	"Event data"
//	@Success		201		{object}	database.Event
//	@Failure		400		{object}	gin.H
//	@Failure		401		{object}	gin.H
//	@Failure		500		{object}	gin.H
//	@Security		BearerAuth
//	@Router			/api/v1/events [post]
func (app *application) createEvents(c *gin.Context) {
	// Handler logic for creating events
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)
	event.OwnerId = user.Id

	err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEvents returns all events
//
//	@Summary		Returns all events
//	@Description	Returns all events
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]database.Event
//	@Router			/api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

// GetEvent returns a single event by ID
//
//	@Summary		Get event by ID
//	@Description	Get a single event by its ID
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	database.Event
//	@Failure		400	{object}	gin.H
//	@Failure		404	{object}	gin.H
//	@Failure		500	{object}	gin.H
//	@Router			/api/v1/events/{id} [get]
func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent updates an existing event
//
//	@Summary		Update an event
//	@Description	Update an existing event (only owner can update)
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Event ID"
//	@Param			event	body		database.Event	true	"Updated event data"
//	@Success		200		{object}	database.Event
//	@Failure		400		{object}	gin.H
//	@Failure		401		{object}	gin.H
//	@Failure		403		{object}	gin.H
//	@Failure		404		{object}	gin.H
//	@Failure		500		{object}	gin.H
//	@Security		BearerAuth
//	@Router			/api/v1/events/{id} [put]
func (app *application) updateEvents(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}

	user := app.GetUserFromContext(c)

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this event"})
		return
	}

	updatedEvent := &database.Event{}

	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = id

	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent deletes an event
//
//	@Summary		Delete an event
//	@Description	Delete an event (only owner can delete)
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		204	{object}	gin.H
//	@Failure		400	{object}	gin.H
//	@Failure		401	{object}	gin.H
//	@Failure		403	{object}	gin.H
//	@Failure		404	{object}	gin.H
//	@Failure		500	{object}	gin.H
//	@Security		BearerAuth
//	@Router			/api/v1/events/{id} [delete]
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this event"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"success": "Delete success"})
}

// AddAttendeeToEvent adds an attendee to an event
//
//	@Summary		Add attendee to event
//	@Description	Add a user as attendee to an event (only event owner can add)
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Event ID"
//	@Param			userId	path		int	true	"User ID"
//	@Success		201		{object}	database.Attendee
//	@Failure		400		{object}	gin.H
//	@Failure		401		{object}	gin.H
//	@Failure		403		{object}	gin.H
//	@Failure		404		{object}	gin.H
//	@Failure		409		{object}	gin.H
//	@Failure		500		{object}	gin.H
//	@Security		BearerAuth
//	@Router			/api/v1/events/{id}/attendees/{userId} [post]
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user := app.GetUserFromContext(c)

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to add attendees to this event"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "attendee already exist"})
		return
	}

	attendee := database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add attendee"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// GetAttendeesForEvent gets all attendees for an event
//
//	@Summary		Get attendees for event
//	@Description	Get all users who are attending a specific event
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	[]database.User
//	@Failure		400	{object}	gin.H
//	@Failure		500	{object}	gin.H
//	@Router			/api/v1/events/{id}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to to retreive attendees for events"})
		return
	}

	c.JSON(http.StatusOK, users)

}

// DeleteAttendeeFromEvent removes an attendee from an event
//
//	@Summary		Remove attendee from event
//	@Description	Remove a user from an event's attendee list (only event owner can remove)
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Event ID"
//	@Param			userId	path		int	true	"User ID"
//	@Success		204		{object}	gin.H
//	@Failure		400		{object}	gin.H
//	@Failure		401		{object}	gin.H
//	@Failure		403		{object}	gin.H
//	@Failure		404		{object}	gin.H
//	@Failure		500		{object}	gin.H
//	@Security		BearerAuth
//	@Router			/api/v1/events/{id}/attendees/{userId} [delete]
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	user := app.GetUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete attendees from this event"})
		return
	}

	err = app.models.Attendees.Delete(userId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee from event"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"success": "Attendee deleted successfully"})
}

// GetEventsByAttendee gets all events that a user is attending
//
//	@Summary		Get events by attendee
//	@Description	Get all events that a specific user is attending
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Attendee/User ID"
//	@Success		200	{object}	[]database.Event
//	@Failure		400	{object}	gin.H
//	@Failure		500	{object}	gin.H
//	@Router			/api/v1/attendees/{id}/events [get]
func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee id"})
		return
	}

	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events for attendee"})
		return
	}

	c.JSON(http.StatusOK, events)
}
