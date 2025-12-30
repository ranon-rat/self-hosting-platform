package projects

import "time"

/*
PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS project (
   ID INTEGER PRIMARY KEY,
   name VARCHAR(255) NOT NULL DEFAULT '',
   dir TEXT NOT NULL DEFAULT '', -- con esto puedo saber el entorno para ejecutar
   command TEXT NOT NULL DEFAULT '', -- con esto lo empiezo a correr
   thumbnail_url TEXT NOT NULL DEFAULT '',
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*/

type NewProject struct {
	Name         string `json:"name" db:"name"`
	Dir          string `json:"dir" db:"dir"`
	Command      string `json:"command" db:"command"`
	ThumbnailURL string `json:"thumbnail_url" db:"thumbnail_url"`
}

type Project struct {
	ID           int       `json:"id" db:"name"`
	Name         string    `json:"name" db:"name"`
	Dir          string    `json:"dir" db:"dir"`
	Command      string    `json:"command" db:"command"`
	ThumbnailURL string    `json:"thumbnail_url" db:"thumbnail_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
