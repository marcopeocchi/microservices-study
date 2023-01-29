import { getHost } from "../utils"

export const useLogin = () => (password) => {
  return new Promise((resolve) => {
    fetch(`${getHost()}/user`, {
      method: 'POST',
      body: JSON.stringify({
        password: password
      })
    })
      .then(res => resolve(res.status === 200))
  })
}