import type { Logs, PaginatedLogs } from "../types/logs.js";
import { baseAPI } from "./base.js";

export async function GetLogs(oldId: number) {
    return baseAPI<PaginatedLogs>("/logs", "GET", { oldId })
}
// the websocket connectoin i will handle that on a different part