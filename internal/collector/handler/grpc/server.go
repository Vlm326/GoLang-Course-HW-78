package grpc

import (
	"context"
	"errors"

	"github.com/vlm326/golang-course-hw-78/internal/collector/domain"
	"github.com/vlm326/golang-course-hw-78/internal/collector/usecase"
	"github.com/vlm326/golang-course-hw-78/internal/shared/repositoryrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	getRepository usecase.GetRepositoryUseCase
}

func NewServer(getRepository usecase.GetRepositoryUseCase) *Server {
	return &Server{getRepository: getRepository}
}

func (s *Server) GetRepository(ctx context.Context, req *repositoryrpc.GetRepositoryRequest) (*repositoryrpc.RepositoryResponse, error) {
	repository, err := s.getRepository.Execute(ctx, req.Owner, req.Repo)
	if err != nil {
		return nil, mapError(err)
	}

	return &repositoryrpc.RepositoryResponse{
		Name:        repository.Name,
		Description: repository.Description,
		Stars:       repository.Stars,
		Forks:       repository.Forks,
		CreatedAt:   repository.CreatedAt,
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidRepository):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrRepositoryNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
