import React, { useState } from 'react'
import "./selectboxes.css"

const SelectBoxes = ({ airports }) => {
    const [srcValue, setSrcValue] = useState("Origem")
    const [destValue, setDestValue] = useState("Origem")
    const [isAnimating, setIsAnimating] = useState(false);

    const swap = () => {
        const temp = srcValue;
        setSrcValue(destValue)
        setDestValue(temp)
        setIsAnimating(true);
        setTimeout(() => setIsAnimating(false), 900);
    }

    return (
        <div className="select-boxes">
            <select
                className="route-input source-input"
                value={srcValue}
                onChange={(e) => setSrcValue(e.target.value)}
            >
                <option value="" disabled>
                    Origem
                </option>
                {airports.map((airport, i) => (
                    <option key={i} value={airport.City.CityName}>
                        {airport.City.CityName}
                    </option>
                ))}
            </select>

            <button onClick={swap}
            className={isAnimating ? 'spin-animation' : ''}>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
                </svg>
            </button>

            <select
                className="route-input dest-input"
                value={destValue}
                onChange={(e) => setDestValue(e.target.value)}
            >
                <option value="" disabled>
                    Destino
                </option>
                {airports.map((airport, i) => (
                    <option key={i} value={airport.City.CityName}>
                        {airport.City.CityName}
                    </option>
                ))}
            </select>


        </div>
    )
}

export default SelectBoxes