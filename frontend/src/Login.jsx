import { useNavigate } from 'react-router-dom'
import { useLogin } from './hooks/useLogin'
import Logo from './components/Logo'
import { useState } from 'react'

export default function Login() {
  const [password, setPassword] = useState('')
  const navigate = useNavigate()
  const login = useLogin()

  const performLogin = async () => {
    await login('admin', password)
    navigate('/')
  }

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Logo />
      <div className="mt-6 flex justify-center items-center">
        <input className="bg-neutral-800 rounded-l text-xl h-10 w-64 text-center" type="password" onChange={
          (e) => setPassword(e.target.value)
        } />
        <button className="h-10 bg-blue-400 p-2 rounded-r" onClick={performLogin}>Login</button>
      </div>
    </div>
  )
}