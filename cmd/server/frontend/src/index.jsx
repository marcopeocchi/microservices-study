import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import {
  createBrowserRouter,
  RouterProvider
} from "react-router-dom"

import './index.css'
import { routes } from './routes'

const root = ReactDOM.createRoot(document.getElementById('root'))
const router = createBrowserRouter(routes)

root.render(
  <StrictMode>
    <main className='bg-neutral-900 text-neutral-100 min-h-screen'>
      <RouterProvider router={router} />
    </main>
  </StrictMode>
);
