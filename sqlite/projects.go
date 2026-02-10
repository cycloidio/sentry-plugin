package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cycloidio/sentry-plugin/project"
	"github.com/cycloidio/sqlr"
)

type ProjectRepository struct {
	querier sqlr.Querier
}

func NewProjectRepository(db sqlr.Querier) *ProjectRepository {
	return &ProjectRepository{
		querier: db,
	}
}

type dbProject struct {
	ID     sql.NullString
	Name   sql.NullString
	Slug   sql.NullString
	Status sql.NullString
}

func newDBProject(p project.Project) dbProject {
	return dbProject{
		ID:     toNullString(p.ID),
		Name:   toNullString(p.Name),
		Slug:   toNullString(p.Slug),
		Status: toNullString(p.Status),
	}
}

func (r *ProjectRepository) Create(ctx context.Context, orgSlug string, p project.Project) (uint32, error) {
	dbp := newDBProject(p)
	_, err := r.querier.ExecContext(ctx, `
		INSERT INTO projects(id, name, slug, status, organization_id)
		VALUES (?, ?, ?, ?, 
			-- organization_id
			(
				SELECT o.id
				FROM organizations AS o
				WHERE o.slug = ?
			))
	`, dbp.ID, dbp.Name, dbp.Slug, dbp.Status, orgSlug)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return 0, nil
}
