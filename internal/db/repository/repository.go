package repository

import (
	"context"
	"database/sql"
	"time"

	"dppg/internal/db/model"

	"github.com/gofrs/uuid/v5"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, userID uuid.UUID, name, content string) (*model.UserPrompt, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	query := psql.Insert(
		im.Into("prompts", "id", "user_id", "name", "content", "created_at", "updated_at"),
		im.Values(psql.Arg(id), psql.Arg(userID), psql.Arg(name), psql.Arg(content), psql.Arg(now), psql.Arg(now)),
		im.Returning("id", "user_id", "name", "content", "created_at", "updated_at"),
	)

	var prompt model.UserPrompt

	// Bob's Scan is usually via bob.Exec or similar, but here we use stdlib with bob's query string
	// Or we can use bob's built-in execution if we wrap sql.DB?
	// Let's use standard sql.DB with generated SQL

	sqlQuery, args, err := query.Build(ctx)
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, sqlQuery, args...)
	err = row.Scan(&prompt.ID, &prompt.UserID, &prompt.Name, &prompt.Content, &prompt.CreatedAt, &prompt.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &prompt, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.UserPrompt, error) {
	query := psql.Select(
		sm.Columns("id", "user_id", "name", "content", "created_at", "updated_at"),
		sm.From("prompts"),
		sm.Where(psql.Quote("user_id").EQ(psql.Arg(userID))),
		sm.OrderBy("created_at").Desc(),
	)

	sqlQuery, args, err := query.Build(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []model.UserPrompt
	for rows.Next() {
		var p model.UserPrompt
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		prompts = append(prompts, p)
	}

	return prompts, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.UserPrompt, error) {
	query := psql.Select(
		sm.Columns("id", "user_id", "name", "content", "created_at", "updated_at"),
		sm.From("prompts"),
		sm.Where(psql.Quote("id").EQ(psql.Arg(id)).And(psql.Quote("user_id").EQ(psql.Arg(userID)))),
	)

	sqlQuery, args, err := query.Build(ctx)
	if err != nil {
		return nil, err
	}

	var p model.UserPrompt
	row := r.db.QueryRowContext(ctx, sqlQuery, args...)
	if err := row.Scan(&p.ID, &p.UserID, &p.Name, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := psql.Delete(
		dm.From("prompts"),
		dm.Where(psql.Quote("id").EQ(psql.Arg(id)).And(psql.Quote("user_id").EQ(psql.Arg(userID)))),
	)

	sqlQuery, args, err := query.Build(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlQuery, args...)
	return err
}
