package persistence

import "learn_zinx/Cobra.mayfly/internal/project/domain/repository"

var (
	projectRepo       repository.Project        = newProjectRepo()
	projectEnvRepo    repository.ProjectEnv     = newProjectEnvRepo()
	projectMemberRepo repository.ProjectMemeber = newProjectMemberRepo()
)

func GetProjectRepo() repository.Project {
	return projectRepo
}

func GetProjectEnvRepo() repository.ProjectEnv {
	return projectEnvRepo
}

func GetProjectMemberRepo() repository.ProjectMemeber {
	return projectMemberRepo
}
