import axios from './axios'

export const fetchShowtimesByMovie = (movieId) => axios.get(`/movies/${movieId}/showtimes`)
export const fetchSeatsByShowtime = (showtimeId) => axios.get(`/showtimes/${showtimeId}/seats`)