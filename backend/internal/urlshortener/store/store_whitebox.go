package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Store) DeleteUrlEntityByID(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", URLSTable)
	_, err := s.pool.Exec(ctx, query, id)

	return err
}
