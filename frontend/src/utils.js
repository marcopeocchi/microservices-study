export const getHost = () => {
  const host = import.meta.env.PROD ?
    `${window.location.hostname}:${window.location.port}` :
    `localhost:4456`
  return `${window.location.protocol}//${host}`
}

export const getHostStatic = () => `${getHost()}/static`
export const getHostThumbnails = () => `${getHost()}/thumbs`

export const isOrderedByDate = () => localStorage.getItem("fetch-mode") === "date"
export const isOrderedByName = () => localStorage.getItem("fetch-mode") === "name"

export const isImage = (str) => str.match(/(\.jpg)|(\.jpeg)|(\.png)|(\.webp)|(\.avif)|(\.bmp)|(\.gif)/gmi)