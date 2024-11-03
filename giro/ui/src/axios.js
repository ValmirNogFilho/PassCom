import axios from "axios";

const host = "localhost";
const port = 9999;

const client = axios.create({
  baseURL: `http://${host}:${port}/`,

  headers: {
    'Content-Type': 'application/json',
  },
})

client.interceptors.request.use(
  (config) => {
    if (config.url.includes('login')) {
      return config;
    }
    const token = sessionStorage.getItem("token")
    if (token) {
      config.headers.Authorization = `${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export const apiService = {
  login: async ({ Username, Password }) => await client.post("login", { Username, Password }),
  logout: async () => await client.get("logout"),
  getUser: async () => await client.get("user"),
  getFlights: async ({ FlightIds }) => await client.post("flights", {FlightIds}),
  buyTicket: async ({ FlightId }) => await client.post("ticket", {FlightId}),
  cancelTicket: async({ TicketId }) => await client.delete("ticket", {params:{id:TicketId}}),
  getTickets: async () => await client.get("tickets"),
  getAirports: async () => await client.get("airports"),
  getRoute: async ({src, dest}) => await client.get("route", {params: {src, dest}}),
  getWishlist: async () => await client.get("wishlist"),
  addToWishlist: async ({FlightId}) => await client.post("wishlist", {FlightId}),
  removeFromWishlist: async ({FlightId}) => await client.delete("wishlist", {params:{id:FlightId}})
}
