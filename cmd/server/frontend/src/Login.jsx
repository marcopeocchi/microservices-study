import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Logo from './components/Logo'
import { useLogin } from './hooks/useLogin'

export default function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const navigate = useNavigate()
  const login = useLogin()

  const performLogin = async () => {
    await login(username, password)
    navigate('/')
  }

  const detectEnterKey = (event) => {
    if (event.key === 'Enter') {
      performLogin()
    }
  }

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Logo hideSubText />
      <div className='mt-8 flex flex-col items-center'>
        <input
          placeholder='username'
          name="fuu-user"
          className="bg-neutral-800 rounded h-10 w-64 text-center"
          type="text"
          onChange={(e) => setUsername(e.target.value)}
          onKeyDown={detectEnterKey}
        />
        <input
          placeholder='password'
          name="fuu-pass"
          className="mt-2 bg-neutral-800 rounded h-10 w-64 text-center "
          type="password"
          onChange={(e) => setPassword(e.target.value)}
          onKeyDown={detectEnterKey}
        />
        <button className="mt-2 h-10 w-full font-semibold bg-blue-400 p-2 rounded" onClick={performLogin}>
          Login
        </button>
      </div>
    </div>
  )
}