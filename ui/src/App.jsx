import { useState, useEffect } from 'react'
import './App.css'
import api from './api'
import Client from './Client'

function App() {
  const [clients, setClients] = useState([])

  const fetchData = async () => {
    const res = await api.getClients()
    setClients(res.clients.sort((a, b) => {
      return b.access_logs.length - a.access_logs.length
      // let aLastSeen = 0;
      // a.access_logs.forEach((k, v) => {
      //   aLastSeen = Math.max(aLastSeen, v.last_seen)
      // })
      // let bLastSeen = 0;
      // b.access_logs.forEach((k, v) => {
      //   bLastSeen = Math.max(bLastSeen, v.last_seen)
      // })
      // return bLastSeen - aLastSeen
    }))
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
