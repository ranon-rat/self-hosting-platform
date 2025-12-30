package projectServices

import (
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
)

var pRepo projectsD.ProjectsRepoDB
var logRepo executionlogs.ExecutionLogsRepoDB

func Setup(projectRepo projectsD.ProjectsRepoDB, executionLogs executionlogs.ExecutionLogsRepoDB) {
	pRepo = projectRepo
	logRepo = executionLogs
}
