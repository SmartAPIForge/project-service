package projectserver

import (
	"context"
	"encoding/json"
	"strconv"
	
	projectProto "github.com/SmartAPIForge/protos/gen/go/project"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	
	"project-service/internal/domain/models"
)

type ProjectService interface {
	GetUniqueUserProject(
		ctx context.Context,
		owner string,
		name string,
	) (*models.Project, error)
	GetAllUserProjects(
		ctx context.Context,
		owner string,
		page, limit int64,
	) ([]*models.Project, error)
	InitProject(
		ctx context.Context,
		owner string,
		name string,
	) (*models.Project, error)
	UpdateProject(
		ctx context.Context,
		owner string,
		name string,
		data map[string]interface{},
	) (*models.Project, error)
	GetProjectStatus(
		ctx context.Context,
		owner string,
		name string,
	) (string, error)
}

type ProjectServer struct {
	projectProto.UnsafeProjectServiceServer
	projectService ProjectService
}

func RegisterProjectServer(
	gRPCServer *grpc.Server,
	project ProjectService,
) {
	projectProto.RegisterProjectServiceServer(gRPCServer, &ProjectServer{projectService: project})
}

func (s *ProjectServer) GetUniqueUserProject(
	ctx context.Context,
	in *projectProto.GetUniqueUserProjectRequest,
) (*projectProto.ProjectResponse, error) {
	if in.ComposeId == nil {
		return nil, status.Error(codes.InvalidArgument, "не указан идентификатор проекта")
	}
	
	project, err := s.projectService.GetUniqueUserProject(ctx, in.ComposeId.Owner, in.ComposeId.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	return projectToResponse(project)
}

func (s *ProjectServer) GetAllUserProjects(
	ctx context.Context,
	in *projectProto.GetAllUserProjectsRequest,
) (*projectProto.ListOfProjectsResponse, error) {
	if in.Owner == "" {
		return nil, status.Error(codes.InvalidArgument, "не указан владелец проектов")
	}
	
	// Парсим параметры страницы и лимита
	page := int64(1)
	limit := int64(10)
	
	if in.Page != "" {
		parsedPage, err := strconv.ParseInt(in.Page, 10, 64)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	
	if in.Limit != "" {
		parsedLimit, err := strconv.ParseInt(in.Limit, 10, 64)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	projects, err := s.projectService.GetAllUserProjects(ctx, in.Owner, page, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	// Преобразуем список моделей проектов в ответ gRPC
	protoProjects := make([]*projectProto.ProjectResponse, 0, len(projects))
	for _, proj := range projects {
		protoProj, err := projectToResponse(proj)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		protoProjects = append(protoProjects, protoProj)
	}
	
	return &projectProto.ListOfProjectsResponse{
		Projects: protoProjects,
	}, nil
}

func (s *ProjectServer) InitProject(
	ctx context.Context,
	in *projectProto.InitProjectRequest,
) (*projectProto.ProjectResponse, error) {
	if in.ComposeId == nil {
		return nil, status.Error(codes.InvalidArgument, "не указан идентификатор проекта")
	}
	
	project, err := s.projectService.InitProject(ctx, in.ComposeId.Owner, in.ComposeId.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	return projectToResponse(project)
}

func (s *ProjectServer) UpdateProject(
	ctx context.Context,
	in *projectProto.UpdateProjectRequest,
) (*projectProto.ProjectResponse, error) {
	if in.ComposeId == nil {
		return nil, status.Error(codes.InvalidArgument, "не указан идентификатор проекта")
	}
	
	if in.Data == nil {
		return nil, status.Error(codes.InvalidArgument, "не указаны данные для обновления")
	}
	
	// Преобразуем structpb.Struct в map[string]interface{}
	data := in.Data.AsMap()
	
	project, err := s.projectService.UpdateProject(ctx, in.ComposeId.Owner, in.ComposeId.Name, data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	return projectToResponse(project)
}

func (s *ProjectServer) WatchProjectStatus(
	req *projectProto.ProjectUniqueIdentifier,
	stream projectProto.ProjectService_WatchProjectStatusServer,
) error {
	// TODO: Реализовать отслеживание статуса проекта
	// Этот метод должен устанавливать постоянное соединение и отправлять обновления о статусе проекта
	return status.Error(codes.Unimplemented, "метод еще не реализован")
}

// Вспомогательная функция для преобразования модели проекта в ответ gRPC
func projectToResponse(project *models.Project) (*projectProto.ProjectResponse, error) {
	if project == nil {
		return nil, status.Error(codes.NotFound, "проект не найден")
	}
	
	var data map[string]interface{}
	if project.Data != nil {
		err := json.Unmarshal(project.Data, &data)
		if err != nil {
			return nil, status.Error(codes.Internal, "ошибка при парсинге данных проекта")
		}
	} else {
		data = make(map[string]interface{})
	}
	
	// Преобразуем map[string]interface{} в structpb.Struct
	dataStruct, err := structpb.NewStruct(data)
	if err != nil {
		return nil, status.Error(codes.Internal, "ошибка при преобразовании данных проекта")
	}
	
	return &projectProto.ProjectResponse{
		ComposeId: &projectProto.ProjectUniqueIdentifier{
			Owner: project.Owner,
			Name:  project.Name,
		},
		Data: dataStruct,
	}, nil
}
