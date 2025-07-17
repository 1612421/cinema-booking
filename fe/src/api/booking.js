import axios from './axios'

export const holdSeat = (payload) => axios.post('/bookings/hold-seat', payload)
export const releaseSeat = (payload) => axios.post('/bookings/release-seat', payload)
export const createBooking = (payload) => axios.post('/bookings', payload)
export const myTickets = () => axios.get('/bookings/my-tickets')