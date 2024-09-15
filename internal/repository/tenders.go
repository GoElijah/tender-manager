package repository

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
	"time"
)

func (p *PostgresDB) GetTender(ctx context.Context, tenderID string) (entity.Tender, error) {
	stmt, err := p.db.PrepareContext(ctx, `SELECT id,name,organization,status,version,description FROM tenders WHERE id=$1`)
	if err != nil {
		return entity.Tender{}, err
	}

	var id, name, organization, status, version, description string

	if err := stmt.QueryRow(tenderID).Scan(&id, &name, &organization, &status, &version, &description); err != nil {
		return entity.Tender{}, err
	}

	return entity.Tender{
		ID:           id,
		Name:         name,
		Organization: organization,
		Status:       status,
		Version:      version,
		Description:  description,
	}, nil
}
func (p *PostgresDB) CreateTender(ctx context.Context, t entity.Tender) (string, error) {
	var tenderId string
	stmt, err := p.db.PrepareContext(ctx, `INSERT INTO tenders
													(name, description, organization, service_type, created_by, status, version)
													VALUES($1, $2, $3, $4, $5, 'CREATED', 1) RETURNING id`)
	if err != nil {
		return tenderId, err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(t.Name, t.Description, t.Organization, t.ServiceType, t.CreatedBy).Scan(&tenderId); err != nil {
		return tenderId, err
	}

	return tenderId, nil
}
func (p *PostgresDB) PublishTender(ctx context.Context, t entity.Tender) (string, error) {
	stmt, err := p.db.PrepareContext(ctx, `UPDATE tenders
												SET Status='PUBLISHED',
												updated_by=$1,
												updated_at=$2
												WHERE id=$3
												RETURNING id`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	tender, err := p.GetTender(ctx, t.ID)
	if err != nil {
		return "", err
	}

	if tender.Status != "CREATED" {
		return "", fmt.Errorf("status should be CREATED")
	}

	var tenderId string
	if err := stmt.QueryRow(t.UpdatedBy, time.Now(), t.ID).Scan(&tenderId); err != nil {
		return "", err
	}
	return tenderId, nil
}
func (p *PostgresDB) PatchTender(ctx context.Context, t entity.Tender) (entity.Tender, error) {
	var tender entity.Tender
	previousTenderConditions, _ := p.GetTender(ctx, t.ID)
	err := p.SaveSnapshot(ctx, previousTenderConditions)
	if err != nil {
		return tender, err
	}

	stmt, err := p.db.PrepareContext(ctx, `UPDATE tenders
												SET name=$1, description=$2, updated_by=$3, updated_at=$4,
												version=version+1 WHERE id=$5 RETURNING id,name,description,version`)
	if err != nil {
		return tender, err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(t.Name, t.Description, t.UpdatedBy, time.Now(), t.ID).Scan(&tender.ID, &tender.Name, &tender.Description, &tender.Version); err != nil {
		return tender, err
	}

	return tender, nil
}
func (p *PostgresDB) CloseTender(ctx context.Context, t entity.Tender) (string, error) {
	stmt, err := p.db.PrepareContext(ctx, `UPDATE tenders
												SET Status='CLOSED',
												updated_by=$1,
												updated_at=$2
												WHERE id=$3`)

	if err != nil {
		return "", err
	}
	defer stmt.Close()
	tender, err := p.GetTender(ctx, t.ID)
	if err != nil {
		return "", err
	}

	if tender.Status == "CLOSED" {
		return "", fmt.Errorf("tender has already CLOSED")
	}

	var tenderId string
	if err := stmt.QueryRow(t.UpdatedBy, time.Now(), t.ID).Scan(&tenderId); err != nil {
		return "", err
	}
	return tenderId, nil
}
func (p *PostgresDB) ListTenders(ctx context.Context, t entity.Tender) ([]gen.Tender, error) {
	var list []gen.Tender

	var args []interface{}
	args = append(args, t.Organization)

	query := `SELECT id,name,description,service_type FROM tenders
						WHERE (tenders.organization = $1 OR status = 'PUBLISHED')`
	if t.ServiceType != "" {
		query += " AND service_type = $2"
		args = append(args, t.ServiceType)
	}

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t gen.Tender
		err := rows.Scan(&t.Id, &t.Name, &t.Description, &t.ServiceType)
		if err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, nil
}
func (p *PostgresDB) ListMyTenders(ctx context.Context, t entity.Tender) ([]gen.Tender, error) {
	var list []gen.Tender

	stmt, err := p.db.PrepareContext(ctx, `SELECT id,name,description,service_type FROM tenders
						WHERE tenders.created_by = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, t.CreatedBy)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t gen.Tender
		err := rows.Scan(&t.Id, &t.Name, &t.Description, &t.ServiceType)
		if err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, nil
}
func (p *PostgresDB) RollbackTender(ctx context.Context, t entity.Tender) (entity.Tender, error) {
	var patchedTender entity.Tender
	stmt, err := p.db.PrepareContext(ctx, `SELECT name,description,version
												FROM tenders_snapshots
												WHERE tender_id=$1 AND version=$2`)
	if err != nil {
		return patchedTender, err
	}
	defer stmt.Close()
	var name, description, version string
	if err = stmt.QueryRow(t.ID, t.Version).Scan(&name, &description, &version); err != nil {
		return patchedTender, err
	}

	snapshot := entity.Tender{
		ID:          t.ID,
		Name:        name,
		Description: description,
		UpdatedBy:   t.UpdatedBy,
	}

	patchedTender, err = p.PatchTender(ctx, snapshot)
	if err != nil {
		return patchedTender, err
	}

	return patchedTender, nil
}
func (p *PostgresDB) SaveSnapshot(ctx context.Context, t entity.Tender) error {
	stmt, err := p.db.PrepareContext(ctx, `INSERT INTO tenders_snapshots
													(name, description, version, tender_id)
													VALUES($1, $2, $3, $4) `)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(t.Name, t.Description, t.Version, t.ID); err != nil {
		return err
	}

	return nil

}
func (p *PostgresDB) Close() error {
	return p.db.Close()
}
