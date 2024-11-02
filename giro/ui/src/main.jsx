import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import Login from './apps/Login.jsx'
import Home from './apps/Home.jsx'
import "./main.css"
export const url = "http://localhost:9999"

const router = createBrowserRouter(
    [
        {
            path: "/",
            element: <Login />
        },
        {
            path: "/home",
            element: <Home />
        }
    ]
)

createRoot(document.getElementById('root')).render(
  
    <RouterProvider router={router} /> 
  
)