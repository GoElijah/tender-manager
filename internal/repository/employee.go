package repository

import (
	"context"
	_ "github.com/lib/pq"
	"tender-manager/internal/entity"
)

func (p *PostgresDB) GetByUsername(ctx context.Context, username string) (entity.Employee, error) {
	stmt, err := p.db.PrepareContext(ctx, `SELECT id, organization_id FROM employee WHERE username=$1`)
	if err != nil {
		return entity.Employee{}, err
	}

	var id, organizationID string
	if err := stmt.QueryRow(username).Scan(&id, &organizationID); err != nil {
		return entity.Employee{}, err
	}
	return entity.Employee{
		ID:             id,
		Username:       username,
		OrganizationID: organizationID,
	}, nil
}
