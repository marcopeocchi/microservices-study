import { Link } from "react-router-dom";

export default function ListTile({ entry }) {
  return (
    <Link to={`/gallery/${entry.name}`}>
      <div className='h-14 rounded flex items-center p-3 bg-neutral-800 hover:text-pink-400 hover:bg-neutral-700 duration-75 cursor-pointer text-center'>
        {entry.name}
      </div>
    </Link>
  )
}