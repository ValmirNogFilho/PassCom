import "./flights.css"
import * as utils from "../utils/utils"
const Flights = ({ flights, addToCart, setCartItemCount}) => {

  const handleAdd = (ID) => {
    addToCart(ID)
    setCartItemCount((prev) => prev + 1)
  }

  return (
    <div className='flights'>
      {flights.map(f => {
        const imgUrl = utils.findCompany(f.Company)
        return (
          <div className="flight" key={f.ID}>
            <div className="row">
              <b className="flight-route">
                {f.OriginAirport.City.Name} =&gt; {f.DestinationAirport.City.Name}
              </b>
              <img src={imgUrl} className='company-brand' width={"50px"} />
            </div>
            <div className="span seats">Passagens: {f.Seats}</div>
            <div className="row">
              <span className="price">Valor: <b>R${f.Price},00</b></span>
              <button className='buy' onClick={() => handleAdd(f.ID)}>ADICIONAR</button>
            </div>
          </div>
        )
      }
      )}
    </div>
  )
}

export default Flights