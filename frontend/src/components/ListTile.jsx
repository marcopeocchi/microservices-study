import { Link } from "react-router-dom";

export default function ListTile({ entry }) {
  return (
    <Link to={`/gallery/${entry.name}`}>
      <div className='h-14 rounded flex items-center justify-between p-3 bg-neutral-800 hover:text-blue-400 hover:bg-neutral-700 duration-75 cursor-pointer text-center'>
        <div>{entry.name}</div>
        <div className="text-neutral-400">
          <span>{new Date(entry.lastModified).toLocaleDateString()}</span>
          <span>{' - '}</span>
          <span>{new Date(entry.lastModified).toLocaleTimeString()}</span>
        </div>
      </div>
    </Link>
  )
}