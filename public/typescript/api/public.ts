import { baseAPI } from "./base.js";

export async function Login(password: string) {
    return baseAPI<{message: string}>("/public/login?password=" + password, "GET")
}