package projectservice

import (
	"context"
	"fmt"
	"log/slog"
	"project-service/internal/domain/models"
	"project-service/internal/dto"
	"project-service/internal/repository/project"
)

type ProjectRepository interface {
	GetAllUserProjects(ctx context.Context, owner string, page, limit int64) ([]*models.Project, error)
	GetFilteredProjects(ctx context.Context, owner, status, namePrefix string, page, limit int64) ([]*models.Project, error)
	InitProject(ctx context.Context, composeId, owner, name string) (*models.Project, error)
	UpdateProject(ctx context.Context, composeId string, data string) (*models.Project, error)
	UpdateProjectStatus(ctx context.Context, string, status string) (*models.Project, error)
	UpdateProjectUrlZip(ctx context.Context, composeId string, url string) (*models.Project, error)
	UpdateProjectUrlDeploy(ctx context.Context, composeId string, url string) (*models.Project, error)
	DeleteProject(ctx context.Context, composeId string) error
}

type ProjectService struct {
	log               *slog.Logger
	projectRepository ProjectRepository
	projectUpdater    *ProjectUpdater
}

func NewProjectService(
	log *slog.Logger,
	projectRepository *project.ProjectRepository,
	projectUpdater *ProjectUpdater,
) *ProjectService {
	return &ProjectService{
		log:               log,
		projectRepository: projectRepository,
		projectUpdater:    projectUpdater,
	}
}

func (s *ProjectService) GetAllUserProjects(
	ctx context.Context,
	owner string,
	page, limit int64,
) ([]*models.Project, error) {
	projects, err := s.projectRepository.GetAllUserProjects(ctx, owner, page, limit)
	if err != nil {
		s.log.Error("ошибка при получении списка проектов", "error", err)
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) GetFilteredProjects(
	ctx context.Context,
	owner, status, namePrefix string,
	page, limit int64,
) ([]*models.Project, error) {
	projects, err := s.projectRepository.GetFilteredProjects(ctx, owner, status, namePrefix, page, limit)
	if err != nil {
		s.log.Error("ошибка при получении отфильтрованных проектов", "error", err)
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) InitProject(
	ctx context.Context,
	owner string,
	name string,
) (*models.Project, error) {
	composeId := toComposeId(owner, name)
	projectEntity, err := s.projectRepository.InitProject(ctx, composeId, owner, name)
	if err != nil {
		s.log.Error("ошибка при инициализации проекта", "error", err)
		return nil, err
	}

	return projectEntity, nil
}

func (s *ProjectService) UpdateProject(
	ctx context.Context,
	owner string,
	name string,
	data string,
) (*models.Project, error) {
	composeId := toComposeId(owner, name)
	projectEntity, err := s.projectRepository.UpdateProject(ctx, composeId, data)
	if err != nil {
		s.log.Error("ошибка при обновлении проекта", "error", err)
		return nil, err
	}

	return projectEntity, nil
}

func (s *ProjectService) DeleteProject(
	ctx context.Context,
	owner string,
	name string,
) error {
	composeId := toComposeId(owner, name)
	err := s.projectRepository.DeleteProject(ctx, composeId)
	if err != nil {
		s.log.Error("ошибка при удалении проекта", "error", err)
		return err
	}

	return nil
}

func (s *ProjectService) UpdateProjectStatus(
	ctx context.Context,
	dto dto.ProjectStatusDTO,
) (bool, error) {
	updProject, err := s.projectRepository.UpdateProjectStatus(ctx, dto.Id, dto.Status)
	if err != nil {
		s.log.Error("ошибка при обновлении статуса проекта", "error", err)
		return false, err
	}

	s.projectUpdater.Publish(updProject)

	return true, nil
}

func (s *ProjectService) UpdateProjectUrlZip(
	ctx context.Context,
	dto dto.NewZipDTO,
) (bool, error) {
	composeId := toComposeId(dto.Owner, dto.Name)
	updProject, err := s.projectRepository.UpdateProjectUrlZip(ctx, composeId, dto.Url)
	if err != nil {
		s.log.Error("ошибка при обновлении url zip проекта", "error", err)
		return false, err
	}

	s.projectUpdater.Publish(updProject)

	return true, nil
}

func (s *ProjectService) UpdateProjectUrlDeploy(
	ctx context.Context,
	dto dto.DeployPayloadDTO,
) (bool, error) {
	composeId := toComposeId(dto.Owner, dto.Name)
	updProject, err := s.projectRepository.UpdateProjectUrlDeploy(ctx, composeId, dto.Url)
	if err != nil {
		s.log.Error("ошибка при обновлении url zip проекта", "error", err)
		return false, err
	}

	s.projectUpdater.Publish(updProject)

	return true, nil
}

func toComposeId(owner, name string) string {
	return fmt.Sprintf("%s_%s", owner, name)
}
