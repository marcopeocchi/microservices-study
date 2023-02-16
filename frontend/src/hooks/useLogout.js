import { getHost } from "../utils"

export const useLogout = () => async () => {
  await fetch(`${getHost()}/user/logout`)
}