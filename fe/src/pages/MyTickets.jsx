import React, { useEffect, useState } from 'react'
import { myTickets } from '../api/booking'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import Typography from '@mui/material/Typography'

const MyTickets = () => {
  const [tickets, setTickets] = useState([])

  useEffect(() => {
    myTickets().then(res => setTickets(res.data.data || [])).catch(console.error)
  }, [])

  return (
    <div>
      <Typography variant="h5" gutterBottom>My Tickets</Typography>
      <List>
        {tickets.map(t => (
          <ListItem key={t.id}>
            {t.movie_title} - {new Date(t.showtime).toLocaleString()} - Seat {t.seat}
          </ListItem>
        ))}
      </List>
    </div>
  )
}

export default MyTickets