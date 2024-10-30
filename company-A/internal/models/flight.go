package models

import (
	"gorm.io/gorm"
)

type Flight struct {
	gorm.Model
	OriginAirportID      uint    `gorm:"not null"`
	OriginAirport        Airport `gorm:"foreignKey:OriginAirportID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	DestinationAirportID uint    `gorm:"not null"`
	DestinationAirport   Airport `gorm:"foreignKey:DestinationAirportID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Seats                int
}

// func (f *Flight) AcceptReservation() (*Ticket, error) {

// 	if f.Seats > 0 {
// 		f.Seats--
// 		ticket := new(Ticket)
// 		ticket.Id = uuid.New()
// 		ticket.FlightId = f.Id
// 		f.Passengers = append(f.Passengers, ticket)
// 		return ticket, nil
// 	}
// 	return nil, errors.New("no seats available")
// }

// func (f *Flight) ProcessReservations() {
// 	for session := range f.Queue {
// 		session.Mu.Lock()
// 		ticket, err := f.AcceptReservation()
// 		if err != nil {
// 			fmt.Printf("Session %s: error reserving for flight %s - %s\n", session.ID, f.Id, err)
// 			session.FailedReservations <- f.Id.String()
// 		} else {
// 			session.FailedReservations <- "success"
// 			fmt.Println("session" + session.ID.String())
// 			ticket.ClientId = session.ClientID
// 			id := uuid.New()
// 			session.Reservations[id] = Reservation{
// 				Id:        id,
// 				CreatedAt: time.Now(),
// 				Ticket:    ticket,
// 			}
// 			fmt.Printf("Session %s: flight %s reserved successfully!\n", session.ID, f.Id)
// 		}
// 		fmt.Println(session.FailedReservations)
// 		session.Mu.Unlock()
// 	}
// }
