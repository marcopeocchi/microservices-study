import { lazy, Suspense } from 'react'
import App from './App'
import Spinner from './components/Spinner'
import Login from './Login'

const Gallery = lazy(() => import('./Gallery'))
const Help = lazy(() => import('./Help'))

export const routes = [
  {
    path: "/",
    element: <App />,
  },
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/gallery/:id",
    element: (
      <Suspense fallback={<Spinner />}>
        <Gallery />
      </Suspense>
    ),
  },
  {
    path: "/help",
    element: (
      <Suspense fallback={<Spinner />}>
        <Help />
      </Suspense>
    ),
  },
]