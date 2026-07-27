import type { NewProject, UpdateProject, Project } from "../types/projects.js";
import { baseAPI } from "./base.js";


export async function CreateProject(project: NewProject) {
    return baseAPI<{message: string}>("/projects", "POST", project)
}
export async function UpdateProject(project: UpdateProject) {
    return baseAPI<{message: string}>("/projects", "PUT", project)
}
export async function PauseProject(id: number, pause: boolean) {
    return baseAPI<{message: string}>("/projects/pause?id=" + id + "&pause=" + pause, "PUT")
}
export async function GetProjectById(id: number) {
    return baseAPI<Project>("/projects/by-id?id=" + id, "GET" )
}
export async function SearchProjects(search: string) {
    return baseAPI<Project[]>("/projects?search=" + search, "GET")
}
export async function DeleteProject(id: number) {
    return baseAPI<{message: string}>("/projects?id=" + id, "DELETE")
}