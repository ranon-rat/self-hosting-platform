import { Login } from "./api/public.js"

window.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("login-form")
    if (!loginForm) {
        return
    }
    
    const passwordInput = document.getElementById("password") as HTMLInputElement
    const submitButton = loginForm.querySelector("button[type='submit']") as HTMLButtonElement
    
    // Create error message container
    let errorMessage: HTMLElement | null = null
    
    const showError = (message: string) => {
        if (!errorMessage) {
            errorMessage = document.createElement("div")
            errorMessage.className = "bg-red-500/20 border border-red-500/50 text-red-200 px-4 py-3 rounded-lg text-sm"
            loginForm.insertBefore(errorMessage, loginForm.firstChild)
        }
        errorMessage.textContent = message
        errorMessage.classList.remove("hidden")
        passwordInput.classList.add("border-red-500")
    }
    
    const hideError = () => {
        if (errorMessage) {
            errorMessage.classList.add("hidden")
        }
        passwordInput.classList.remove("border-red-500")
    }
    
    loginForm.addEventListener("submit", async (event) => {
        event.preventDefault()
        if (!passwordInput) {
            return
        }
        
        hideError()
        
        // Disable button and show loading state
        submitButton.disabled = true
        const originalText = submitButton.textContent
        submitButton.textContent = "Signing in..."
        submitButton.classList.add("opacity-75", "cursor-not-allowed")
        
        localStorage.setItem("password", passwordInput.value)
        try {

            await Login(passwordInput.value)
            debugger
            window.location.href = "/index.html"
        } catch (error) {
            localStorage.removeItem("password")
            showError("Invalid password. Please try again.")
            passwordInput.value = ""
            passwordInput.focus()
        } finally {
            submitButton.disabled = false
            submitButton.textContent = originalText || "Sign In"
            submitButton.classList.remove("opacity-75", "cursor-not-allowed")
        }
    })
})