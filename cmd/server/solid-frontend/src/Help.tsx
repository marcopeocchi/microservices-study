import { useNavigate } from "@solidjs/router";
import { Component } from "solid-js";
import Button from "./components/Button";

const Help: Component = () => {
  const navigate = useNavigate()

  return (
    <div class="container mx-auto p-8">
      <div class="text-5xl font-extrabold text-blue-400">
        Help
      </div>
      <Button class="w-16 my-6" onClick={() => navigate('/')}>Back</Button>
      <div class="mt-6">
        Fuu is a simple image viewer.
      </div>
      <div>
        It is oriended for manga/comic display but you can easily handle every kind of image (including videos!).
      </div>
      <div class="text-3xl font-bold mt-8 mb-4">
        Gallery mode shortcuts
      </div>
      <div>
        <code class="bg-neutral-800 p-1 rounded text-blue-400">V</code> &rarr; Vertical Split / Split view
      </div>
      <div>
        <code class="bg-neutral-800 p-1 rounded text-blue-400">S</code> &rarr; Span horizontally
      </div>
      <div>
        <code class="bg-neutral-800 p-1 rounded text-blue-400">P</code> &rarr; Start/Stop slideshow
      </div>
      <div class="mt-2">
        <code class="bg-neutral-800 p-1 rounded text-blue-400">Backspace</code> &rarr; Go back
      </div>
    </div>
  )
}

export default Help