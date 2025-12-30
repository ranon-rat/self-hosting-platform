package projectsDB

import (
	"fmt"

	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
)

/*
CREATE TABLE IF NOT EXISTS project (
   ID INTEGER PRIMARY KEY,
   name VARCHAR(255) NOT NULL DEFAULT '',
   dir TEXT NOT NULL DEFAULT '', -- con esto puedo saber el entorno para ejecutar
   command TEXT NOT NULL DEFAULT '', -- con esto lo empiezo a correr
   thumbnail_url TEXT NOT NULL DEFAULT '',
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   is_paused bool NOT NULL default false
);
*/

const BaseSelectProjectQuery = `
SELECT 
 p.*
FROM project p

%s
`

func (p Repository) Create(project *projectsD.NewProject) (int, error) {
	query := `
	INSERT INTO project (
		name,
		dir,
		command,
		thumbnail_url
	) VALUES (
		:name,
		:dir,
		:command,
		:thumbnail_url
	)
	RETURNING id
	`
	row, err := p.DB.NamedQuery(query, project)
	if err != nil {
		return 0, err
	}
	if !row.Next() {
		return 0, fmt.Errorf("no id returned from insert")
	}
	var id int
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil

}

func (p Repository) UpdateProject(project *projectsD.UpdateProject) error {
	query := `
	UPDATE project 
	SET 
		name=:name,
		dir=:dir,
		command=:command,
		thumbnail_url=:thumbnail_url
	WHERE id=:id
	`
	_, err := p.DB.NamedExec(query, project)
	return err
}
func (p Repository) GetByID(id int) (*projectsD.Project, error) {
	query := fmt.Sprintf(BaseSelectProjectQuery, "WHERE p.id=?1")
	project := new(projectsD.Project)
	err := p.DB.Get(project, query, id)
	return project, err
}

func (p Repository) Search(search string) ([]projectsD.Project, error) {
	query := fmt.Sprintf(BaseSelectProjectQuery, "WHERE LOWER(p.name) LIKE LOWER(?1) ")
	projects := []projectsD.Project{}
	err := p.DB.Select(&projects, query, search)
	return projects, err
}
func (p Repository) PauseProject(pause bool, id int) error {
	query := `
	UPDATE project SET
		is_paused=?1
	WHERE id=?2
	`
	_, err := p.DB.Exec(query, pause, id)
	return err
}
