import React, { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import Typography from '@mui/material/Typography'
import Button from '@mui/material/Button'
import { fetchShowtimesByMovie } from '../api/showtime'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'

const Showtimes = () => {
  const { id, showtimeId } = useParams()
  const [showtimes, setShowtimes] = useState([])

  useEffect(() => {
    fetchShowtimesByMovie(id).then(res => setShowtimes(res.data.data || [])).catch(console.error)
  }, [id])

  return (
    <div>
      <Typography variant="h5" gutterBottom>Showtimes</Typography>
      <List>
        {showtimes.map(s => (
          <ListItem key={s.id}>
            <Button variant="outlined" component={Link} to={`/showtimes/${s.id}/seats`}>
              {new Date(s.start_time).toLocaleString()}
            </Button>
          </ListItem>
        ))}
      </List>
    </div>
  )
}

export default Showtimes