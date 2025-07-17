import React, { useState } from 'react'
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import { login as loginApi } from '../api/auth'
import { useNavigate, Link } from 'react-router-dom'
import {useAuth} from "../contexts/AuthContext";

const Login = () => {
  const navigate = useNavigate()
  const { login } = useAuth()
  const [form, setForm] = useState({ email: '', password: '' })

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      const res = await loginApi(form)
      login(res.data.data.user, res.data.data.access_token)
      navigate('/')
    } catch (err) {
      console.error(err)
    }
  }

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ maxWidth: 400, mx: 'auto' }}>
      <Typography variant="h5" gutterBottom>Login</Typography>
      <TextField fullWidth name="username" label="Username" margin="normal" value={form.username} onChange={handleChange} />
      <TextField fullWidth name="password" label="Password" type="password" margin="normal" value={form.password} onChange={handleChange} />
      <Button fullWidth variant="contained" type="submit" sx={{ mt: 2 }}>Login</Button>
      <Typography variant="body2" sx={{ mt: 2 }}>Don't have an account? <Link to="/register">Register</Link></Typography>
    </Box>
  )
}

export default Login