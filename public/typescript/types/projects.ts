export interface NewProject {
    name: string;
    dir: string;
    command: string;
    thumbnail_url: string;
}

export interface UpdateProject extends NewProject {
    id: number;
}

export interface Project {
    id: number;
    name: string;
    dir: string;
    command: string;
    thumbnail_url: string;
    created_at: string; // ISO 8601 date string
    is_paused: boolean;
}