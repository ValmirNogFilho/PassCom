import React, { useState } from 'react'
import "./container.css"
import src from "../assets/srcwhite.svg"
import cart from "../assets/cart.svg"
import Flights from './Flights'
import Tickets from './Tickets'
const Container = () => {
    const [page, setPage] = useState(0)
    return (
        <div className="route">
            <div className="route-content">
                {page === 0 ?
                    <Flights /> :
                    <Tickets />
                }
            </div>
            <menu>
                <div className="flight icon" onClick={() => setPage(0)}
                    style={{backgroundColor: page===0?"#bd1616":"transparent"}}>
                    <img src={src} alt="" />
                </div>
                <div className="cart icon" onClick={() => setPage(1)}
                    style={{backgroundColor: page===1?"#bd1616":"transparent"}}>
                    <img src={cart} alt="" />
                </div>
            </menu>
        </div>

    )
}

export default Container