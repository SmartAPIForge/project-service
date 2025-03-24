package codegenservice

import (
	"fmt"
	"log/slog"
	"project-service/internal/domain/models"
	"project-service/internal/dto"
	"project-service/internal/repository/project"
)

type ProjectRepository interface {
}

type ProjectService struct {
	log               *slog.Logger
	projectRepository ProjectRepository
}

func NewProjectService(
	log *slog.Logger,
	projectRepository *project.ProjectRepository,
) *ProjectService {
	return &ProjectService{
		log:               log,
		projectRepository: projectRepository,
	}
}

func (*ProjectService) GetUniqueUserProject(
	owner string,
	name string,
) (*models.Project, error) {
	return nil, nil
}

func (*ProjectService) GetAllUserProjects(
	owner string,
) ([]*models.Project, error) {
	return nil, nil
}

func (*ProjectService) CreateNewProject(
	owner string,
	name string,
) (*models.Project, error) {
	return nil, nil
}

func (*ProjectService) UpdateProject(
	owner string,
	name string,
	data map[string]interface{},
) (*models.Project, error) {
	return nil, nil
}

func (*ProjectService) GetProjectStatus(
	owner string,
	name string,
) (string, error) {
	return "nil", nil
}

func (*ProjectService) UpdateProjectStatus(
	dto dto.ProjectStatusDTO,
) (bool, error) {
	fmt.Println(dto)
	return true, nil
}
