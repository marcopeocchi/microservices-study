import { getHost } from "../utils/url"

export const useLogout = () => async () => {
  await fetch(`${getHost()}/user/logout`)
}