import React from 'react'
import {useLocation, useNavigate} from 'react-router-dom'
import Typography from '@mui/material/Typography'
import Button from '@mui/material/Button'
import {createBooking} from '../api/booking'
import {toast} from 'react-toastify'
import Box from "@mui/material/Box";

const Checkout = () => {
    const navigate = useNavigate()
    const {state} = useLocation()
    const {seats, showtimeId} = state || {seats: [], showtimeId: null}

    const total = seats.length * 80000

    const handlePay = () => {
        createBooking({showtime_id: showtimeId, seat_ids: seats.map(s => s.id)})
            .then(() => {
                toast.success('🎉 Payment successful!')
                navigate('/')
            })
            .catch(err => {
                const message =
                    err.response?.data?.error?.error || '❌ Payment failed. Please try again.'
                toast.error(message, {
                    position: 'top-right',
                    autoClose: 4000,
                })
            })
    }

    if (!showtimeId) return null

    return (
        <Box sx={{ px: 4, py: 4 }}>
            <Typography variant="h4" gutterBottom>
                Checkout
            </Typography>

            <Typography variant="h6" gutterBottom>
                Seats: {seats.map(s => `${s.row}${s.number}`).join(', ')}
            </Typography>

            <Typography variant="h6" gutterBottom>
                Total: {total.toLocaleString('vi-VN')}đ
            </Typography>

            <Box sx={{ mt: 4, display: 'flex', justifyContent: 'space-between', gap: 2 }}>
                <Button
                    variant="outlined"
                    onClick={() => navigate(`/showtimes/${showtimeId}/seats`, { state: { seats, showtimeId } })}
                >
                    ⬅️ Back to seat selection
                </Button>

                <Button variant="contained" onClick={handlePay} disabled={seats.length === 0}>
                    PAY
                </Button>
            </Box>
        </Box>
    )
}

export default Checkout