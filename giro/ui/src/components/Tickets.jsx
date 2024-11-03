import React, { useEffect, useState } from 'react'
import "./tickets.css"
import { apiService } from '../axios'
import * as utils from "../utils/utils"

const Tickets = () => {
  const [tickets, setTickets] = useState([])
  useEffect(() => {
    const fetchTickets = async () => {
      try {
        const res = await apiService.getTickets();
        setTickets(res.data.Data.Tickets)
      } catch (error) {
        console.error(error)
      }
    }
    fetchTickets()
  }, [])

  const cancelBuy = async (ID) => {
    try {
      const res = await apiService.cancelTicket({TicketId: ID});
      setTickets(tickets.filter(w => w.ID !== ID))
    } catch (error) {
      console.error(error)
    }
  }


  return (
    <div className='tickets'>{
      (tickets.map((f) =>
        {
          const imgUrl = utils.findCompany(f.Company)
          return (
              <div className="flight" key={f.ID}>
                  <div className="row">
                      <b className="flight-route">
                          {f.Src.Name} =&gt; {f.Dest.Name}
                      </b>
                      <img src={imgUrl} className='company-brand' width={"50px"} />
                  </div>
                  <div className="row">
                      <button className='cancel' onClick={() => cancelBuy(f.ID)}>REMOVER</button>
                  </div>
              </div>
          )
      }
      ))
    }</div>
  )
}

export default Tickets