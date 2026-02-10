package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cycloidio/sentry-plugin/organization"
	"github.com/cycloidio/sqlr"
)

type OrganizationRepository struct {
	querier sqlr.Querier
}

func NewOrganizationRepository(db sqlr.Querier) *OrganizationRepository {
	return &OrganizationRepository{
		querier: db,
	}
}

type dbOrganization struct {
	ID   sql.NullString
	Name sql.NullString
	Slug sql.NullString
}

func newDBOrganization(o organization.Organization) dbOrganization {
	return dbOrganization{
		ID:   toNullString(o.ID),
		Name: toNullString(o.Name),
		Slug: toNullString(o.Slug),
	}
}

func (r *OrganizationRepository) Create(ctx context.Context, o organization.Organization) (uint32, error) {
	dbo := newDBOrganization(o)
	_, err := r.querier.ExecContext(ctx, `
		INSERT INTO organizations(id, name, slug)
		VALUES (?, ?, ?)
	`, dbo.ID, dbo.Name, dbo.Slug)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return 0, nil
}

func (r *OrganizationRepository) DeleteAll(ctx context.Context) error {
	_, err := r.querier.ExecContext(ctx, `
		DELETE 
		FROM organizations
	`)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
