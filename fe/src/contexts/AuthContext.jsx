import { createContext, useContext, useState, useEffect, useCallback } from 'react'

const AuthContext = createContext(null)

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null)
    const [isAuthLoaded, setIsAuthLoaded] = useState(false)

    useEffect(() => {
        try {
            const stored = localStorage.getItem('user')
            if (stored && stored !== 'undefined') {
                setUser(JSON.parse(stored))
            }
        } catch (e) {
            localStorage.removeItem('user')
        } finally {
            setIsAuthLoaded(true)
        }
    }, [])

    const login = useCallback((userData, token) => {
        setUser(userData)
        localStorage.setItem('user', JSON.stringify(userData))
        localStorage.setItem('access_token', token)
    }, [])

    const logout = useCallback(() => {
        setUser(null)
        localStorage.removeItem('user')
        localStorage.removeItem('access_token')
    }, [])

    return (
        <AuthContext.Provider value={{ user, login, logout, isAuthLoaded }}>
            {children}
        </AuthContext.Provider>
    )
}

export const useAuth = () => useContext(AuthContext)