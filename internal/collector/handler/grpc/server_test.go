package grpc

import (
	"context"
	"testing"

	"github.com/vlm326/golang-course-hw-78/internal/collector/domain"
	collectorusecase "github.com/vlm326/golang-course-hw-78/internal/collector/usecase"
	"github.com/vlm326/golang-course-hw-78/internal/shared/repositoryrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type providerStub struct {
	result domain.Repository
	err    error
}

func (p *providerStub) GetRepository(_ context.Context, _, _ string) (domain.Repository, error) {
	return p.result, p.err
}

func TestServerGetRepositoryReturnsMappedResponse(t *testing.T) {
	t.Parallel()

	server := NewServer(collectorusecase.NewGetRepositoryUseCase(&providerStub{
		result: domain.Repository{
			Name:        "go",
			Description: "The Go programming language",
			Stars:       1,
			Forks:       2,
			CreatedAt:   "2009-11-10T23:00:00Z",
		},
	}))

	resp, err := server.GetRepository(context.Background(), &repositoryrpc.GetRepositoryRequest{
		Owner: "golang",
		Repo:  "go",
	})
	if err != nil {
		t.Fatalf("GetRepository returned error: %v", err)
	}

	if resp.Name != "go" || resp.Stars != 1 || resp.Forks != 2 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestServerGetRepositoryMapsNotFoundToGRPCStatus(t *testing.T) {
	t.Parallel()

	server := NewServer(collectorusecase.NewGetRepositoryUseCase(&providerStub{
		err: domain.ErrRepositoryNotFound,
	}))

	_, err := server.GetRepository(context.Background(), &repositoryrpc.GetRepositoryRequest{
		Owner: "golang",
		Repo:  "missing",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	if status.Code(err) != codes.NotFound {
		t.Fatalf("expected NotFound, got %s", status.Code(err))
	}
}
