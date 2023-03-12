import { A } from "@solidjs/router";
import { Component } from "solid-js";

type Props = {
  entry: Directory
}

const ListTile: Component<Props> = (props) => {
  return (
    <A href={`/gallery/${props.entry.name}`}>
      <div class='h-14 rounded flex items-center justify-between p-3 bg-neutral-800 hover:text-blue-400 hover:bg-neutral-700 duration-75 cursor-pointer text-center'>
        <div>{props.entry.name}</div>
        <div class="text-neutral-400">
          <span>{new Date(props.entry.lastModified).toLocaleDateString()}</span>
          <span>{' - '}</span>
          <span>{new Date(props.entry.lastModified).toLocaleTimeString()}</span>
        </div>
      </div>
    </A>
  )
}

export default ListTile