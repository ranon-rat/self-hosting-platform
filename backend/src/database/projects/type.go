package projectsDB

import (
	"github.com/jmoiron/sqlx"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepo(db *sqlx.DB) projectsD.ProjectsRepoDB {
	return Repository{
		DB: db,
	}
}
