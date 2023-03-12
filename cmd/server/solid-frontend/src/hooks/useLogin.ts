import { getHost } from "../utils/url"

export const useLogin = () => async (username: string, password: string) => {
  const res = await fetch(`${getHost()}/user/login`, {
    method: 'POST',
    body: JSON.stringify({
      username: username,
      password: password
    })
  })
  const token = await res.text()
  return token
}