package shiftslot

import "context"

type Service struct {
	Store *Store
}

func NewService(store *Store) *Service {
	return &Service{Store: store}
}

func (s *Service) GenerateAllShiftSlots(ctx context.Context, tenantID int64) error {

	shiftTimings, err := s.Store.GetShiftTimingsByTenant(ctx, tenantID)
	if err != nil {
		return err
	}

	tx, err := s.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, st := range shiftTimings {

		// check already exists
		exists, err := s.Store.SlotsExist(ctx, tx, st.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if exists {
			continue // ✅ skip existing
		}

		slots := GenerateSlots(st.Start, st.End)

		err = s.Store.InsertSlots(ctx, tx, tenantID, st.ID, slots)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
