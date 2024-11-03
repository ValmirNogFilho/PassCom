import React, { useEffect, useState } from 'react'
import { apiService } from '../axios'
import * as utils from "../utils/utils"
import "./flights.css"
const Cart = () => {
    const [wishes, setWishes] = useState([])

    const removeFromWishlist = async (ID, key) => {
        try {
            const res = await apiService.removeFromWishlist({FlightId:ID})
            setWishes(wishes.filter(w => w.key !== key))
        } catch (error) {
            console.error(error)
        }
    }

    const buy = async (ID) => {
        try {
            apiService.buyTicket({FlightId:ID})
            apiService.removeFromWishlist({FlightId:ID})
            setWishes(wishes.filter(w => w.ID !== ID))
        } catch (error) {
            console.error(error)
        }
    }

    useEffect(
        () => {
            const getCart = async () => {
                try {
                    const fetchedWishes = await apiService.getWishlist()
                    setWishes(fetchedWishes.data.Data.Wishes)
                } catch (error) {
                    console.error(error)
                }
            }

            getCart()
        }, []
    )

    return (
        <div>{
            wishes.map((f, i) => {
                const imgUrl = utils.findCompany(f.Company)
                f.key = i
                return (
                    <div className="flight" key={i}>
                        <div className="row">
                            <b className="flight-route">
                                {f.OriginAirport.City.CityName} =&gt; {f.DestinationAirport.City.CityName}
                            </b>
                            <img src={imgUrl} className='company-brand' width={"50px"} />
                        </div>
                        <div className="span seats">Passagens: {f.Seats}</div>
                        <div className="row">
                            <span className="price">Valor: <b>R${f.Price},00</b></span>
                            <button className='cancel' onClick={() => removeFromWishlist(f.ID, i)}>REMOVER</button>
                            <button className='buy' onClick={() => buy(f.ID, i)}>COMPRAR</button>
                        </div>
                    </div>
                )
            })
        }</div>
    )
}

export default Cart