import React, { useState } from 'react'
import "./container.css"
import src from "../assets/srcwhite.svg"
import cart from "../assets/cart.svg"
import ticket from "../assets/ticket.svg"
import Flights from './Flights'
import Tickets from './Tickets'
import Cart from './Cart'
const Container = ({ flights, addToCart, cartItemCount, setCartItemCount }) => {
    const [page, setPage] = useState(0)
    return (
        <div className="route">
            <div className="route-content">
                {page === 0 ?
                    <Flights flights={flights} addToCart={addToCart}
                        setCartItemCount={setCartItemCount} /> :
                    page === 1 ?
                        <Cart setCartItemCount={setCartItemCount} /> :
                        <Tickets />
                }
            </div>
            <menu>
                <div className="flight-icon icon" onClick={() => setPage(0)}
                    style={{ backgroundColor: page === 0 ? "#3675e2" : "transparent" }}>
                    <img src={src} alt="" />
                </div>
                <div className="cart-icon icon" onClick={() => setPage(1)}
                    style={{ backgroundColor: page === 1 ? "#3675e2" : "transparent" }}>
                    <img src={cart} alt="" />
                    {cartItemCount > 0 && (
                        <span className="badge">{cartItemCount}</span>
                    )}
                </div>
                <div className="ticket-icon icon" onClick={() => setPage(2)}
                    style={{ backgroundColor: page === 2 ? "#3675e2" : "transparent" }}>
                    <img src={ticket} alt="" />
                </div>
            </menu>
        </div>

    )
}

export default Container