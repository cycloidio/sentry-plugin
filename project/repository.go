package project

import "context"

//go:generate go tool mockgen -destination=../mock/project_repository.go -mock_names=Repository=ProjectRepository -package mock github.com/cycloidio/sentry-plugin/project Repository
type Repository interface {
	Create(ctx context.Context, orgSlug string, p Project) (uint32, error)
	//Update(ctx context.Context, orgSlug, prjSlug string, p Project) error
	//Find(ctx context.Context, orgSlug, prjSlug string) (*Project, error)
	//Filter(ctx context.Context, orgSlug string) ([]*Project, error)
	//Delete(ctx context.Context, orgSlug, prjSlug string) error
}
