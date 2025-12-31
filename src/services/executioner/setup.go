package executionerServices

import (
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
	"github.com/ranon-rat/self-hosting-manager/src/domain/repositories"
)

var pRepo projectsD.ProjectsRepoDB
var logRepo executionlogs.ExecutionLogsRepoDB

func Setup(repos *repositories.Repositories) {
	pRepo = repos.ProjectRepo
	logRepo = repos.LogRepo
}
