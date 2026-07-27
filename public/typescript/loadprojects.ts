import type { Project } from "./types/projects.js"
import { SearchProjects } from "./api/projects.js"
import { formatDate } from "./lib/date.js"


function renderLoading(): string {
    return `
        <div class="state-box">
            <div class="state-box-inner">
                <div class="spinner"></div>
                <p class="subtitle text-sm" style="margin-top: 1rem;">Loading projects...</p>
            </div>
        </div>
    `
}

function renderError(message: string): string {
    return `
        <div class="state-box">
            <div class="error-box">
                <p style="font-weight: 600; margin: 0 0 0.5rem;">Error loading projects</p>
                <p class="text-sm" style="margin: 0;">${message}</p>
            </div>
        </div>
    `
}

function renderProjects(projects: Project[] | undefined): string {

    if (projects === undefined || projects.length === 0) {
        return `
            <div class="state-box">
                <div class="state-box-inner">
                    <p class="empty-state-icon">📦</p>
                    <p class="empty-state-title">No projects found</p>
                    <p class="empty-state-text">Create your first project to get started</p>
                </div>
            </div>
        `
    }

    return `
        <div style="margin-bottom: 2rem;">
            <h1 class="title">Projects</h1>
            <p class="subtitle">Manage your self-hosted applications</p>
        </div>
        <div class="projects-grid">
            ${projects.map((project: Project) => `
                <a
                    href="/project.html?id=${project.id}"
                    class="project-card"
                >
                    <div class="project-card-header">
                        <div>
                            <h2 class="project-card-name">
                                ${project.name}
                            </h2>
                            <p class="project-card-id">ID: ${project.id}</p>
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
                        <div class="project-card-thumb">
                            <img
                                src="${project.thumbnail_url}"
                                alt="${project.name}"
                                onerror="this.style.display='none'"
                            />
                        </div>
                    ` : `
                        <div class="project-card-thumb-placeholder">
                            <span>📦</span>
                        </div>
                    `}

                    <div class="project-card-meta">
                        <div class="project-card-meta-row">
                            <span>📁</span>
                            <span class="truncate">${project.dir}</span>
                        </div>
                        <div class="project-card-meta-row">
                            <span>⚡</span>
                            <span class="truncate">${project.command}</span>
                        </div>
                        <div class="project-card-meta-row">
                            <span>📅</span>
                            <span class="text-xs">${formatDate(project.created_at)}</span>
                        </div>
                    </div>

                    <div class="project-card-footer">
                        View details →
                    </div>
                </a>
            `).join("")}
        </div>
    `
}

window.addEventListener("DOMContentLoaded", async () => {
    const app = document.getElementById("app")
    if (!app) {
        return
    }

    // Show loading state
    app.innerHTML = renderLoading()

    try {
        const projects = await SearchProjects("")
        app.innerHTML = renderProjects(projects)
    } catch (error) {
        console.error("Error loading projects:", error)
        const errorMessage = error instanceof Error ? error.message : "Unknown error occurred"
        app.innerHTML = renderError(errorMessage)
    }
})
