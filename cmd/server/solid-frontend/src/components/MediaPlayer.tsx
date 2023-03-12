import clsx from "clsx"
import { Component, createEffect, createSignal, Show } from "solid-js"
import { isImage } from "../utils/files"

type Props = {
  onLoad: VoidFunction
  onEnded?: VoidFunction
  onClick?: VoidFunction
  fullscreen: boolean
  loading?: boolean
  src: string
}

const MediaPlayer: Component<Props> = (props) => {
  const [isVideo, setIsVideo] = createSignal(true)

  createEffect(() => {
    if (props.src) {
      isImage(props.src) ? setIsVideo(false) : setIsVideo(true)
      return
    }
    setIsVideo(false)
  })

  return (
    <div>
      <Show when={isVideo()}>
        <video
          autoplay
          muted
          controls
          src={props.src}
          class={clsx(props.fullscreen ? "" : "h-screen", props.loading ? "opacity-50" : "block")}
          style={{ "user-select": "none" }}
          onClick={props.onClick}
          onPlay={props.onLoad}
          onEnded={props.onEnded}
        />
      </Show>
      <Show when={!isVideo()}>
        <img
          src={props.src}
          class={clsx(props.fullscreen ? "" : "h-screen", props.loading ? "opacity-50" : "block")}
          style={{ "user-select": "none" }}
          onClick={props.onClick}
          onLoad={props.onLoad}
        />
      </Show>
    </div>
  )
}

export default MediaPlayer