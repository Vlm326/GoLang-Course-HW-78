package usecase

import (
	"context"
	"strings"

	"github.com/vlm326/golang-course-hw-78/internal/collector/domain"
)

type GetRepositoryUseCase struct {
	provider domain.RepositoryProvider
}

func NewGetRepositoryUseCase(provider domain.RepositoryProvider) GetRepositoryUseCase {
	return GetRepositoryUseCase{provider: provider}
}

func (uc GetRepositoryUseCase) Execute(ctx context.Context, owner, repo string) (domain.Repository, error) {
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)

	if owner == "" || repo == "" {
		return domain.Repository{}, domain.ErrInvalidRepository
	}

	return uc.provider.GetRepository(ctx, owner, repo)
}
