package executionlogsDB

import (
	"github.com/jmoiron/sqlx"
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepo(db *sqlx.DB) executionlogs.ExecutionLogsRepoDB {
	return Repository{
		DB: db,
	}
}
