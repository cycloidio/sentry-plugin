package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cycloidio/sentry-plugin/issue"
	"github.com/cycloidio/sqlr"
)

type IssueRepository struct {
	querier sqlr.Querier
}

func NewIssueRepository(db sqlr.Querier) *IssueRepository {
	return &IssueRepository{
		querier: db,
	}
}

type dbIssue struct {
	ID        sql.NullString
	Title     sql.NullString
	Permalink sql.NullString
	HasSeen   sql.NullBool
	FirstSeen sql.NullTime
	LastSeen  sql.NullTime
	UserCount sql.NullInt64
	Level     sql.NullString
	Status    sql.NullString
	Type      sql.NullString
}

func newDBIssue(i issue.Issue) dbIssue {
	return dbIssue{
		ID:        toNullString(i.ID),
		Title:     toNullString(i.Title),
		Permalink: toNullString(i.Permalink),
		HasSeen:   toNullBool(i.HasSeen),
		FirstSeen: toNullTime(i.FirstSeen),
		LastSeen:  toNullTime(i.LastSeen),
		UserCount: toNullInt64(i.UserCount),
		Level:     toNullString(i.Level),
		Status:    toNullString(i.Status),
		Type:      toNullString(i.Type),
	}
}

func (r *IssueRepository) Create(ctx context.Context, orgSlug, prjSlug string, i issue.Issue) (uint32, error) {
	dbi := newDBIssue(i)
	_, err := r.querier.ExecContext(ctx, `
		INSERT INTO issues(id, title, permalink, has_seen, first_seen, last_seen, user_count, level, status, type, project_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			-- project_id
			(
				SELECT p.id
				FROM projects AS p
				JOIN organizations AS o
					ON p.organization_id = o.id
				WHERE o.slug = ? AND p.slug = ?
			))
	`, dbi.ID, dbi.Title, dbi.Permalink, dbi.HasSeen, dbi.FirstSeen, dbi.LastSeen, dbi.UserCount, dbi.Level, dbi.Status, dbi.Type, orgSlug, prjSlug)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return 0, nil
}
