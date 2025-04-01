package projectserver

import (
	"context"
	projectProto "github.com/SmartAPIForge/protos/gen/go/project"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project-service/internal/domain/models"
	projectservice "project-service/internal/services/project"
	"strconv"
)

type ProjectService interface {
	GetAllUserProjects(
		ctx context.Context,
		owner string,
		page, limit int64,
	) ([]*models.Project, error)
	GetFilteredProjects(
		ctx context.Context,
		owner, status, namePrefix string,
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
	DeleteProject(
		ctx context.Context,
		owner string,
		name string,
	) error
}

type ProjectServer struct {
	projectProto.UnsafeProjectServiceServer
	projectService ProjectService
	projectUpdater *projectservice.ProjectUpdater
}

func RegisterProjectServer(
	gRPCServer *grpc.Server,
	project ProjectService,
	projectUpdater *projectservice.ProjectUpdater,
) {
	projectProto.RegisterProjectServiceServer(
		gRPCServer,
		&ProjectServer{projectService: project, projectUpdater: projectUpdater},
	)
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

func (s *ProjectServer) GetFilteredProjects(
	ctx context.Context,
	in *projectProto.GetFilteredProjectsRequest,
) (*projectProto.ListOfProjectsResponse, error) {
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

	projects, err := s.projectService.GetFilteredProjects(
		ctx,
		in.Owner,
		in.Status,
		in.NamePrefix,
		page,
		limit,
	)
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

func (s *ProjectServer) StreamUserProjectsUpdates(
	req *projectProto.Owner,
	stream projectProto.ProjectService_StreamUserProjectsUpdatesServer,
) error {
	updates := s.projectUpdater.Subscribe()

	for {
		select {
		case project, ok := <-updates:
			if !ok {
				return status.Error(codes.Internal, "project updates channel closed")
			}

			if project.Owner != req.Owner {
				continue
			}

			projectResponse, _ := projectToResponse(project)
			if err := stream.Send(projectResponse); err != nil {
				return err
			}

		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
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

func (s *ProjectServer) DeleteProject(
	ctx context.Context,
	in *projectProto.ProjectUniqueIdentifier,
) (*projectProto.DeleteProjectResponse, error) {
	if in.Owner == "" || in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "не указаны владелец или имя проекта")
	}

	err := s.projectService.DeleteProject(ctx, in.Owner, in.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &projectProto.DeleteProjectResponse{
		Success: true,
	}, nil
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
		Info: &projectProto.ProjectInfo{
			Data:      project.Data,
			Status:    project.Status,
			UrlZip:    project.UrlZip,
			UrlDeploy: project.UrlDeploy,
			CreatedAt: project.CreatedAt.Time().Unix(),
			UpdatedAt: project.UpdatedAt.Time().Unix(),
		},
	}, nil
}
