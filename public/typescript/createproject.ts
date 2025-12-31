import { CreateProject as CreateProjectAPI } from "./api/projects.js"

window.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("create-project-form")
    if (!form) {
        console.error("Form not found!")
        return
    }
    
    console.log("Form found, adding event listener")
    
    form.addEventListener("submit", async (event) => {
        event.preventDefault()
        console.log("Form submitted!")
        
        const name = form.querySelector("input[name='name']") as HTMLInputElement
        const dir = form.querySelector("input[name='dir']") as HTMLInputElement
        const command = form.querySelector("input[name='command']") as HTMLInputElement
        const thumbnail_url = form.querySelector("input[name='thumbnail_url']") as HTMLInputElement
        
        if (!name || !dir || !command || !thumbnail_url) {
            console.error("One or more form fields not found")
            return
        }
        
        const submitButton = form.querySelector("button[type='submit']") as HTMLButtonElement
        const originalText = submitButton?.textContent
        
        try {
            if (submitButton) {
                submitButton.disabled = true
                submitButton.textContent = "Creating..."
            }
            
            const projectData = {
                name: name.value,
                dir: dir.value,
                command: command.value,
                thumbnail_url: thumbnail_url.value
            }
            
            console.log("Sending project data:", projectData)
            
            const result = await CreateProjectAPI(projectData)
            console.log("Project created successfully!", result)
            
            // Redirect to index page after successful creation
            window.location.href = "/index.html"
        } catch (error) {
            console.error("Error creating project:", error)
            if (submitButton) {
                submitButton.disabled = false
                submitButton.textContent = originalText || "Create Project"
            }
            alert("Error creating project: " + (error instanceof Error ? error.message : "Unknown error"))
        }
    })
})
