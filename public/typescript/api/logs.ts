import type { Logs, PaginatedLogs } from "../types/logs.js";
import { baseAPI } from "./base.js";

export async function GetLogs(id: number) {
    return baseAPI<PaginatedLogs>("/logs?id=" + id, "GET")
}
// the websocket connectoin i will handle that on a different part