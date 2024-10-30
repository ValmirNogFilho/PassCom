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
    const token = "a06e8023-9cad-45db-991f-9b18d490c8b2"
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
  getFlights: async ({ FlightIds }) => await client.post("flights", {FlightIds})
}

const user = { Username: "pedrocosta", Password: "senhaSegura789" };

try {
  // let res = await apiService.login(user)
  // console.log(res.data)
  let res = await apiService.getFlights({FlightIds:[1, 5, 7]})
  console.log(res.data.Data.Flights[0])

} catch (error) {
  console.error(error)
}
// let res = await apiService.getUser()
// console.log(res.data)
// let res = await apiService.logout()
// console.log(res.data)