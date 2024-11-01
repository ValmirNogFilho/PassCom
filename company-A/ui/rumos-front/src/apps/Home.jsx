import React, { useEffect, useState } from 'react'
import BrazilMap from '../components/BrazilMap'
import SelectBoxes from '../components/SelectBoxes'
import Header from '../components/Header'
import { apiService } from '../axios'
import Container from '../components/Container'
import Footer from '../components/Footer'

const Home = () => {
  const [airports, setAirports] = useState([])

  useEffect(() => {
    const fetchCapitals = async () => {
      try {
        const res = await apiService.getAirports()
        setAirports(res.data.Data.Airports)
      } catch (error) {
        console.error(error)
      }
    }
    fetchCapitals()
  }, [])

  return (
    <div className="home">
      <Header />
      <div className="title">
        <h1>Voar nunca foi tão fácil.</h1>
        <h3>Qual o seu destino?</h3>
      </div>
      <div className="content">

        <div className="map">
          <BrazilMap capitals={airports} />
        </div>
        <div className="search">
          <SelectBoxes airports={airports} />
          <Container />
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default Home