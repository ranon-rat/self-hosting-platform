
export async function baseAPI<T>(url: string, method: string, body?: any) {
    const password = localStorage.getItem("password")
    if (!password) {
        throw new Error("Password not found")
    }
    const response = await fetch(url, {
        method,
        body: body ? JSON.stringify(body) : undefined,
        headers: {
            "Content-Type": "application/json",
            "Password": password
        },
    })
    if (!response.ok) {
        if (response.status === 401) {
            localStorage.removeItem("password")
            window.location.href = "/login.html"
        }
        throw new Error("Failed to fetch data")
    }
    return response.json() as Promise<T>
}