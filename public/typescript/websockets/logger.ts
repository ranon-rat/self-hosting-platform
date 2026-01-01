import type { LogMessage } from "../types/logs.js";
export function createLoggerWebSocket(idProject: number, onMessage: (data: LogMessage) => void) {
    const host = window.location.host
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
    
    // Get password from localStorage (like baseAPI does)
    const password = localStorage.getItem("password")
    if (!password) {
        throw new Error("Password not found")
    }
    
    // WebSocket in browsers doesn't support custom headers, so we pass password as query parameter
    // Note: The backend expects it in headers, but this is the only way in browsers
    const ws = new WebSocket(`${protocol}//${host}/logs/ws?id-project=${idProject}&password=${encodeURIComponent(password)}`)

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data) as LogMessage
        console.log(data)
        onMessage(data)
    }

    ws.onopen = () => {
        console.log("Connected to the logger websocket")
    }

    ws.onclose = () => {
        console.log("Disconnected from the logger websocket")
    }

    return ws
}