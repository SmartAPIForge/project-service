package codegenservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"project-service/internal/domain/models"
	"project-service/internal/dto"
	"project-service/internal/repository/project"
)

type ProjectRepository interface {
	GetProjectByID(ctx context.Context, id project.ProjectUniqueIdentifier) (*project.Project, error)
	GetAllUserProjects(ctx context.Context, owner string, page, limit int64) ([]*project.Project, error)
	InitProject(ctx context.Context, id project.ProjectUniqueIdentifier) (*project.Project, error)
	UpdateProject(ctx context.Context, id project.ProjectUniqueIdentifier, data map[string]interface{}) (*project.Project, error)
	UpdateProjectStatus(ctx context.Context, id project.ProjectUniqueIdentifier, status project.ProjectStatus) error
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

// GetUniqueUserProject получает уникальный проект пользователя
func (s *ProjectService) GetUniqueUserProject(
	ctx context.Context,
	owner string,
	name string,
) (*models.Project, error) {
	id := project.ProjectUniqueIdentifier{
		Owner: owner,
		Name:  name,
	}

	projectEntity, err := s.projectRepository.GetProjectByID(ctx, id)
	if err != nil {
		s.log.Error("ошибка при получении проекта", "error", err)
		return nil, err
	}

	if projectEntity == nil {
		return nil, errors.New("проект не найден")
	}

	return mapProjectEntityToModel(projectEntity), nil
}

// GetAllUserProjects получает все проекты пользователя с пагинацией
func (s *ProjectService) GetAllUserProjects(
	ctx context.Context,
	owner string,
	page, limit int64,
) ([]*models.Project, error) {
	projectEntities, err := s.projectRepository.GetAllUserProjects(ctx, owner, page, limit)
	if err != nil {
		s.log.Error("ошибка при получении списка проектов", "error", err)
		return nil, err
	}

	projects := make([]*models.Project, len(projectEntities))
	for i, entity := range projectEntities {
		projects[i] = mapProjectEntityToModel(entity)
	}

	return projects, nil
}

// InitProject инициализирует новый проект
func (s *ProjectService) InitProject(
	ctx context.Context,
	owner string,
	name string,
) (*models.Project, error) {
	id := project.ProjectUniqueIdentifier{
		Owner: owner,
		Name:  name,
	}

	// Создаем новый проект
	projectEntity, err := s.projectRepository.InitProject(ctx, id)
	if err != nil {
		s.log.Error("ошибка при инициализации проекта", "error", err)
		return nil, err
	}

	return mapProjectEntityToModel(projectEntity), nil
}

// UpdateProject обновляет данные проекта
func (s *ProjectService) UpdateProject(
	ctx context.Context,
	owner string,
	name string,
	data map[string]interface{},
) (*models.Project, error) {
	id := project.ProjectUniqueIdentifier{
		Owner: owner,
		Name:  name,
	}

	// Обновляем проект
	projectEntity, err := s.projectRepository.UpdateProject(ctx, id, data)
	if err != nil {
		s.log.Error("ошибка при обновлении проекта", "error", err)
		return nil, err
	}

	return mapProjectEntityToModel(projectEntity), nil
}

// GetProjectStatus получает статус проекта
func (s *ProjectService) GetProjectStatus(
	ctx context.Context,
	owner string,
	name string,
) (string, error) {
	id := project.ProjectUniqueIdentifier{
		Owner: owner,
		Name:  name,
	}

	projectEntity, err := s.projectRepository.GetProjectByID(ctx, id)
	if err != nil {
		s.log.Error("ошибка при получении статуса проекта", "error", err)
		return "", err
	}

	if projectEntity == nil {
		return "", errors.New("проект не найден")
	}

	return mapStatusToString(projectEntity.Status), nil
}

// UpdateProjectStatus обновляет статус проекта
func (s *ProjectService) UpdateProjectStatus(
	ctx context.Context,
	dto dto.ProjectStatusDTO,
) (bool, error) {
	// Парсим ID из DTO
	parts := dto.Id
	if parts == "" {
		return false, errors.New("некорректный ID проекта")
	}

	owner, name := parseProjectID(parts)
	id := project.ProjectUniqueIdentifier{
		Owner: owner,
		Name:  name,
	}

	// Конвертируем строковый статус в ProjectStatus
	status, err := parseStatusFromString(dto.Status)
	if err != nil {
		return false, err
	}

	// Обновляем статус
	err = s.projectRepository.UpdateProjectStatus(ctx, id, status)
	if err != nil {
		s.log.Error("ошибка при обновлении статуса проекта", "error", err)
		return false, err
	}

	return true, nil
}

// Вспомогательные функции

// mapProjectEntityToModel конвертирует сущность проекта в модель
func mapProjectEntityToModel(entity *project.Project) *models.Project {
	// Конвертируем данные проекта в JSON
	data, _ := json.Marshal(entity.Data)

	return &models.Project{
		Owner:  entity.ComposeID.Owner,
		Name:   entity.ComposeID.Name,
		Status: mapStatusToString(entity.Status),
		Data:   data,
	}
}

// mapStatusToString конвертирует статус в строку
func mapStatusToString(status project.ProjectStatus) string {
	switch status {
	case project.NEW:
		return "NEW"
	case project.GENERATE_PENDING:
		return "GENERATE_PENDING"
	case project.GENERATE_SUCCESS:
		return "GENERATE_SUCCESS"
	case project.GENERATE_FAIL:
		return "GENERATE_FAIL"
	case project.DEPLOY_PENDING:
		return "DEPLOY_PENDING"
	case project.DEPLOY_SUCCESS:
		return "DEPLOY_SUCCESS"
	case project.DEPLOY_FAIL:
		return "DEPLOY_FAIL"
	case project.RUNNING:
		return "RUNNING"
	case project.STOPPED:
		return "STOPPED"
	case project.FAILED:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

// parseStatusFromString конвертирует строку в статус
func parseStatusFromString(statusStr string) (project.ProjectStatus, error) {
	switch statusStr {
	case "NEW":
		return project.NEW, nil
	case "GENERATE_PENDING":
		return project.GENERATE_PENDING, nil
	case "GENERATE_SUCCESS":
		return project.GENERATE_SUCCESS, nil
	case "GENERATE_FAIL":
		return project.GENERATE_FAIL, nil
	case "DEPLOY_PENDING":
		return project.DEPLOY_PENDING, nil
	case "DEPLOY_SUCCESS":
		return project.DEPLOY_SUCCESS, nil
	case "DEPLOY_FAIL":
		return project.DEPLOY_FAIL, nil
	case "RUNNING":
		return project.RUNNING, nil
	case "STOPPED":
		return project.STOPPED, nil
	case "FAILED":
		return project.FAILED, nil
	default:
		return 0, fmt.Errorf("неизвестный статус: %s", statusStr)
	}
}

// parseProjectID разбирает ID проекта на owner и name
func parseProjectID(id string) (owner, name string) {
	// Логика парсинга ID должна быть реализована согласно формату хранения
	// Например, если ID хранится в формате "owner/name"
	// В данном случае, предполагается, что ID уже содержит owner и name
	// TODO: Реализовать корректный парсинг ID
	return id, id
}
