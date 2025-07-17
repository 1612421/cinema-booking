import {Env} from "../env";

// let socket

export const connectSocket = (token, onMessage) => {
    if (!token) return null

    const wsUrl = `${Env.ws}?access_token=${token}`

    const socket = new WebSocket(wsUrl)

    socket.onopen = () => {
        console.log('WebSocket connected')
    }

    socket.onerror = (err) => {
        console.error('WebSocket error:', err)
    }

    socket.onclose = () => {
        console.warn('🔌 WebSocket closed')
    }

    socket.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data)
            if (onMessage) {
                onMessage(msg)
            }
        } catch (err) {
            console.error('Invalid socket message:', err)
        }
    }

    return socket
}

// export const closeSocket = () => {
//     if (socket) {
//         console.log('testa asdasdasda 123123123123123')
//         socket.close()
//         socket = null
//     }
// }