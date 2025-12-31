import type { Project } from "./types/projects.js"
import { SearchProjects } from "./api/projects.js"
import { formatDate } from "./lib/date.js"


function renderLoading(): string {
    return `
        <div class="flex items-center justify-center min-h-[400px]">
            <div class="text-center space-y-4">
                <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-white/20 border-t-purple-500"></div>
                <p class="text-slate-300 text-sm">Loading projects...</p>
            </div>
        </div>
    `
}

function renderError(message: string): string {
    return `
        <div class="flex items-center justify-center min-h-[400px]">
            <div class="bg-red-500/20 border border-red-500/50 text-red-200 px-6 py-4 rounded-lg max-w-md text-center">
                <p class="font-semibold mb-2">Error loading projects</p>
                <p class="text-sm">${message}</p>
            </div>
        </div>
    `
}

function renderProjects(projects: Project[] | undefined): string {

    if (projects === undefined || projects.length === 0) {
        return `
            <div class="flex items-center justify-center min-h-[400px]">
                <div class="text-center space-y-2">
                    <p class="text-4xl mb-4">📦</p>
                    <p class="text-white text-xl font-semibold">No projects found</p>
                    <p class="text-slate-400 text-sm">Create your first project to get started</p>
                </div>
            </div>
        `
    }

    return `
        <div class="mb-8">
            <h1 class="text-4xl font-bold text-white mb-2">Projects</h1>
            <p class="text-slate-400">Manage your self-hosted applications</p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            ${projects.map((project: Project) => `
                <a 
                    href="/project.html?id=${project.id}" 
                    class="group block bg-white/10 backdrop-blur-xl rounded-xl border border-white/20 p-6 hover:bg-white/15 hover:border-purple-500/50 transition-all duration-200 transform hover:scale-[1.02] hover:shadow-xl hover:shadow-purple-500/20"
                >
                    <div class="flex items-start justify-between mb-4">
                        <div class="flex-1">
                            <h2 class="text-xl font-bold text-white mb-1 group-hover:text-purple-300 transition-colors">
                                ${project.name}
                            </h2>
                            <p class="text-slate-400 text-xs font-mono">ID: ${project.id}</p>
                        </div>
                        <span class="px-3 py-1 rounded-full text-xs font-semibold ${
                            project.is_paused 
                                ? "bg-yellow-500/20 text-yellow-300 border border-yellow-500/50" 
                                : "bg-green-500/20 text-green-300 border border-green-500/50"
                        }">
                            ${project.is_paused ? "⏸ Paused" : "▶ Running"}
                        </span>
                    </div>
                    
                    ${project.thumbnail_url ? `
                        <div class="mb-4 rounded-lg overflow-hidden bg-slate-800/50">
                            <img 
                                src="${project.thumbnail_url}" 
                                alt="${project.name}" 
                                class="w-full h-32 object-cover group-hover:scale-105 transition-transform duration-200"
                                onerror="this.style.display='none'"
                            />
                        </div>
                    ` : `
                        <div class="mb-4 rounded-lg bg-gradient-to-br from-purple-600/20 to-indigo-600/20 h-32 flex items-center justify-center">
                            <span class="text-4xl">📦</span>
                        </div>
                    `}
                    
                    <div class="space-y-2 text-sm">
                        <div class="flex items-center text-slate-400">
                            <span class="mr-2">📁</span>
                            <span class="truncate font-mono text-xs">${project.dir}</span>
                        </div>
                        <div class="flex items-center text-slate-400">
                            <span class="mr-2">⚡</span>
                            <span class="truncate font-mono text-xs">${project.command}</span>
                        </div>
                        <div class="flex items-center text-slate-400">
                            <span class="mr-2">📅</span>
                            <span class="text-xs">${formatDate(project.created_at)}</span>
                        </div>
                    </div>
                    
                    <div class="mt-4 pt-4 border-t border-white/10">
                        <span class="text-purple-400 text-sm font-medium group-hover:text-purple-300">
                            View details →
                        </span>
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