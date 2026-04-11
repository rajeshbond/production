package productoperationsequence

import (
	"context"
)

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (s *Service) Create(
	ctx context.Context,
	tenantID, userID int64,
	req CreateProductOperationRequest,
) (*ProductOperationResponse, error) {

	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 🔐 lock rows (important)
	if err := s.Store.LockRows(ctx, tx, tenantID, req.ProductID); err != nil {
		return nil, err
	}

	// 🔢 get next sequence
	seq, err := s.Store.GetNextSequenceNo(ctx, tx, tenantID, req.ProductID)
	if err != nil {
		return nil, err
	}

	p := &ProductOperationSequence{
		TenantID:    tenantID,
		ProductID:   req.ProductID,
		OperationID: req.OperationID,
		SequenceNo:  seq,
		CreatedBy:   &userID,
		UpdatedBy:   &userID,
	}

	// ➕ insert
	if err := s.Store.Create(ctx, tx, p); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// ✅ build response
	return &ProductOperationResponse{
		ID:          p.ID,
		ProductID:   p.ProductID,
		OperationID: p.OperationID,
		SequenceNo:  p.SequenceNo,
	}, nil
}

//

func (s *Service) BulkCreate(
	ctx context.Context,
	tenantID, userID int64,
	req BulkCreateProductOperationRequest,
) ([]ProductOperationResponse, error) {

	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 🔐 Lock rows (avoid race condition)
	if err := s.Store.LockRows(ctx, tx, tenantID, req.ProductID); err != nil {
		return nil, err
	}

	// 🔢 Get current max sequence
	maxSeq, err := s.Store.GetCurrentMaxSequence(ctx, tx, tenantID, req.ProductID)
	if err != nil {
		return nil, err
	}

	var records []ProductOperationSequence

	// 🧠 Build records with sequence
	for i, opID := range req.OperationIDs {
		records = append(records, ProductOperationSequence{
			TenantID:    tenantID,
			ProductID:   req.ProductID,
			OperationID: opID,
			SequenceNo:  maxSeq + i + 1,
			CreatedBy:   &userID,
			UpdatedBy:   &userID,
		})
	}

	// ➕ Insert into DB
	if err := s.Store.BulkCreate(ctx, tx, records); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 🎯 Convert to response
	var resp []ProductOperationResponse
	for _, r := range records {
		resp = append(resp, ProductOperationResponse{
			ID:          r.ID,
			ProductID:   r.ProductID,
			OperationID: r.OperationID,
			SequenceNo:  r.SequenceNo,
		})
	}

	return resp, nil
}
