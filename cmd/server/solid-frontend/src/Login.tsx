import { useNavigate } from '@solidjs/router'
import { Component, createSignal } from 'solid-js'
import Logo from './components/Logo'
import { useLogin } from './hooks/useLogin'

const Login: Component = () => {
  const [username, setUsername] = createSignal('')
  const [password, setPassword] = createSignal('')

  const navigate = useNavigate()
  const login = useLogin()

  const performLogin = async () => {
    await login(username(), password())
    navigate('/')
  }

  const detectEnterKey = (event) => {
    if (event.key === 'Enter') {
      performLogin()
    }
  }

  return (
    <div class="flex flex-col items-center justify-center h-screen">
      <Logo hideSubText />
      <div class='mt-8 flex flex-col items-center'>
        <input
          placeholder='username'
          name="fuu-user"
          class="bg-neutral-800 rounded h-10 w-64 text-center"
          type="text"
          onChange={(e) => setUsername(e.currentTarget.value)}
          onKeyDown={detectEnterKey}
        />
        <input
          placeholder='password'
          name="fuu-pass"
          class="mt-2 bg-neutral-800 rounded h-10 w-64 text-center "
          type="password"
          onChange={(e) => setPassword(e.currentTarget.value)}
          onKeyDown={detectEnterKey}
        />
        <button class="mt-2 h-10 w-full font-semibold bg-blue-400 p-2 rounded" onClick={performLogin}>
          Login
        </button>
      </div>
    </div>
  )
}

export default Login