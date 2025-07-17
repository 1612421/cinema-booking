import axios from './axios'

export const fetchMovies = () => axios.get('/movies')
export const fetchMovie = (id) => axios.get(`/movies/${id}`)