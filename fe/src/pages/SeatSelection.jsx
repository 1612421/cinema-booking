import React, {useEffect, useRef, useState} from 'react'
import {useParams, useNavigate, useLocation} from 'react-router-dom'
import {fetchSeatsByShowtime} from '../api/showtime'
import {holdSeat, releaseSeat} from '../api/booking'
import Button from '@mui/material/Button'
import Grid from '@mui/material/Grid'
import Typography from '@mui/material/Typography'
import Box from '@mui/material/Box'
import {useAuth} from '../contexts/AuthContext'
import {toast} from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import {connectSocket} from "../utils/socket";

const SeatSelection = () => {
    const {showtimeId} = useParams()
    const navigate = useNavigate()
    const {user, logout, isAuthLoaded} = useAuth()
    const location = useLocation()

    const [seats, setSeats] = useState([])
    const [selected, setSelected] = useState([])
    const [hoveredSeat, setHovered] = useState(null)        // seatId đang hover
    const [interest, setInterest] = useState({})
    const socketRef = useRef(null)
    const selectedRef = useRef(null)

    /* ───────────────── Redirect nếu chưa login ───────────────── */
    useEffect(() => {
        if (!isAuthLoaded) return
        if (!user) navigate('/login', {replace: true})
    }, [user, navigate, isAuthLoaded])


    useEffect(() => {
        const restoredSeats = location.state?.seats
        if (restoredSeats?.length > 0) {
            setSelected(restoredSeats)
        }
    }, [location.state])

    /* ───────────────── Lấy danh sách ghế ───────────────── */
    useEffect(() => {
        fetchSeatsByShowtime(showtimeId)
            .then(res => {
                const allSeats = res.data?.data || []
                setSeats(allSeats)

                // Nếu có state (từ trang Checkout quay lại)
                if (location.state?.seats?.length) {
                    const previouslySelected = location.state.seats

                    // Lọc bỏ những ghế đã bị đặt (is_available = false)
                    const validSelected = previouslySelected.filter(sel =>
                        allSeats.find(s => s.id === sel.id && s.is_available)
                    )

                    setSelected(validSelected)
                }
            })
            .catch(console.error)
    }, [showtimeId])

    /* ───────────────── Toggle chọn ghế ───────────────── */
    const toggleSeat = (seat) => {
        const picked = selected.find(s => s.id === seat.id)
        if (picked) {
            setSelected(prev => prev.filter(s => s.id !== seat.id))
            releaseSeat({seat_id: seat.id, showtime_id: showtimeId}).catch(err => {
                if (err.status === 401) {
                    logout();
                    navigate('/login');
                } else {
                    console.log(err)
                }
            })
        } else {
            setSelected(prev => [...prev, seat])
            holdSeat({seat_id: seat.id, showtime_id: showtimeId}).then(r => r.data).then(data => {
                const qty = data?.data?.qty - 1
                if (qty > 0) {
                    toast.info(`🔥 Hurry up! There are ${data?.data?.qty - 1} people interested in seat ${seat.row}${seat.number}`, {
                        icon: false,
                        position: 'top-right',
                        autoClose: 4000,
                        closeOnClick: true,
                        draggable: false,
                        style: {
                            background: '#fff3cd',
                            color: '#856404',
                            border: '1px solid #ffeeba',
                            fontWeight: 500,
                            fontSize: '16px',
                            textAlign: 'center',
                        },
                    })
                }
            }).catch(err => {
                if (err.status === 401) {
                    logout();
                    navigate('/login');
                } else {
                    console.log(err)
                }
            })
        }
    }

    /* ───────────────── Hover: tăng + lấy interest ───────────────── */
    const handleMouseEnter = (seat) => {
        setHovered(seat.id)

        // Step 2 – lấy interest count hiện tại cho seat
        fetch(`/api/v1/showtimes/${showtimeId}/seats/${seat.id}/interest`)
            .then(r => r.json())
            .then(data => {
                const isPicked = selected.find(s => s.id === seat.id)
                const qty = isPicked ? data?.data?.qty - 1 : data?.data?.qty
                setInterest(prev => ({
                    ...prev,
                    [seat.id]: qty
                }))
            })
            .catch(console.error)
    }

    const getSeatStyle = (seat) => {
        if (!seat.is_available) return {
            variant: 'outlined',
            color: 'inherit',
            disabled: true,
            sx: {borderColor: '#ccc', color: '#aaa'}
        }

        const isSelected = selected.find(s => s.id === seat.id)
        return {
            variant: isSelected ? 'contained' : 'outlined',
            color: isSelected ? 'primary' : 'info',
            disabled: false,
            sx: {},
        }
    }

    useEffect(() => {
        selectedRef.current = selected
    }, [selected])

    useEffect(() => {
        if (!showtimeId) return

        if (socketRef.current && (socketRef.current.readyState === WebSocket.OPEN || socketRef.current.readyState === WebSocket.CONNECTING)) {
            console.log('🔁 Socket already connected')
            return
        }

        const token = localStorage.getItem('access_token') // hoặc từ context
        socketRef.current = connectSocket(token, (msg) => {
            try {
                console.log(msg)

                if (msg.event === 'booking_created') {
                    const {seat_ids} = msg.data
                    const conflictSeats = selectedRef.current.filter(s => seat_ids.includes(s.id))
                    // 👉 Cập nhật trạng thái ghế bị book
                    setSeats(prev =>
                        prev.map(seat =>
                            seat_ids.includes(seat.id) ? {...seat, is_available: false} : seat
                        )
                    )

                    // 🔥 Hiển thị cảnh báo nếu user đang chọn 1 trong các ghế đó
                    if (conflictSeats.length > 0) {
                        toast.warn(
                            ` There ${conflictSeats.length > 1 ? `are ${conflictSeats.length} seats` : `is ${conflictSeats.length} seat`} you selected has just been booked by someone else. Please choose again.`,
                            {autoClose: 5000}
                        )
                        setSelected(prev => prev.filter(s => !seat_ids.includes(s.id)))
                    } else {
                        toast.info(
                            `A few seats have been booked by someone else`,
                            {autoClose: 5000}
                        )
                    }
                }
            } catch (err) {
                console.error('Invalid socket message:', err)
            }
        })

        // 🧹 Clean up khi component unmount hoặc showtimeId thay đổi
        return () => {
            if (socketRef.current) {
                socketRef.current.close()
                socketRef.current = null
            }
        }
    }, [showtimeId])

    return (
        <Box>
            <Typography variant="h5" gutterBottom>Select seats</Typography>

            <Grid container spacing={1}>
                {seats.map(seat => {
                    const isSel = selected.some(s => s.id === seat.id)
                    return (
                        <Grid item key={seat.id} sx={{textAlign: 'center', position: 'relative'}}>
                            {/* Nút ghế */}
                            <Button
                                {...getSeatStyle(seat)}
                                onClick={() => toggleSeat(seat)}
                                onMouseEnter={() => handleMouseEnter(seat)}
                                onMouseLeave={() => setHovered(null)}
                                size="small"
                            >
                                {seat.row}{seat.number}
                            </Button>

                            {/* Pop‑up interest chỉ khi hover */}
                            {(hoveredSeat === seat.id) && interest[seat.id] > 0 && (
                                <Box
                                    sx={{
                                        position: 'absolute',
                                        top: '-26px',
                                        left: '50%',
                                        transform: 'translateX(-50%)',
                                        bgcolor: 'white',
                                        color: 'red',
                                        border: '1px solid #f0b3b3',
                                        px: 1,
                                        py: 0.3,
                                        borderRadius: '12px',
                                        fontSize: '12px',
                                        whiteSpace: 'nowrap',
                                        zIndex: 10,
                                        fontWeight: 500,
                                    }}
                                >
                                    🔥 There are {interest[seat.id]} people interested in!
                                </Box>
                            )}
                        </Grid>
                    )
                })}
            </Grid>

            <Button
                sx={{mt: 2}}
                variant="contained"
                disabled={selected.length === 0}
                onClick={() => navigate('/checkout', {state: {seats: selected, showtimeId}})}
            >
                Continue ({selected.length})
            </Button>
        </Box>
    )
}

export default SeatSelection
