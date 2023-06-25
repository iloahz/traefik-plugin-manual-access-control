import { useState, useEffect } from 'react'
import './App.css'
import api from './api'
import Client from './Client'

function App() {
  const [clients, setClients] = useState([])

  const fetchData = async () => {
    const res = await api.getClients()
    setClients(res.clients.sort((a, b) => b.stats.last_seen - a.stats.last_seen))
  }

  useEffect(() => {
    fetchData()
    
    const interval = setInterval(() => {
      fetchData()
    }, 500)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className='board'>
      {clients.map((client) => (
        <Client client={client} key={client.id}></Client>
      ))}
    </div>
  )
}

export default App
