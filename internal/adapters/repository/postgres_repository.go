package repository

import (
	"context"
	"errors"
	"meme-bot/internal/domain"
	"meme-bot/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) ports.Repository {
	return &PostgresRepository{db: db}
}
func (r *PostgresRepository) GetByID(id int) (domain.Meme, error) {
	query := `
		SELECT id, phash, status, source, source_id, created_at
		FROM memes
		WHERE id = $1
	`

	var m domain.Meme

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&m.ID,
		&m.PHash,
		&m.Status,
		&m.Source,
		&m.SourceID,
		&m.CreatedAt,
	)

	return m, err
}

func (r *PostgresRepository) Save(m *domain.Meme) error {
	query := `
		INSERT INTO memes (phash, status, source, source_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (phash) DO NOTHING
		RETURNING id, created_at
	`

	err := r.db.QueryRow(context.Background(),
		query,
		m.PHash,
		m.Status,
		m.Source,
		m.SourceID,
	).Scan(&m.ID, &m.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) ExistsByHash(hash string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM memes WHERE phash = $1)`

	err := r.db.QueryRow(context.Background(), query, hash).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) ExistsBySourceID(sourceID string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM memes WHERE source_id = $1)`

	err := r.db.QueryRow(context.Background(), query, sourceID).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) GetByStatus(status domain.MemeStatus) ([]domain.Meme, error) {
	query := `
		SELECT id, phash, status, source, source_id, created_at
		FROM memes
		WHERE status = $1
		ORDER BY created_at
	`

	rows, err := r.db.Query(context.Background(), query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memes []domain.Meme

	for rows.Next() {
		var m domain.Meme

		if err := rows.Scan(
			&m.ID,
			&m.PHash,
			&m.Status,
			&m.Source,
			&m.SourceID,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}

		memes = append(memes, m)
	}

	return memes, rows.Err()
}

func (r *PostgresRepository) UpdateStatus(id int, status domain.MemeStatus) error {
	query := `
		UPDATE memes
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(context.Background(), query, status, id)
	return err
}

func (r *PostgresRepository) GetOldestApproved() (domain.Meme, error) {
	query := `
		SELECT id, phash, status, source, source_id, created_at
		FROM memes
		WHERE status = $1
		ORDER BY created_at
		LIMIT 1
	`

	var m domain.Meme

	if err := r.db.QueryRow(
		context.Background(),
		query,
		domain.Approved,
	).Scan(
		&m.ID,
		&m.PHash,
		&m.Status,
		&m.Source,
		&m.SourceID,
		&m.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Meme{}, ports.ErrMemesEnded
		}
		return domain.Meme{}, err
	}

	return m, nil
}

func (r *PostgresRepository) GetOldestPending() (domain.Meme, error) {
	query := `
		SELECT id, phash, status, source, source_id, created_at
		FROM memes
		WHERE status = $1
		ORDER BY created_at
		LIMIT 1
	`

	var m domain.Meme

	if err := r.db.QueryRow(
		context.Background(),
		query,
		domain.Pending,
	).Scan(
		&m.ID,
		&m.PHash,
		&m.Status,
		&m.Source,
		&m.SourceID,
		&m.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Meme{}, ports.ErrMemesEnded
		}
		return domain.Meme{}, err
	}

	return m, nil
}

func (r *PostgresRepository) Delete(ID int) error {
	query := `
		DELETE FROM memes
		WHERE id = $1
	`

	_, err := r.db.Exec(context.Background(), query, ID)
	return err
}
