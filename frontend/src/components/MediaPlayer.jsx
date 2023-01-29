import clsx from "clsx"
import { useEffect, useState } from "react"
import { isImage } from "../utils"

export default function MediaPlayer({
  src,
  onClick,
  onImageLoad,
  onVideoLoad,
  fullscreen,
  loading
}) {
  const [isVideo, setIsVideo] = useState(true)

  useEffect(() => {
    if (src) {
      isImage(src) ?
        setIsVideo(false) :
        setIsVideo(true)
      return
    }
    setIsVideo(false)
  }, [src])

  return (
    isVideo ?
      <video
        autoPlay
        muted
        controls
        src={src}
        className={clsx(fullscreen ? "" : "h-screen", loading ? "opacity-50" : "block")}
        style={{ userSelect: "none" }}
        onClick={onClick}
        onPlay={onVideoLoad}
      /> :
      <img
        src={src}
        className={clsx(fullscreen ? "" : "h-screen", loading ? "opacity-50" : "block")}
        style={{ userSelect: "none" }}
        onClick={onClick}
        onLoad={onImageLoad}
      />
  )
}