import { getHost } from "../utils"

export const useLogin = () => async (username, password) => {
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