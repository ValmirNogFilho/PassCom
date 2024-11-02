import React, { useEffect, useState } from 'react'
import { apiService } from '../axios'
import "./flights.css"
import giro_brand from "../assets/brand_giro.svg"
import rumos_brand from "../assets/brand_rumos.svg"
import boreal_brand from "../assets/brand_boreal.svg"
const Flights = ({flights}) => {
  

  const findCompany = (company) => {
    switch (company) {
      case "giro":
        return giro_brand
      case "boreal":
        return boreal_brand
      default:
        return rumos_brand
    }
  }


  return (
    <div className='flights'>
      {flights.map(f => {
        const imgUrl = findCompany(f.Company)
        return (
          <div className="flight" key={f.ID}>
            <div className="row">
              <b className="flight-route">
                {f.OriginAirport.City.CityName} =&gt; {f.DestinationAirport.City.CityName}
              </b>
              <img src={imgUrl} className='company-brand' width={"50px"} />
            </div>
            <div className="span seats">Passagens: {f.Seats}</div>
            <div className="row">
              <span className="price">Valor: <b>R${f.Price},00</b></span>
              <button className='buy'>COMPRAR</button>
            </div>
          </div>
        )
      }
      )}
    </div>
  )
}

export default Flights