import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import App from './App'
import CssBaseline from '@mui/material/CssBaseline'
import { AuthProvider } from './contexts/AuthContext'

ReactDOM.createRoot(document.getElementById('root')).render(
    // <React.StrictMode>
        <BrowserRouter>
            <AuthProvider>
                <CssBaseline />
                <App />
            </AuthProvider>
        </BrowserRouter>
    // </React.StrictMode>
)