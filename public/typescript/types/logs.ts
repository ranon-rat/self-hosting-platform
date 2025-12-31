export interface Logs {
    id: number;
    id_project: number;
    content: string;
    created_at: string; // ISO 8601 date string
}

export interface PaginatedLogs {
    has_more: boolean;
    old_id: number;
    logs: Logs[];
}
export interface LogMessage {
    content: string;
    id_project: number;
}