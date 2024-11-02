import React, { useState } from 'react'
import "./container.css"
import src from "../assets/srcwhite.svg"
import cart from "../assets/cart.svg"
import ticket from "../assets/ticket.svg"
import Flights from './Flights'
import Tickets from './Tickets'
const Container = () => {
    const [page, setPage] = useState(0)
    return (
        <div className="route">
            <div className="route-content">
                {page === 0 ?
                    <Flights /> :
                page === 1 ? 
                <div></div> :
                <Tickets />
                }
            </div>
            <menu>
                <div className="flight icon" onClick={() => setPage(0)}
                    style={{backgroundColor: page===0?"#5ad733":"transparent"}}>
                    <img src={src} alt="" />
                </div>
                <div className="cart icon" onClick={() => setPage(1)}
                    style={{backgroundColor: page===1?"#5ad733":"transparent"}}>
                    <img src={cart} alt="" />
                </div>
                <div className="ticket icon" onClick={() => setPage(2)}
                    style={{backgroundColor: page===2?"#5ad733":"transparent"}}>
                    <img src={ticket} alt="" />                     
                </div>
            </menu>
        </div>

    )
}

export default Container