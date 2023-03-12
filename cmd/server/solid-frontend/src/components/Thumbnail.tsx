import { A } from "@solidjs/router";
import { Component } from "solid-js";
import { getHostThumbnails } from "../utils/url";

type Props = {
  entry: Directory
  isHidden?: boolean
}

const Thumbnail: Component<Props> = (props) => {
  return (
    <A href={`/gallery/${props.entry.name}`}>
      <div class='mb-3 hover:text-blue-400 duration-75 cursor-pointer text-center'>
        {!props.isHidden ?
          <img
            alt=''
            class='rounded hover:border-2 border-blue-400 duration-75'
            loading='lazy'
            src={`${getHostThumbnails()}/${props.entry.thumbnail}`}
            style={{
              "object-fit": "cover",
              width: "285px",
              height: "400px",
              cursor: "pointer"
            }} /> :
          <div style={{
            "object-fit": "cover",
            height: "400px",
            cursor: "pointer"
          }}
            class="bg-neutral-800 rounded hover:border-2 border-blue-400 duration-75">
          </div>
        }
        <div class='mt-2 fw-semibold'>
          {props.isHidden ? '...' : props.entry.name}
        </div>
      </div>
    </A>
  )
}

export default Thumbnail