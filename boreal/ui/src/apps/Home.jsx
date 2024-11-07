import React, { useEffect, useState } from 'react'
import BrazilMap from '../components/BrazilMap'
import SelectBoxes from '../components/SelectBoxes'
import Header from '../components/Header'
import { apiService } from '../axios'
import Container from '../components/Container'

const Home = () => {
  const [cartItemCount, setCartItemCount] = useState(0)
  const [airports, setAirports] = useState([])
  const [name, setName] = useState("")
  const [srcValue, setSrcValue] = useState("Origem")
  const [destValue, setDestValue] = useState("Destino")
  const [flights, setFlights] = useState([])

  const addToCart = async (ID) => {
    try {
      const res = await apiService.addToWishlist({ FlightId: ID })
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    const fetchCapitals = async () => {
      try {
        const res = await apiService.getAirports()
        setAirports(res.data.Data.Airports)
      } catch (error) {
        console.error(error)
      }
    }
    const fetchName = async () => {
      try {
        const res = await apiService.getUser()
        const name = `, ${res.data.Data.user.Name}`
        setName(name)
      } catch (error) {
        console.error(error)
      }
    }
    fetchCapitals()
    fetchName()
  }, [])

  useEffect(() => {
    const fetchFlights = async () => {
      try {
        const res = await apiService.getRoute({
          src: srcValue,
          dest: destValue,
        });
        setFlights(res.data.Data.paths.filter(f => f.Seats > 0));
      } catch (error) {
        console.error(error);
      }
    };

    if (srcValue !== "Origem" && destValue !== "Destino") {
      fetchFlights();
    }
  }, [srcValue, destValue])

  return (
    <div className="home">
      <Header />
      <div className="title">
        <h1>A liberdade de explorar começa com a gente{name}.</h1>
        <h3>Qual o seu destino?</h3>
      </div>
      <div className="content">

        <div className="map">
          <BrazilMap capitals={airports} flights={flights} />
        </div>
        <div className="search">
          <SelectBoxes airports={airports} srcValue={srcValue}
            destValue={destValue} setSrcValue={setSrcValue} setDestValue={setDestValue} />
          <Container flights={flights} addToCart={addToCart} 
          cartItemCount={cartItemCount} setCartItemCount={setCartItemCount} />
        </div>
      </div>
    </div>
  )
}

export default Home