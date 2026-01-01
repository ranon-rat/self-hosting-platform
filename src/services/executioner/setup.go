package executionerServices

import (
	"github.com/ranon-rat/self-hosting-manager/src/domain/executionlogs"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
	"github.com/ranon-rat/self-hosting-manager/src/domain/repositories"
)

var pRepo projectsD.ProjectsRepoDB
var logRepo executionlogs.ExecutionLogsRepoDB
var executableEnv = []string{}

func Setup(repos *repositories.Repositories) {
	executableEnv = LoadLocalEnviroments()
	pRepo = repos.ProjectRepo
	logRepo = repos.LogRepo
	StartServices()
	go DeletingOldShit()
}
