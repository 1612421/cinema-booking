import React, { useState } from 'react'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import { register as registerApi } from '../api/auth'
import { useNavigate, Link } from 'react-router-dom'

const Register = () => {
    const navigate = useNavigate()
    const [form, setForm] = useState({
        username: '',
        password: '',
        phone_number: '',
        address: '',
    })

    const handleChange = (e) => {
        setForm({ ...form, [e.target.name]: e.target.value })
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        try {
            await registerApi(form)
            navigate('/login')
        } catch (err) {
            console.error(err)
        }
    }

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ maxWidth: 400, mx: 'auto' }}>
            <Typography variant="h5" gutterBottom>Register</Typography>

            <TextField
                fullWidth
                name="username"
                label="Username"
                margin="normal"
                value={form.username}
                onChange={handleChange}
                required
            />
            <TextField
                fullWidth
                name="password"
                label="Password"
                type="password"
                margin="normal"
                value={form.password}
                onChange={handleChange}
                required
            />
            <TextField
                fullWidth
                name="phone_number"
                label="Phone Number"
                margin="normal"
                value={form.phone_number}
                onChange={handleChange}
                required
            />
            <TextField
                fullWidth
                name="address"
                label="Address"
                margin="normal"
                value={form.address}
                onChange={handleChange}
                required
            />

            <Button fullWidth variant="contained" type="submit" sx={{ mt: 2 }}>Register</Button>

            <Typography variant="body2" sx={{ mt: 2 }}>
                Already have an account? <Link to="/login">Login</Link>
            </Typography>
        </Box>
    )
}

export default Register
