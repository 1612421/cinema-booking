import React, { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import Typography from '@mui/material/Typography'
import Button from '@mui/material/Button'
import { fetchMovie } from '../api/movie'

const MovieDetail = () => {
  const { id } = useParams()
  const [movie, setMovie] = useState(null)

  useEffect(() => {
    fetchMovie(id).then(res => setMovie(res.data.data)).catch(console.error)
  }, [id])

  if (!movie) return null

  return (
    <div>
      <Typography variant="h4" gutterBottom>{movie.title}</Typography>
      <Typography variant="body1" gutterBottom>{movie.description}</Typography>
      <Button variant="contained" component={Link} to={`/movies/${id}/showtimes/`}>
        View Showtimes
      </Button>
    </div>
  )
}

export default MovieDetail