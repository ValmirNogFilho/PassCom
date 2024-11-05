import React, { useEffect, useState } from 'react';
import { apiService } from '../axios';
import * as utils from "../utils/utils";
import "./flights.css";

const Cart = ({setCartItemCount}) => {
    const [wishes, setWishes] = useState([]);

    const removeFromWishlist = async (FlightId, index) => {
        try {
            await apiService.removeFromWishlist({ FlightId });
            setWishes((prevWishes) => prevWishes.filter((wish, i) => i !== index));
            setCartItemCount((prev) => prev - 1)
        } catch (error) {
            console.error(error);
        }
    };

    const buy = async (FlightId, index) => {
        try {
            // Realiza a compra no servidor
            await apiService.buyTicket({ FlightId });

            // Remove o item do servidor
            await apiService.removeFromWishlist({ FlightId });

            // Remove apenas a primeira ocorrência no estado local (usando o índice passado)
            setWishes((prevWishes) => {
                const newWishes = [...prevWishes];
                newWishes.splice(index, 1); // Remove o item pelo índice
                return newWishes;
            });
            setCartItemCount((prev) => prev - 1)
        } catch (error) {
            console.error(error);
        }
    };

    useEffect(() => {
        const getCart = async () => {
            try {
                const fetchedWishes = await apiService.getWishlist();
                setWishes(fetchedWishes.data.Data.Wishes);
            } catch (error) {
                console.error(error);
            }
        };

        getCart();
    }, []);

    return (
        <div>
            {wishes.map((f, i) => {
                const imgUrl = utils.findCompany(f.Company);
                return (
                    <div className="flight" key={`${f.ID}-${i}`}>
                        <div className="row">
                            <b className="flight-route">
                                {f.OriginAirport.City.Name} =&gt; {f.DestinationAirport.City.Name}
                            </b>
                            <img src={imgUrl} className='company-brand' width={"50px"} alt="company logo" />
                        </div>
                        <div className="span seats">Passagens: {f.Seats}</div>
                        <div className="row">
                            <span className="price">Valor: <b>R${f.Price},00</b></span>
                            <button className='cancel' onClick={() => removeFromWishlist(f.ID, i)}>REMOVER</button>
                            <button className='buy' onClick={() => buy(f.ID, i)}>COMPRAR</button>
                        </div>
                    </div>
                );
            })}
        </div>
    );
};

export default Cart;
