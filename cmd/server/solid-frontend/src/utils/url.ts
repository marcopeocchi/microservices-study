export const getHost = () => {
  const host = import.meta.env.PROD ?
    `${window.location.hostname}:${window.location.port}` :
    `localhost:4456`
  return `${window.location.protocol}//${host}`
}

export const getHostStatic = () => `${getHost()}/static`
export const getHostOverlay = () => `${getHost()}/overlay`
export const getHostThumbnails = () => `${getHost()}/thumbs`

export const composeQuery = (fetchMode: string, filter: string, page: number, pageSize: number) =>
  `list?fetchBy=${fetchMode}&filter=${filter}&page=${page}&pageSize=${pageSize}`