package executionlogsDB

import (
	"fmt"
	"time"

	"github.com/ranon-rat/self-hosting-manager/src/domain"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
)

/*
CREATE TABLE IF NOT EXISTS execution_logs(

	ID INTEGER PRIMARY KEY,
	id_project INTEGER NOT NULL REFERENCES PROJECT(id) ON DELETE CASCADE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	content TEXT NOT NULL DEFAULT ''

);
*/
const baseSelectLogQuery = `
SELECT 
	el.*
FROM execution_logs el
`

func (r Repository) Create(log *executionlogs.NewLog) error {
	query := `
	INSERT INTO execution_logs(id_project, content)
	VALUES (:id_project,:content)
	`
	_, err := r.DB.NamedExec(query, log)
	return err
}

// creo que esto deberia de funcionar?
func (r Repository) Get(oldId, projectID int) ([]executionlogs.Logs, error) {
	logs := []executionlogs.Logs{}
	query := baseSelectLogQuery
	args := []any{}
	query += ` WHERE el.id_project = ? `
	args = append(args, projectID)
	if projectID != 0 {
		query += ` AND el.id < ? `
		args = append(args, oldId)
	}
	query += ` ORDER BY el.id DESC LIMIT ?`
	args = append(args, domain.LIMIT_PAGE)
	fmt.Println(query, args)
	err := r.DB.Select(&logs, query, args...)
	return logs, err
}
func (r Repository) DeleteOldMessages(days int) error {
	query := `
	DELETE FROM execution_logs WHERE created_at < ?
	`
	_, err := r.DB.Exec(query, time.Now().AddDate(0, 0, -days))
	return err
}
