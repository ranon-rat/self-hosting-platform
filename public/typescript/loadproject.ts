import { GetProjectById, PauseProject, UpdateProject } from "./api/projects.js"
import type { Project, UpdateProject as UpdateProjectType } from "./types/projects.js"
import { formatDate } from "./lib/date.js"

let currentProject: Project | null = null
let isEditing = false

function escapeHtmlContent(text: string): string {
    return text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
}

function renderProjectDetails(project: Project, editing: boolean = false): string {
    if (editing) {
        return renderEditForm(project)
    }
    
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
                    <div class="px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white font-mono text-sm overflow-x-auto">
                        <pre class="whitespace-pre-wrap break-all">${escapeHtmlContent(project.command)}</pre>
                    </div>
                </div>
            </div>
            
            <div class="space-y-2">
                <label class="block text-sm font-medium text-slate-300">Thumbnail URL</label>
                <div class="px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white text-sm break-all">
                    ${project.thumbnail_url || "No thumbnail URL"}
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

function renderEditForm(project: Project): string {
    return `
        <div class="bg-white/10 backdrop-blur-xl rounded-2xl shadow-2xl border border-white/20 p-8 space-y-6">
            <div class="flex items-start justify-between">
                <div class="flex-1">
                    <h1 class="text-4xl font-bold text-white mb-2">Edit Project</h1>
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
            
            <form id="edit-project-form" class="space-y-6">
                <div class="space-y-2">
                    <label for="edit-name" class="block text-sm font-medium text-slate-300">Project Name</label>
                    <input 
                        type="text" 
                        id="edit-name"
                        name="name"
                        value="${project.name}"
                        class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200"
                        required
                    >
                </div>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div class="space-y-2">
                        <label for="edit-dir" class="block text-sm font-medium text-slate-300">Directory</label>
                        <input 
                            type="text" 
                            id="edit-dir"
                            name="dir"
                            value="${project.dir}"
                            class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white font-mono text-sm placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200"
                            required
                        >
                    </div>
                    
                    <div class="space-y-2">
                        <label for="edit-command" class="block text-sm font-medium text-slate-300">Command</label>
                        <textarea
                            id="edit-command"
                            name="command"
                            class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white font-mono text-sm placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200 resize-y min-h-[80px]"
                            required
                        >${escapeHtmlContent(project.command)}</textarea>
                    </div>
                </div>
                
                <div class="space-y-2">
                    <label for="edit-thumbnail_url" class="block text-sm font-medium text-slate-300">Thumbnail URL</label>
                    <input 
                        type="url" 
                        id="edit-thumbnail_url"
                        name="thumbnail_url"
                        value="${project.thumbnail_url || ''}"
                        class="w-full px-4 py-3 bg-white/5 border border-white/10 rounded-lg text-white text-sm placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200"
                    >
                </div>
                
                <div class="flex gap-4 pt-4">
                    <button 
                        type="button"
                        id="cancel-edit-button"
                        class="flex-1 bg-slate-700/50 hover:bg-slate-700 text-white font-semibold py-3 px-4 rounded-lg transition-all duration-200 transform hover:scale-[1.02] active:scale-[0.98]"
                    >
                        Cancel
                    </button>
                    <button 
                        type="submit"
                        class="flex-1 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-semibold py-3 px-4 rounded-lg transition-all duration-200 transform hover:scale-[1.02] active:scale-[0.98] shadow-lg hover:shadow-indigo-500/50"
                    >
                        Save Changes
                    </button>
                </div>
            </form>
        </div>
    `
}

function updateButtons(project: Project) {
    const pauseButton = document.getElementById("pause-project-button")
    const resumeButton = document.getElementById("resume-project-button")
    
    if (!pauseButton || !resumeButton) {
        return
    }
    
    // Hide pause/resume buttons when editing
    if (isEditing) {
        pauseButton.classList.add("hidden")
        resumeButton.classList.add("hidden")
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

async function handleEdit() {
    if (!currentProject) return
    
    isEditing = true
    renderProject()
    
    // Add event listeners for edit form
    const editForm = document.getElementById("edit-project-form")
    const cancelButton = document.getElementById("cancel-edit-button")
    
    cancelButton?.addEventListener("click", () => {
        isEditing = false
        renderProject()
    })
    
    editForm?.addEventListener("submit", async (event) => {
        event.preventDefault()
        await handleSaveEdit()
    })
}

async function handleSaveEdit() {
    if (!currentProject) return
    
    const form = document.getElementById("edit-project-form") as HTMLFormElement
    if (!form) return
    
    const name = form.querySelector("input[name='name']") as HTMLInputElement
    const dir = form.querySelector("input[name='dir']") as HTMLInputElement
    const command = form.querySelector("input[name='command']") as HTMLInputElement
    const thumbnail_url = form.querySelector("input[name='thumbnail_url']") as HTMLInputElement
    
    const submitButton = form.querySelector("button[type='submit']") as HTMLButtonElement
    const originalText = submitButton?.textContent
    
    try {
        if (submitButton) {
            submitButton.disabled = true
            submitButton.textContent = "Saving..."
        }
        
        const updateData: UpdateProjectType = {
            id: currentProject.id,
            name: name.value,
            dir: dir.value,
            command: command.value,
            thumbnail_url: thumbnail_url.value
        }
        
        await UpdateProject(updateData)
        
        // Reload project data
        currentProject = await GetProjectById(currentProject.id)
        if (currentProject) {
            isEditing = false
            renderProject()
        }
    } catch (error) {
        console.error("Error updating project:", error)
        alert("Error updating project. Please try again.")
    } finally {
        if (submitButton) {
            submitButton.disabled = false
            submitButton.textContent = originalText || "Save Changes"
        }
    }
}

function renderProject() {
    if (!currentProject) return
    
    const detailsContainer = document.getElementById("project-details")
    if (detailsContainer) {
        detailsContainer.innerHTML = renderProjectDetails(currentProject, isEditing)
    }
    
    updateButtons(currentProject)
    updateEditButton()
}

function updateEditButton() {
    const editButton = document.getElementById("edit-project-button")
    if (editButton) {
        if (isEditing) {
            editButton.classList.add("hidden")
        } else {
            editButton.classList.remove("hidden")
        }
    }
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
            const editButton = document.getElementById("edit-project-button")
            
            pauseButton?.addEventListener("click", handlePause)
            resumeButton?.addEventListener("click", handleResume)
            editButton?.addEventListener("click", handleEdit)
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
