package projectsD

type ProjectsRepoDB interface {
	Create(project *NewProject) (int, error)
	UpdateProject(project *UpdateProject) error
	PauseProject(pause bool, id int) error

	GetByID(id int) (*Project, error)
	Search(search string) ([]Project, error)

	//
	Delete(id int) error
}
