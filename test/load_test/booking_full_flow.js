import http from 'k6/http'
import { sleep, check } from 'k6'

export const options = {
    stages: [
        { duration: '5s', target: 40 }, // tăng từ 0 lên 20 VUs
        { duration: '10s', target: 200 }, // tăng tiếp lên 100
        { duration: '30s', target: 500 }, // spike lên 200 users
        { duration: '15s', target: 500 }, // giữ ở mức cao 200 VUs
        { duration: '10s', target: 100 }, // giảm xuống
        { duration: '5s', target: 0 },   // dừng
    ],
}

const BASE_URL = 'http://localhost:8080' // 🔁 chỉnh lại nếu khác

const USERNAME = 'user1'
const PASSWORD = 'test123'

const randItem = (arr) => arr[Math.floor(Math.random() * arr.length)]

export default function () {
    // Step 1: Login
    const loginRes = http.post(`${BASE_URL}/api/v1/users/login`, JSON.stringify({
        username: USERNAME,
        password: PASSWORD,
    }), {
        headers: { 'Content-Type': 'application/json' },
    })

    check(loginRes, { 'login success': (r) => r.status === 200 })
    const token = loginRes.json('data.access_token')
    const authHeaders = {
        headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
        }
    }

    // Step 2: Get movies
    const moviesRes = http.get(`${BASE_URL}/api/v1/movies`, authHeaders)
    check(moviesRes, { 'fetched movies': (r) => r.status === 200 })
    const movies = moviesRes.json('data')
    if (!movies || movies.length === 0) return

    const movie = randItem(movies)

    // Step 3: Get showtimes
    const showtimesRes = http.get(`${BASE_URL}/api/v1/movies/${movie.id}/showtimes`, authHeaders)
    check(showtimesRes, { 'fetched showtimes': (r) => r.status === 200 })
    const showtimes = showtimesRes.json('data')
    if (!showtimes || showtimes.length === 0) return

    const showtime = randItem(showtimes)

    // Step 4: Get seats for showtime
    const seatsRes = http.get(`${BASE_URL}/api/v1/showtimes/${showtime.id}/seats`, authHeaders)
    check(seatsRes, { 'fetched seats': (r) => r.status === 200 })
    const seats = seatsRes.json('data')
    const availableSeats = seats?.filter(s => s.is_available)
    if (!availableSeats || availableSeats.length === 0) return

    const seat = randItem(availableSeats)
    // 5. Hold seat
    const holdRes = http.post(`${BASE_URL}/api/v1/bookings/hold-seat`, JSON.stringify({
        showtime_id: showtime.id,
        seat_id: seat.id,
    }), authHeaders)

    check(holdRes, {
        'seat held': (r) => r.status === 200,
    })

    // Step 6: Book seat
    const bookingRes = http.post(`${BASE_URL}/api/v1/bookings`, JSON.stringify({
        showtime_id: showtime.id,
        seat_ids: [seat.id],
    }), authHeaders)

    check(bookingRes, {
        'booking success or conflict': (r) => r.status === 200 || r.status === 400,
    })

    if (bookingRes.status === 200) {
        console.log(`✅ Booked seat ${seat.id} in showtime ${showtime.id}`)
    } else if (bookingRes.status === 400) {
        console.warn(`⚠️ Seat ${seat.id} already booked by someone else`)
    }

    sleep(1) // Optional delay
}
