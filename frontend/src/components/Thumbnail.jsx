import { Link } from "react-router-dom";
import { getHostThumbnails } from "../utils";

export default function Thumbnail({ entry, isHidden }) {
  return (
    <Link to={`/gallery/${entry.name}`}>
      <div className='mb-3 hover:text-pink-400 duration-75 cursor-pointer text-center'>
        {!isHidden ?
          <img
            alt=''
            className='rounded hover:border-2 hover:border-pink-400 duration-75'
            loading='lazy'
            src={`${getHostThumbnails()}/${entry.thumbnail}`}
            style={{
              objectFit: "cover",
              width: "285px",
              height: "400px",
              cursor: "pointer"
            }} /> :
          <div style={{
            objectFit: "cover",
            height: "400px",
            cursor: "pointer"
          }}
            className="bg-neutral-800 rounded hover:border-2 duration-75">
          </div>
        }
        <div className='mt-2 fw-semibold'>
          {isHidden ? '...' : entry.name}
        </div>
      </div>
    </Link>
  )
}