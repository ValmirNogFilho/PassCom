import React, { useEffect, useState } from "react";
import {
    MapContainer,
    TileLayer,
    Marker,
    Popup,
    Polyline,
} from "react-leaflet";
import "leaflet/dist/leaflet.css";
import { apiService } from "../axios";


const BrazilMap = ({capitals, flights}) => {
    const [selectedCity, setSelectedCity] = useState("");
    const position = [-15.7801, -47.9292] 

    function selectCity(city) {
        setSelectedCity(city);
    }

    const switchLineColor = (company) => {
        switch (company) {
            case "rumos":
                return "red"
            case "giro":
                return "green"
            default:
                return "blue"
        }
    }

    return (
        <>
            <MapContainer
                minZoom={3.5}
                center={position}
                zoom={3.5}
                style={{ height: "100%", width: "100%", borderRadius:'20px', overflow:'hidden' }}
            >
                <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                />
                {capitals.map((capital) => {
                    console.log(capital.City)
                    return (
                        <Marker
                            key={capital.name}
                            position={[capital.City.Latitude, capital.City.Longitude]}
                            eventHandlers={{
                                click: () => selectCity(capital.City.Name),
                            }}
                        >
                            <Popup key={capital.Name}>
                                {capital.Name}, {capital.City.State}
                            </Popup>
                        </Marker>
                    )
                })}
                {flights && (
                    <>
                        {flights.map((line, i) => {
                            const lineColor = switchLineColor(line.Company)
                            return (
                                <Polyline
                                    key={i}
                                    pathOptions={{color: lineColor}}
                                    positions={[
                                        [
                                            line.OriginAirport.City.Latitude,
                                            line.OriginAirport.City.Longitude,
                                        ],
                                        [
                                            line.DestinationAirport.City.Latitude,
                                            line.DestinationAirport.City.Longitude,
                                        ],
                                    ]}
                                    color="blue"
                                />
                            )
                        }
                        )}
                    </>
                )}
            </MapContainer>
        </>
    );
};

export default BrazilMap;