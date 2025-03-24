package projectserver

import (
	"context"
	projectProto "github.com/SmartAPIForge/protos/gen/go/project"
	"google.golang.org/grpc"
	"project-service/internal/domain/models"
)

type ProjectService interface {
	GetUniqueUserProject(
		owner string,
		name string,
	) (*models.Project, error)
	GetAllUserProjects(
		owner string,
	) ([]*models.Project, error)
	CreateNewProject(
		owner string,
		name string,
	) (*models.Project, error)
	UpdateProject(
		owner string,
		name string,
		data map[string]interface{},
	) (*models.Project, error)
	GetProjectStatus(
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
	return &projectProto.ProjectResponse{}, nil
}

func (s *ProjectServer) GetAllUserProjects(
	ctx context.Context,
	in *projectProto.GetAllUserProjectsRequest,
) (*projectProto.ListOfProjectsResponse, error) {
	return &projectProto.ListOfProjectsResponse{}, nil
}

func (s *ProjectServer) InitProject(
	ctx context.Context,
	in *projectProto.InitProjectRequest,
) (*projectProto.ProjectResponse, error) {
	return &projectProto.ProjectResponse{}, nil
}

func (s *ProjectServer) UpdateProject(
	ctx context.Context,
	in *projectProto.UpdateProjectRequest,
) (*projectProto.ProjectResponse, error) {
	return &projectProto.ProjectResponse{}, nil
}

func (s *ProjectServer) WatchProjectStatus(
	req *projectProto.ProjectUniqueIdentifier,
	stream projectProto.ProjectService_WatchProjectStatusServer,
) error {
	return nil
}
