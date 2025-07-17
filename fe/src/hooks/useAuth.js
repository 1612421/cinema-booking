import { useState, useEffect, useCallback } from 'react'

export const useAuth = () => {
  const [user, setUser] = useState(null)

  useEffect(() => {
    const stored = localStorage.getItem('user')

    if (stored) {
      setUser(JSON.parse(stored))
    }
  }, [])

  const login = useCallback((userData, token) => {
    setUser(userData)
    localStorage.setItem('user', JSON.stringify(user))
    localStorage.setItem('access_token', token)
  }, [])

  const logout = useCallback(() => {
    setUser(null)
    localStorage.removeItem('user')
    localStorage.removeItem('access_token')
  }, [])

  return { user, login, logout }
}