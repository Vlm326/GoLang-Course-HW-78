package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/vlm326/golang-course-hw-78/internal/collector/domain"
)

type providerStub struct {
	result domain.Repository
	err    error
	owner  string
	repo   string
}

func (p *providerStub) GetRepository(_ context.Context, owner, repo string) (domain.Repository, error) {
	p.owner = owner
	p.repo = repo

	return p.result, p.err
}

func TestGetRepositoryUseCaseExecuteRejectsEmptyInput(t *testing.T) {
	t.Parallel()

	uc := NewGetRepositoryUseCase(&providerStub{})

	_, err := uc.Execute(context.Background(), " ", "")
	if !errors.Is(err, domain.ErrInvalidRepository) {
		t.Fatalf("expected invalid repository error, got %v", err)
	}
}

func TestGetRepositoryUseCaseExecuteTrimsInput(t *testing.T) {
	t.Parallel()

	provider := &providerStub{
		result: domain.Repository{Name: "go"},
	}
	uc := NewGetRepositoryUseCase(provider)

	result, err := uc.Execute(context.Background(), " golang ", " go ")
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if provider.owner != "golang" || provider.repo != "go" {
		t.Fatalf("expected trimmed input, got owner=%q repo=%q", provider.owner, provider.repo)
	}

	if result.Name != "go" {
		t.Fatalf("unexpected repository result: %+v", result)
	}
}
