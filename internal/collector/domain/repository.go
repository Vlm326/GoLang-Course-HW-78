package domain

import "context"

type Repository struct {
	Name        string
	Description string
	Stars       int64
	Forks       int64
	CreatedAt   string
}

type RepositoryProvider interface {
	GetRepository(ctx context.Context, owner, repo string) (Repository, error)
}
