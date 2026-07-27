import { GetProjectById, PauseProject, UpdateProject, DeleteProject } from "./api/projects.js"
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
        <div class="card stack">
            <div class="details-header">
                <div>
                    <h1 class="title">${project.name}</h1>
                    <p class="subtitle text-sm mono">ID: ${project.id}</p>
                </div>
                <span class="badge ${
                    project.is_paused
                        ? "badge-paused"
                        : "badge-running"
                }">
                    ${project.is_paused ? "⏸ Paused" : "▶ Running"}
                </span>
            </div>

            ${project.thumbnail_url ? `
                <div class="details-thumb">
                    <img
                        src="${project.thumbnail_url}"
                        alt="${project.name}"
                        onerror="this.style.display='none'"
                    />
                </div>
            ` : ''}

            <div class="details-grid">
                <div class="field">
                    <label>Directory</label>
                    <div class="value-box mono">
                        ${project.dir}
                    </div>
                </div>

                <div class="field">
                    <label>Command</label>
                    <div class="value-box mono" style="overflow-x: auto;">
                        <pre style="white-space: pre-wrap; word-break: break-all; margin: 0;">${escapeHtmlContent(project.command)}</pre>
                    </div>
                </div>
            </div>

            <div class="field">
                <label>Thumbnail URL</label>
                <div class="value-box">
                    ${project.thumbnail_url || "No thumbnail URL"}
                </div>
            </div>

            <div class="field">
                <label>Created At</label>
                <div class="value-box">
                    ${formatDate(project.created_at)}
                </div>
            </div>
        </div>
    `
}

function renderEditForm(project: Project): string {
    return `
        <div class="card stack">
            <div class="details-header">
                <div>
                    <h1 class="title">Edit Project</h1>
                    <p class="subtitle text-sm mono">ID: ${project.id}</p>
                </div>
                <span class="badge ${
                    project.is_paused
                        ? "badge-paused"
                        : "badge-running"
                }">
                    ${project.is_paused ? "⏸ Paused" : "▶ Running"}
                </span>
            </div>

            <form id="edit-project-form" class="form-fields">
                <div class="field">
                    <label for="edit-name">Project Name</label>
                    <input
                        type="text"
                        id="edit-name"
                        name="name"
                        value="${project.name}"
                        class="input"
                        required
                    >
                </div>

                <div class="details-grid">
                    <div class="field">
                        <label for="edit-dir">Directory</label>
                        <input
                            type="text"
                            id="edit-dir"
                            name="dir"
                            value="${project.dir}"
                            class="input mono text-sm"
                            required
                        >
                    </div>

                    <div class="field">
                        <label for="edit-command">Command</label>
                        <textarea
                            id="edit-command"
                            name="command"
                            class="input mono text-sm"
                            required
                        >${escapeHtmlContent(project.command)}</textarea>
                    </div>
                </div>

                <div class="field">
                    <label for="edit-thumbnail_url">Thumbnail URL</label>
                    <input
                        type="url"
                        id="edit-thumbnail_url"
                        name="thumbnail_url"
                        value="${project.thumbnail_url || ''}"
                        class="input text-sm"
                    >
                </div>

                <div class="buttons-row">
                    <button
                        type="button"
                        id="cancel-edit-button"
                        class="btn btn-secondary btn-grow"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        class="btn btn-primary btn-grow"
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
    const deleteButton = document.getElementById("delete-project-button")

    if (!pauseButton || !resumeButton) {
        return
    }

    // Hide pause/resume/delete buttons when editing
    if (isEditing) {
        pauseButton.classList.add("hidden")
        resumeButton.classList.add("hidden")
        deleteButton?.classList.add("hidden")
        return
    }

    deleteButton?.classList.remove("hidden")
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

async function handleDelete() {
    if (!currentProject) return

    if (!confirm(`Are you sure you want to delete "${currentProject.name}"? This cannot be undone.`)) {
        return
    }

    const deleteButton = document.getElementById("delete-project-button") as HTMLButtonElement

    try {
        if (deleteButton) {
            deleteButton.disabled = true
            deleteButton.textContent = "Deleting..."
        }

        await DeleteProject(currentProject.id)

        window.location.href = "/index.html"
    } catch (error) {
        console.error("Error deleting project:", error)
        alert("Error deleting project. Please try again.")
        if (deleteButton) {
            deleteButton.disabled = false
            deleteButton.textContent = "🗑 Delete Project"
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
    const command = form.querySelector("textarea[name='command']") as HTMLTextAreaElement
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
                <div class="card" style="text-align: center;">
                    <p class="empty-state-title">No project ID provided</p>
                    <a href="/index.html" class="link" style="margin-top: 1rem; display: inline-block;">Go back to projects</a>
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
            const deleteButton = document.getElementById("delete-project-button")

            pauseButton?.addEventListener("click", handlePause)
            resumeButton?.addEventListener("click", handleResume)
            editButton?.addEventListener("click", handleEdit)
            deleteButton?.addEventListener("click", handleDelete)
        }
    } catch (error) {
        console.error("Error loading project:", error)
        const detailsContainer = document.getElementById("project-details")
        if (detailsContainer) {
            detailsContainer.innerHTML = `
                <div class="card" style="text-align: center;">
                    <p class="empty-state-title" style="color: var(--color-red-400);">Error loading project</p>
                    <a href="/index.html" class="link" style="margin-top: 1rem; display: inline-block;">Go back to projects</a>
                </div>
            `
        }
    }
})
