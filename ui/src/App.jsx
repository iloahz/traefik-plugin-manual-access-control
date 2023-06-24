import { useState, useEffect } from 'react'
import './App.css'
import api from './api'
import Client from './Client'

function App() {
  const [clients, setClients] = useState([])

  useEffect(() => {
    api.getClients().then((res) => setClients(res.clients))
  }, [])

  return (
    <>
      {clients.map((client) => (
        <Client client={client} key={client.id}></Client>
      ))}
    </>
  )
}

export default App
