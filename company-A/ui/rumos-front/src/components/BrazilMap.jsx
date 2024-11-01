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


const BrazilMap = ({capitals}) => {
    const [selectedCity, setSelectedCity] = useState("");
    const position = [-15.7801, -47.9292] 

    function selectCity(city) {
        setSelectedCity(city);
    }

    return (
        <>
            <MapContainer
                minZoom={3.5}
                center={position}
                zoom={3.5}
                style={{ height: "100%", width: "100%" }}
            >
                <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                />
                {capitals.map((capital) => (
                    <Marker
                        key={capital.name}
                        position={[capital.City.Latitude, capital.City.Longitude]}
                        eventHandlers={{
                            click: () => selectCity(capital.City.CityName),
                        }}
                    >
                        <Popup key={capital.Name}>
                            {capital.Name}, {capital.City.State}
                        </Popup>
                    </Marker>
                ))}
                {/* {paths && (
                    <>
                        {paths.map((line, i) => (
                            <Polyline
                                key={i}
                                positions={[
                                    [
                                        line.Path[0].Latitude,
                                        line.Path[0].Longitude,
                                    ],
                                    [
                                        line.Path[1].Latitude,
                                        line.Path[1].Longitude,
                                    ],
                                ]}
                                color="blue"
                            />
                        ))}
                    </>
                )} */}
            </MapContainer>
        </>
    );
};

export default BrazilMap;