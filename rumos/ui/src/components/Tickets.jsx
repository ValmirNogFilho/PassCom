import React, { useEffect, useState } from 'react'
import "./tickets.css"
import { apiService } from '../axios'
const Tickets = () => {
  const [tickets, setTickets] = useState([])
  useEffect(() => {
    const fetchTickets = async () => {
      try {
        const res = await apiService.getTickets();
        setTickets(res.data.Data.Tickets)
        console.log(tickets)
      } catch (error) {
        console.error(error)
      }
    }
    fetchTickets()
  }, [])

  return (
    <div className='tickets'>{
      (tickets.map((ticket) =>
        <div className="ticket">
          {ticket.Src.CityName}-{ticket.Src.State}/{ticket.Dest.CityName}-{ticket.Dest.State}
        </div>
      ))
    }</div>
  )
}

export default Tickets