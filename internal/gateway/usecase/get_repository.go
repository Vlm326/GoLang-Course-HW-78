package usecase

import (
	"context"

	"github.com/vlm326/golang-course-hw-78/internal/shared/repositoryrpc"
)

type GetRepositoryUseCase struct {
	client repositoryrpc.RepositoryServiceClient
}

func NewGetRepositoryUseCase(client repositoryrpc.RepositoryServiceClient) GetRepositoryUseCase {
	return GetRepositoryUseCase{client: client}
}

func (uc GetRepositoryUseCase) Execute(ctx context.Context, owner, repo string) (*repositoryrpc.RepositoryResponse, error) {
	return uc.client.GetRepository(ctx, &repositoryrpc.GetRepositoryRequest{
		Owner: owner,
		Repo:  repo,
	})
}
