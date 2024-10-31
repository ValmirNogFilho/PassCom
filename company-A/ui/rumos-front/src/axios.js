import axios from "axios";

const host = "localhost"
const port = 9999

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
    const token = "31f884ab-a9cf-4f0b-9634-6a37fe4c39ae"
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
  cancelTicket: async({ TicketId }) => await client.delete("ticket", {TicketId}),
  getTickets: async () => await client.get("tickets")
}

const user = { Username: "pedrocosta", Password: "senhaSegura789" };

try {
  // let res = await apiService.login(user)
  // console.log(res.data)
  let res = await apiService.getTickets()
  console.log(res.data.Data.Tickets)
  // let res = await apiService.getUser()
  // console.log(res.data)
  // let res = await apiService.buyTicket({FlightId:5})
  // console.log(res.data)
  // let res = await apiService.getFlights({FlightIds:[1, 5, 7]})
  // console.log(res.data.Data.Flights)

} catch (error) {
  console.error(error)
}