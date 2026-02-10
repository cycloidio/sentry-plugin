package organization

import "context"

//go:generate go tool mockgen -destination=../mock/organization_repository.go -mock_names=Repository=OrganizationRepository -package mock github.com/cycloidio/sentry-plugin/organization Repository
type Repository interface {
	Create(ctx context.Context, o Organization) (uint32, error)
	//Update(ctx context.Context, orgSlug string, o Organization) error
	//Find(ctx context.Context, orgSlug string) (*Organization, error)
	//Filter(ctx context.Context) ([]*Organization, error)
	//Delete(ctx context.Context, orgSlug string) error
	DeleteAll(ctx context.Context) error
}
