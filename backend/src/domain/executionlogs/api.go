package executionlogs

import "time"

/*
CREATE TABLE IF NOT EXISTS execution_logs(
   ID INTEGER PRIMARY KEY,
   id_project INTEGER REFERENCES PROJECT(id) ON DELETE CASCADE,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   content TEXT
);
*/

type NewLog struct {
	IdProject int    `db:"id_project"`
	Content   string `db:"content"`
}

type Logs struct {
	IdProject int       `db:"id_project" json:"id_project"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type PaginatedLogs struct {
	HasMore bool `json:"has_more"`
	// lo del first id seria el primero que se ha cargado
	FirstID int    `json:"first_id"`
	Logs    []Logs `json:"logs"`
}
