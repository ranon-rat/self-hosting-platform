package repositories

import (
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
)

type Repositories struct {
	ProjectRepo projectsD.ProjectsRepoDB
	LogRepo     executionlogs.ExecutionLogsRepoDB
}
