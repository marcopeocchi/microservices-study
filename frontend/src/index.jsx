import React, { Suspense } from 'react';
import ReactDOM from 'react-dom/client';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";

import './index.css';
import App from './App';
import Login from './Login';
import Spinner from './components/Spinner';

const Gallery = React.lazy(() => import('./Gallery'))
const Help = React.lazy(() => import('./Help'))

const router = createBrowserRouter([
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
]);

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <main className='bg-neutral-900 text-neutral-100 min-h-screen'>
      <RouterProvider router={router} />
    </main>
  </React.StrictMode>
);
