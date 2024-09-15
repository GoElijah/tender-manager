package repository

import (
	"context"
	"errors"
	"fmt"
	"tender-manager/internal/entity"
	gen "tender-manager/internal/generated"
	"time"
)

func (p *PostgresDB) CreateBid(ctx context.Context, b entity.Bid) (string, error) {
	var bidId string
	stmt, err := p.db.PrepareContext(ctx, `INSERT INTO bids(name, description, tender_id, created_by, tender_organization, bid_organization, status, version, votes)
													VALUES($1,$2,$3,$4,$5,$6,'CREATED',1,0) RETURNING id`)
	if err != nil {
		return bidId, err
	}
	defer stmt.Close()

	if err != nil {
		return bidId, err
	}
	if err := stmt.QueryRow(b.Name, b.Description, b.TenderID, b.CreatedBy, b.TenderOrganization, b.BidOrganization).Scan(&bidId); err != nil {
		return bidId, err
	}
	return bidId, nil
}
func (p *PostgresDB) GetBid(ctx context.Context, bidId string) (entity.Bid, error) {
	stmt, err := p.db.PrepareContext(ctx, `SELECT id,name,status,version,created_by,tender_organization,bid_organization,description
												FROM bids WHERE id=$1`)
	if err != nil {
		return entity.Bid{}, err
	}

	var bid entity.Bid
	if err := stmt.QueryRow(bidId).Scan(
		&bid.ID, &bid.Name, &bid.Status, &bid.Version, &bid.CreatedBy, &bid.TenderOrganization, &bid.BidOrganization, &bid.Description); err != nil {
		return entity.Bid{}, err
	}

	return bid, nil
}
func (p *PostgresDB) PublishBid(ctx context.Context, b entity.Bid) error {
	stmt, err := p.db.PrepareContext(ctx, `UPDATE bids
												SET Status='PUBLISHED',
												updated_by=$1,
												updated_at=$2
												WHERE id=$3`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	bid, err := p.GetBid(ctx, b.ID)
	if err != nil {
		return err
	}

	if bid.Status != "CREATED" {
		return fmt.Errorf("status should be 'CREATED'")
	}

	if _, err := stmt.Exec(b.UpdatedBy, time.Now(), b.ID); err != nil {
		return err
	}
	return nil
}
func (p *PostgresDB) CancelBid(ctx context.Context, b entity.Bid) error {
	stmt, err := p.db.PrepareContext(ctx, `UPDATE bids
												SET Status='CANCELED',
												updated_by=$1,
												updated_at=$2
												WHERE id=$3`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	bid, err := p.GetBid(ctx, b.ID)
	if err != nil {
		return err
	}

	if bid.Status == "CANCELED" {
		return errors.New("status should be diffrent from CANCELED")
	}

	if _, err := stmt.Exec(b.UpdatedBy, time.Now(), b.ID); err != nil {
		return err
	}
	return nil
}
func (p *PostgresDB) PatchBid(ctx context.Context, b entity.Bid) (entity.Bid, error) {
	var bid entity.Bid
	previousBidConditions, _ := p.GetBid(ctx, b.ID)
	err := p.SaveBidSnapshot(ctx, previousBidConditions)
	if err != nil {
		return bid, err
	}

	stmt, err := p.db.PrepareContext(ctx, `UPDATE bids
												SET name=$1, description=$2, updated_by=$3, updated_at=$4,
												version=version+1 WHERE id=$5 
												RETURNING id,name,description,version`)
	if err != nil {
		return bid, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(b.Name, b.Description, b.UpdatedBy, time.Now(), b.ID).Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Version); err != nil {
		return bid, err
	}

	return bid, nil
}
func (p *PostgresDB) ListBids(ctx context.Context, b entity.Bid) ([]gen.Bid, error) {
	var list []gen.Bid

	query := `SELECT id,name,description,tender_id FROM bids
						WHERE bids.created_by = $1 
						OR (status='PUBLISHED' AND tender_id=$2)`

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, b.CreatedBy, b.TenderID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b gen.Bid
		err := rows.Scan(&b.Id, &b.Name, &b.Description, &b.TenderId)
		if err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}
func (p *PostgresDB) ListMyBids(ctx context.Context, username string) ([]gen.Bid, error) {
	var list []gen.Bid

	stmt, err := p.db.PrepareContext(ctx, `SELECT id,name,description,tender_id FROM bids
						WHERE bids.created_by = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b gen.Bid
		err = rows.Scan(&b.Id, &b.Name, &b.Description, &b.TenderId)
		if err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}
func (p *PostgresDB) RollbackBid(ctx context.Context, t entity.Bid) (entity.Bid, error) {
	var patchedBid entity.Bid
	stmt, err := p.db.PrepareContext(ctx, `SELECT name,description,version
												FROM bids_snapshots
												WHERE bid_id=$1 AND version=$2`)
	if err != nil {
		return patchedBid, err
	}
	defer stmt.Close()
	var name, description, version string
	if err = stmt.QueryRow(t.ID, t.Version).Scan(&name, &description, &version); err != nil {
		return patchedBid, err
	}

	snapshot := entity.Bid{
		ID:          t.ID,
		Name:        name,
		Description: description,
		UpdatedBy:   t.UpdatedBy,
	}

	patchedBid, err = p.PatchBid(ctx, snapshot)
	if err != nil {
		return patchedBid, err
	}

	return patchedBid, nil
}
func (p *PostgresDB) SubmitBid(ctx context.Context, b entity.Bid) (entity.Bid, error) {

	var bid entity.Bid
	stmt, err := p.db.PrepareContext(ctx, `UPDATE bids 
												SET votes=votes+1 WHERE id=$1
													RETURNING id,name,votes,status`)
	if err != nil {
		return bid, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(b.ID).Scan(&bid.ID, &bid.Name, &bid.Votes, &bid.Status); err != nil {
		return bid, err
	}

	return bid, nil
}
func (p *PostgresDB) RejectBid(ctx context.Context, b entity.Bid) (entity.Bid, error) {
	var bid entity.Bid
	err := p.CancelBid(ctx, b)
	if err != nil {
		return bid, err
	}

	stmt, err := p.db.PrepareContext(ctx, `UPDATE bids SET votes=-1 WHERE id=$1
													RETURNING id,name,votes,status`)
	if err != nil {
		return bid, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(b.ID).Scan(&bid.ID, &bid.Name, &bid.Version, &bid.Status, &bid.Votes); err != nil {
		return bid, err
	}

	return bid, nil
}
func (p *PostgresDB) ApproveBid(ctx context.Context, b entity.Bid) (entity.Bid, error) {
	var bid entity.Bid

	stmt, err := p.db.PrepareContext(ctx, `UPDATE bids
												SET status='APPROVED', updated_by=$1, updated_at=$2
												WHERE id=$3 RETURNING id,name,status`)
	if err != nil {
		return bid, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(b.UpdatedBy, time.Now(), b.ID).Scan(&bid.Votes, &bid.Name, &bid.Status); err != nil {
		return bid, err
	}

	return bid, nil

}
func (p *PostgresDB) SaveBidSnapshot(ctx context.Context, b entity.Bid) error {
	stmt, err := p.db.PrepareContext(ctx, `INSERT INTO bids_snapshots
													(name, description, version, bid_id)
													VALUES($1, $2, $3, $4) `)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(b.Name, b.Description, b.Version, b.ID); err != nil {
		return err
	}

	return nil

}
func (p *PostgresDB) Feedback(ctx context.Context, b entity.Bid) error {
	stmt, err := p.db.PrepareContext(ctx, `UPDATE bids
												SET feedback=$1,
												updated_by=$2,
												updated_at=$3
												WHERE id=$4`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(b.Feedback, b.UpdatedBy, time.Now(), b.ID); err != nil {
		return err
	}
	return nil
}
func (p *PostgresDB) ListFeedback(ctx context.Context, b entity.Bid) ([]gen.Feedback, error) {
	var list []gen.Feedback

	query := `SELECT b.id, b.feedback FROM bids b
			INNER JOIN tenders t ON b.tender_id = t.id
        	INNER JOIN organization o ON t.organization = o.id
			WHERE
				t.id = $1 AND b.created_by = $2 AND t.organization = $3`

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, b.TenderID, b.CreatedBy, b.TenderOrganization)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var f gen.Feedback
		err := rows.Scan(&f.BidId, &f.Review)
		if err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, nil
}
