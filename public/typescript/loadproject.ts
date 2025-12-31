import { GetProjectById, PauseProject } from "./api/projects.js"
import type { Project } from "./types/projects.js"
import { formatDate } from "./lib/date.js"

let currentProject: Project | null = null

function renderProjectDetails(project: Project): string {
    return `
        <div class="bg-white/10 backdrop-blur-xl rounded-2xl shadow-2xl border border-white/20 p-8 space-y-6">
            <div class="flex items-start justify-between">
                <div class="flex-1">
                    <h1 class="text-4xl font-bold text-white mb-2">${project.name}</h1>
                    <p class="text-slate-400 text-sm font-mono">ID: ${project.id}</p>
                </div>
                <span class="px-4 py-2 rounded-full text-sm font-semibold ${
                    project.is_paused 
                        ? "bg-yellow-500/20 text-yellow-300 border border-yellow-500/50" 
                        : "bg-green-500/20 text-green-300 border border-green-500/50"
                }">
                    ${project.is_paused ? "⏸ Paused" : "▶ Running"}
                </span>
            </div>
            
            ${project.thumbnail_url ? `
                <div class="rounded-lg overflow-hidden bg-slate-800/50">
                    <img 
                        src="${project.thumbnail_url}" 
                        alt="${project.name}" 
                        class="w-full h-64 object-cover"
                        onerror="this.style.display='none'"
                    />
                </div>
            ` : ''}
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-2">
                    <label class="block text-sm font-medium text-slate-300">Directory</label>
                    <div class="px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white font-mono text-sm">
                        ${project.dir}
                    </div>
                </div>
                
                <div class="space-y-2">
                    <label class="block text-sm font-medium text-slate-300">Command</label>
                    <div class="px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white font-mono text-sm">
                        ${project.command}
                    </div>
                </div>
            </div>
            
            <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Created At</label>
                <div class="px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white text-sm">
                    ${formatDate(project.created_at)}
                </div>
            </div>
        </div>
    `
}

function updateButtons(project: Project) {
    const pauseButton = document.getElementById("pause-project-button")
    const resumeButton = document.getElementById("resume-project-button")
    
    if (!pauseButton || !resumeButton) {
        return
    }
    
    if (project.is_paused) {
        pauseButton.classList.add("hidden")
        resumeButton.classList.remove("hidden")
    } else {
        pauseButton.classList.remove("hidden")
        resumeButton.classList.add("hidden")
    }
}

async function handlePause() {
    if (!currentProject) return
    
    const submitButton = document.getElementById("pause-project-button") as HTMLButtonElement
    const originalText = submitButton?.textContent
    
    try {
        if (submitButton) {
            submitButton.disabled = true
            submitButton.textContent = "Pausing..."
        }
        
        await PauseProject(currentProject.id, true)
        
        // Reload project data
        currentProject = await GetProjectById(currentProject.id)
        if (currentProject) {
            renderProject()
        }
    } catch (error) {
        console.error("Error pausing project:", error)
        alert("Error pausing project. Please try again.")
    } finally {
        if (submitButton) {
            submitButton.disabled = false
            submitButton.textContent = originalText || "Pause Project"
        }
    }
}

async function handleResume() {
    if (!currentProject) return
    
    const submitButton = document.getElementById("resume-project-button") as HTMLButtonElement
    const originalText = submitButton?.textContent
    
    try {
        if (submitButton) {
            submitButton.disabled = true
            submitButton.textContent = "Resuming..."
        }
        
        await PauseProject(currentProject.id, false)
        
        // Reload project data
        currentProject = await GetProjectById(currentProject.id)
        if (currentProject) {
            renderProject()
        }
    } catch (error) {
        console.error("Error resuming project:", error)
        alert("Error resuming project. Please try again.")
    } finally {
        if (submitButton) {
            submitButton.disabled = false
            submitButton.textContent = originalText || "Resume Project"
        }
    }
}

function renderProject() {
    if (!currentProject) return
    
    const detailsContainer = document.getElementById("project-details")
    if (detailsContainer) {
        detailsContainer.innerHTML = renderProjectDetails(currentProject)
    }
    
    updateButtons(currentProject)
}

window.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search)
    const id = urlParams.get("id")
    
    if (!id) {
        const detailsContainer = document.getElementById("project-details")
        if (detailsContainer) {
            detailsContainer.innerHTML = `
                <div class="bg-white/10 backdrop-blur-xl rounded-2xl shadow-2xl border border-white/20 p-8 text-center">
                    <p class="text-white text-xl">No project ID provided</p>
                    <a href="/index.html" class="text-indigo-400 hover:text-indigo-300 mt-4 inline-block">Go back to projects</a>
                </div>
            `
        }
        return
    }
    
    try {
        currentProject = await GetProjectById(parseInt(id))
        if (currentProject) {
            renderProject()
            
            // Add event listeners for buttons
            const pauseButton = document.getElementById("pause-project-button")
            const resumeButton = document.getElementById("resume-project-button")
            
            pauseButton?.addEventListener("click", handlePause)
            resumeButton?.addEventListener("click", handleResume)
        }
    } catch (error) {
        console.error("Error loading project:", error)
        const detailsContainer = document.getElementById("project-details")
        if (detailsContainer) {
            detailsContainer.innerHTML = `
                <div class="bg-white/10 backdrop-blur-xl rounded-2xl shadow-2xl border border-white/20 p-8 text-center">
                    <p class="text-red-400 text-xl">Error loading project</p>
                    <a href="/index.html" class="text-indigo-400 hover:text-indigo-300 mt-4 inline-block">Go back to projects</a>
                </div>
            `
        }
    }
})
