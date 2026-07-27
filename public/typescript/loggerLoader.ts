import { createLoggerWebSocket } from "./websockets/logger.js"
import type { LogMessage, Logs } from "./types/logs.js"
import { GetLogs } from "./api/logs.js"
import { formatDate } from "./lib/date.js"

let oldId = 0
let loggerWebSocket: WebSocket | null = null

function escapeHtml(text: string): string {
    const div = document.createElement('div')
    div.textContent = text
    return div.innerHTML
}

function renderLog(log: Logs): string {
    const timestamp = formatDate(log.created_at)
    return `
        <div class="log-entry">
            <span class="log-timestamp">[${timestamp}]</span>
            <span class="log-content">${escapeHtml(log.content)}</span>
        </div>
    `
}

async function loadLogs(idProject: number) {
    const container = document.getElementById("project-logs-container")
    const moreLogsButton = document.getElementById("more-logs-button") as HTMLButtonElement
    
    if (!container) return
    
    let originalText = "📜 Load More Logs"
    
    try {
        if (moreLogsButton) {
            moreLogsButton.disabled = true
            originalText = moreLogsButton.textContent || "📜 Load More Logs"
            moreLogsButton.textContent = "Loading..."
        }
        
        const paginatedLogs = await GetLogs(oldId, idProject)
        
        if (paginatedLogs.logs.length > 0) {
            // Prepend older logs at the top (logs are already reversed from server)
            const logsHtml = paginatedLogs.logs.map(log => renderLog(log)).join("")
            container.innerHTML = logsHtml + container.innerHTML
            
            // Update oldId to the oldest ID we just loaded
            oldId = paginatedLogs.old_id
            
            // Update button visibility based on has_more
            if (moreLogsButton) {
                if (!paginatedLogs.has_more) {
                    moreLogsButton.classList.add("hidden")
                }
                moreLogsButton.disabled = false
                moreLogsButton.textContent = originalText
            }
        } else {
            if (moreLogsButton) {
                moreLogsButton.classList.add("hidden")
                moreLogsButton.disabled = false
                moreLogsButton.textContent = originalText
            }
        }
    } catch (error) {
        console.error("Error loading logs:", error)
        if (moreLogsButton) {
            moreLogsButton.disabled = false
            moreLogsButton.textContent = originalText
        }
    }
}

function renderRealtimeLog(logMessage: LogMessage) {
    const container = document.getElementById("project-logs-container")
    if (!container) return
    
    // Add new log at the bottom
    const logHtml = `
        <div class="log-entry">
            <span class="log-timestamp">[${formatDate(new Date().toISOString())}]</span>
            <span class="log-content">${escapeHtml(logMessage.content)}</span>
        </div>
    `
    container.innerHTML += logHtml
    
    // Auto-scroll to bottom
    container.scrollTop = container.scrollHeight
}

function cleanupLogger() {
    if (loggerWebSocket) {
        loggerWebSocket.close()
        loggerWebSocket = null
    }
}

window.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search)
    const id = urlParams.get("id")
    
    if (!id) {
        return
    }
    
    const projectId = parseInt(id)
    
    // Load initial logs
    await loadLogs(projectId)
    
    // Setup "Load More" button
    const moreLogsButton = document.getElementById("more-logs-button") as HTMLButtonElement
    moreLogsButton?.addEventListener("click", async () => {
        await loadLogs(projectId)
    })
    
    // Connect WebSocket for real-time logs
    try {
        loggerWebSocket = createLoggerWebSocket(projectId, (data: LogMessage) => {
            renderRealtimeLog(data)
        })
    } catch (error) {
        console.error("Error connecting to logger WebSocket:", error)
    }
    
    // Cleanup on page unload
    window.addEventListener("beforeunload", cleanupLogger)
})
