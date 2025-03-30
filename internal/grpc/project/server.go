package projectserver

import (
	"context"
	projectProto "github.com/SmartAPIForge/protos/gen/go/project"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project-service/internal/domain/models"
	"strconv"
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
		data string,
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

	project, err := s.projectService.UpdateProject(ctx, in.ComposeId.Owner, in.ComposeId.Name, in.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return projectToResponse(project)
}

func (s *ProjectServer) WatchProjectStatus(
	req *projectProto.ProjectUniqueIdentifier,
	stream projectProto.ProjectService_WatchProjectStatusServer,
) error {
	// TODO
	return status.Error(codes.Unimplemented, "метод еще не реализован")
}

func projectToResponse(project *models.Project) (*projectProto.ProjectResponse, error) {
	if project == nil {
		return nil, status.Error(codes.NotFound, "проект не найден")
	}

	return &projectProto.ProjectResponse{
		ComposeId: &projectProto.ProjectUniqueIdentifier{
			Owner: project.Owner,
			Name:  project.Name,
		},
		Data: project.Data,
	}, nil
}
