const baseURL = 'http://localhost:9502/api'
// const baseURL = '/api'

export default {
    getClients: async () => {
        const res = await fetch(`${baseURL}/clients`)
        return await res.json()
    },

    allowClient: async (id) => {
        const res = await fetch(`${baseURL}/client/${id}/allow`)
        return await res.json()
    },

    blockClient: async (id) => {
        const res = await fetch(`${baseURL}/client/${id}/block`)
        return await res.json()
    }
}