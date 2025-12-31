package executionlogsDB

import (
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
func (r Repository) Get(oldId int) ([]executionlogs.Logs, error) {
	logs := []executionlogs.Logs{}
	query := baseSelectLogQuery
	args := []any{}
	if oldId != 0 {
		query += ` WHERE el.id < ?`
		args = append(args, oldId)
	}
	query += ` ORDER BY el.id DESC LIMIT ?`
	args = append(args, domain.LIMIT_PAGE)

	err := r.DB.Select(&logs, query, args...)
	return logs, err
}
