import { useState, useEffect } from 'react'
import './App.css'
import api from './api'
import AccessLog from './AccessLog'

function App() {
  const [clients, setClients] = useState([])

  const fetchData = async () => {
    const res = await api.getClients()
    setClients(res.clients.sort((a, b) => {
      let aLastSeen = 0;
      a.access_logs.forEach((log) => {
        aLastSeen = Math.max(aLastSeen,log.last_seen)
      })
      let bLastSeen = 0;
      b.access_logs.forEach((log) => {
        bLastSeen = Math.max(bLastSeen,log.last_seen)
      })
      if (aLastSeen != bLastSeen) {
        return bLastSeen - aLastSeen
      }
      if (a.id < b.id) {
        return 1
      } else if (a.id > b.id) {
        return -1
      } else {
        return 0
      }
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
      {clients.flatMap((client) => {
        return client.access_logs.map((accessLog) => {
          const id = `${client.id}-${accessLog.host}-${accessLog.ip_info.ip}`
          return <AccessLog client={client} accessLog={accessLog} key={id}></AccessLog>
        })
      })}
    </div>
  )
}

export default App
