package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

const URLSTable = "urls"

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Save(ctx context.Context, urlEntity model.URL) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, original, short_code)
				VALUES ($1, $2, $3)`, URLSTable)
	_, err := s.pool.Exec(ctx, query, urlEntity.ID, urlEntity.Original, urlEntity.ShortCode)
	return err
}

func (s *Store) FindByShortCode(ctx context.Context, shortCode string) (model.URL, error) {
	query := fmt.Sprintf(`SELECT id, original, created_at FROM %s WHERE short_code=$1`, URLSTable)
	row := s.pool.QueryRow(ctx, query, shortCode)

	var urlEntity model.URL
	if err := row.Scan(&urlEntity.ID, &urlEntity.Original, &urlEntity.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.URL{}, errors.New("no such record found")
		}
		return model.URL{}, fmt.Errorf("failed to fetch data for %s: %w", shortCode, err)
	}
	urlEntity.ShortCode = shortCode

	return urlEntity, nil
}
