package projects

type ProjectsRepoDB interface {
	Create(project *NewProject) error
	GetByID(id int) (*Project, error)
	Search(search string) ([]Project, error)
}
