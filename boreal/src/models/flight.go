package models

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type FlightJSON struct {
	Id              uuid.UUID
	SourceAirportId uuid.UUID
	DestAirportId   uuid.UUID
	Passengers      []*Ticket
	Seats           uint
}

type Flight struct {
	Id              uuid.UUID
	SourceAirportId uuid.UUID
	DestAirportId   uuid.UUID
	Passengers      []*Ticket
	Seats           uint
	Queue           chan *Session // Canal de fila para reservas
	Mu              sync.Mutex
}

// AcceptReservation reserves a seat for a flight and returns the ticket if successful.
// If there are no seats available, it returns an error.
//
// The function locks the Flight's mutex to ensure thread safety while processing the reservation.
// It checks if there are any available seats by comparing the number of seats with the length of the passengers slice.
// If there are seats available, it decrements the number of seats, creates a new Ticket, assigns the Flight's ID to the ticket,
// appends the ticket to the passengers slice, and returns the ticket along with a nil error.
// If there are no seats available, it returns nil and an error indicating that there are no seats available.
func (f *Flight) AcceptReservation() (*Ticket, error) {
	f.Mu.Lock()
	defer f.Mu.Unlock()

	if f.Seats > 0 {
		f.Seats--
		ticket := new(Ticket)
		ticket.Id = uuid.New()
		ticket.FlightId = f.Id
		f.Passengers = append(f.Passengers, ticket)
		return ticket, nil
	}
	return nil, errors.New("no seats available")
}

// ProcessReservations processes reservations for a flight.
// It iterates over a queue of sessions and attempts to reserve a seat for each session.
// If a seat is available, it creates a new ticket, assigns the client ID to the ticket,
// and adds the reservation to the session's reservations map.
// If no seats are available, it logs an error message.
func (f *Flight) ProcessReservations() {
	for session := range f.Queue {
		session.Mu.Lock()
		ticket, err := f.AcceptReservation()
		if err != nil {
			fmt.Printf("Session %s: error reserving for flight %s - %s\n", session.ID, f.Id, err)
			session.FailedReservations <- f.Id.String()
		} else {
			session.FailedReservations <- "success"
			fmt.Println("session" + session.ID.String())
			ticket.ClientId = session.ClientID
			id := uuid.New()
			session.Reservations[id] = Reservation{
				Id:        id,
				CreatedAt: time.Now(),
				Ticket:    ticket,
			}
			fmt.Printf("Session %s: flight %s reserved successfully!\n", session.ID, f.Id)
		}
		fmt.Println(session.FailedReservations)
		session.Mu.Unlock()
	}
}
