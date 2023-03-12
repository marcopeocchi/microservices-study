import { Route, Routes } from '@solidjs/router';
import { Component } from 'solid-js';
import Home from './Home';
import { lazy } from "solid-js";
import Login from './Login';

const Gallery = lazy(() => import('./Gallery'))
const Help = lazy(() => import('./Help'))

const App: Component = () => {
  return (
    <div class="bg-neutral-900 text-neutral-100 min-h-screen">
      <Routes>
        <Route path="/" component={Home} />
        <Route path="/login" component={Login} />
        <Route path="/gallery/:id" component={Gallery} />
        <Route path="/help" component={Help} />
      </Routes>
    </div>
  )
}

export default App;
