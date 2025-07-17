import axios from './axios'

export const login = (data) => axios.post('/users/login', data)
export const register = (data) => axios.post('/users/register', data)