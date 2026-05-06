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

func (r *IssueRepository) ListAll(ctx context.Context) ([]issue.IssueWithRelations, error) {
	rows, err := r.querier.QueryContext(ctx, `
		SELECT
			i.id, i.title, i.permalink, i.has_seen, i.first_seen, i.last_seen, i.user_count, i.level, i.status, i.type,
			p.name, p.slug,
			o.name, o.slug
		FROM issues AS i
		JOIN projects AS p ON p.id = i.project_id
		JOIN organizations AS o ON o.id = p.organization_id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query issues: %w", err)
	}
	defer rows.Close()

	var results []issue.IssueWithRelations
	for rows.Next() {
		var dbi dbIssue
		var projectName, projectSlug, orgName, orgSlug sql.NullString
		err := rows.Scan(
			&dbi.ID, &dbi.Title, &dbi.Permalink, &dbi.HasSeen, &dbi.FirstSeen, &dbi.LastSeen, &dbi.UserCount, &dbi.Level, &dbi.Status, &dbi.Type,
			&projectName, &projectSlug,
			&orgName, &orgSlug,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan issue row: %w", err)
		}
		results = append(results, issue.IssueWithRelations{
			Issue: issue.Issue{
				ID:        dbi.ID.String,
				Title:     dbi.Title.String,
				Permalink: dbi.Permalink.String,
				HasSeen:   dbi.HasSeen.Bool,
				FirstSeen: dbi.FirstSeen.Time,
				LastSeen:  dbi.LastSeen.Time,
				UserCount: int(dbi.UserCount.Int64),
				Level:     dbi.Level.String,
				Status:    dbi.Status.String,
				Type:      dbi.Type.String,
			},
			ProjectName:      projectName.String,
			ProjectSlug:      projectSlug.String,
			OrganizationName: orgName.String,
			OrganizationSlug: orgSlug.String,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterating issue rows: %w", err)
	}

	return results, nil
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
