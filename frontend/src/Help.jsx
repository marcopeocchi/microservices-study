import { useNavigate } from "react-router-dom";
import { Button } from "./components/Button";

export default function Help() {
  const navigate = useNavigate()

  return (
    <div className="container mx-auto p-8">
      <div className="text-5xl font-extrabold text-blue-400">
        Help
      </div>
      <Button className="w-16 my-6" onClick={() => navigate('/')}>Back</Button>
      <div className="mt-6">
        Fuu is a simple image viewer.
      </div>
      <div>
        It is oriended for manga/comic display but you can easily handle every kind of image (including videos!).
      </div>
      <div className="text-3xl font-bold mt-8 mb-4">
        Gallery mode shortcuts
      </div>
      <div>
        <code className="bg-neutral-800 p-1 rounded text-blue-400">V</code> &rarr; Vertical Split / Split view
      </div>
      <div>
        <code className="bg-neutral-800 p-1 rounded text-blue-400">S</code> &rarr; Span horizontally
      </div>
      <div>
        <code className="bg-neutral-800 p-1 rounded text-blue-400">P</code> &rarr; Start/Stop slideshow
      </div>
      <div className="mt-2">
        <code className="bg-neutral-800 p-1 rounded text-blue-400">Backspace</code> &rarr; Go back
      </div>
    </div>
  )
}