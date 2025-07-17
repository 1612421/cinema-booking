import React, { useEffect, useState } from 'react'
import Grid from '@mui/material/Grid'
import { fetchMovies } from '../api/movie'
import MovieCard from '../components/MovieCard'

const Home = () => {
  const [movies, setMovies] = useState([])

  useEffect(() => {
    fetchMovies().then(res => setMovies(res.data.data || [])).catch(console.error)
  }, [])

  return (
    <Grid container spacing={2}>
      {movies.map(movie => (
        <Grid item key={movie.id}>
          <MovieCard movie={movie} />
        </Grid>
      ))}
    </Grid>
  )
}

export default Home