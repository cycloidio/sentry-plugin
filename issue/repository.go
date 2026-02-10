package issue

import "context"

//go:generate go tool mockgen -destination=../mock/issue_repository.go -mock_names=Repository=IssueRepository -package mock github.com/cycloidio/sentry-plugin/issue Repository
type Repository interface {
	Create(ctx context.Context, orgSlug, prjSlug string, i Issue) (uint32, error)
	//Update(ctx context.Context, orgSlug, prjSlug, isuID string, i Issue) error
	//Find(ctx context.Context, orgSlug, prjSlug, isuID string) (*Issue, error)
	//Filter(ctx context.Context, orgSlug, prjSlug string) ([]*Issue, error)
	//Delete(ctx context.Context, orgSlug, prjSlug, isuID string) error
}
