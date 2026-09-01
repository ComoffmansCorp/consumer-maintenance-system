package marketplace

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/marketplace/db"
	platformdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{queries: queries}
}

func (r *Repository) queriesFor(ctx context.Context) *db.Queries {
	if tx, ok := platformdb.TxFromContext(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

// --- categories & offerings (admin-curated catalog) ---

func (r *Repository) ListActiveCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.queriesFor(ctx).ListActiveCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active categories: %w", err)
	}
	out := make([]Category, 0, len(rows))
	for _, row := range rows {
		out = append(out, toCategory(row))
	}
	return out, nil
}

func (r *Repository) GetCategoryByID(ctx context.Context, id int64) (Category, error) {
	row, err := r.queriesFor(ctx).GetCategoryByID(ctx, id)
	if err != nil {
		return Category{}, fmt.Errorf("get category by id: %w", err)
	}
	return toCategory(row), nil
}

func (r *Repository) CreateCategory(ctx context.Context, name string) (Category, error) {
	row, err := r.queriesFor(ctx).CreateCategory(ctx, name)
	if err != nil {
		return Category{}, fmt.Errorf("create category: %w", err)
	}
	return toCategory(row), nil
}

func (r *Repository) ListActiveServices(ctx context.Context, categoryID int64) ([]Offering, error) {
	rows, err := r.queriesFor(ctx).ListActiveServices(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("list active services: %w", err)
	}
	out := make([]Offering, 0, len(rows))
	for _, row := range rows {
		out = append(out, toOffering(row))
	}
	return out, nil
}

func (r *Repository) GetServiceByID(ctx context.Context, id int64) (Offering, error) {
	row, err := r.queriesFor(ctx).GetServiceByID(ctx, id)
	if err != nil {
		return Offering{}, fmt.Errorf("get service by id: %w", err)
	}
	return toOffering(row), nil
}

func (r *Repository) ListServicesByIDs(ctx context.Context, ids []int64) ([]Offering, error) {
	rows, err := r.queriesFor(ctx).ListServicesByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("list services by ids: %w", err)
	}
	out := make([]Offering, 0, len(rows))
	for _, row := range rows {
		out = append(out, toOffering(row))
	}
	return out, nil
}

func (r *Repository) CreateService(ctx context.Context, categoryID int64, name, description string) (Offering, error) {
	row, err := r.queriesFor(ctx).CreateService(ctx, db.CreateServiceParams{
		CategoryID:  categoryID,
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
	})
	if err != nil {
		return Offering{}, fmt.Errorf("create service: %w", err)
	}
	return toOffering(row), nil
}

// --- master profile & specializations ---

func (r *Repository) UpsertMasterProfile(ctx context.Context, userID int64, city, bio string) (MasterProfile, error) {
	row, err := r.queriesFor(ctx).UpsertMasterProfile(ctx, db.UpsertMasterProfileParams{
		UserID: userID,
		City:   pgtype.Text{String: city, Valid: city != ""},
		Bio:    pgtype.Text{String: bio, Valid: bio != ""},
	})
	if err != nil {
		return MasterProfile{}, fmt.Errorf("upsert master profile: %w", err)
	}
	return toMasterProfile(row), nil
}

func (r *Repository) GetMasterProfile(ctx context.Context, userID int64) (MasterProfile, error) {
	row, err := r.queriesFor(ctx).GetMasterProfile(ctx, userID)
	if err != nil {
		return MasterProfile{}, fmt.Errorf("get master profile: %w", err)
	}
	return toMasterProfile(row), nil
}

func (r *Repository) ReplaceMasterSpecializations(ctx context.Context, masterUserID int64, serviceIDs []int64) error {
	q := r.queriesFor(ctx)
	if err := q.ReplaceMasterSpecializations(ctx, masterUserID); err != nil {
		return fmt.Errorf("clear master specializations: %w", err)
	}
	for _, serviceID := range serviceIDs {
		if err := q.AddMasterSpecialization(ctx, db.AddMasterSpecializationParams{
			MasterUserID: masterUserID,
			ServiceID:    serviceID,
		}); err != nil {
			return fmt.Errorf("add master specialization: %w", err)
		}
	}
	return nil
}

func (r *Repository) ListMasterSpecializationIDs(ctx context.Context, masterUserID int64) ([]int64, error) {
	ids, err := r.queriesFor(ctx).ListMasterSpecializationIDs(ctx, masterUserID)
	if err != nil {
		return nil, fmt.Errorf("list master specialization ids: %w", err)
	}
	return ids, nil
}

func (r *Repository) MasterHasSpecialization(ctx context.Context, masterUserID, serviceID int64) (bool, error) {
	has, err := r.queriesFor(ctx).MasterHasSpecialization(ctx, db.MasterHasSpecializationParams{
		MasterUserID: masterUserID,
		ServiceID:    serviceID,
	})
	if err != nil {
		return false, fmt.Errorf("check master specialization: %w", err)
	}
	return has, nil
}

// --- requests (заявки) ---

func (r *Repository) CreateRequest(ctx context.Context, clientID, serviceID int64, description, addressText string, latitude, longitude *float64) (Request, error) {
	row, err := r.queriesFor(ctx).CreateServiceRequest(ctx, db.CreateServiceRequestParams{
		ClientID:    clientID,
		ServiceID:   serviceID,
		Description: description,
		AddressText: addressText,
		Latitude:    float8OrNil(latitude),
		Longitude:   float8OrNil(longitude),
	})
	if err != nil {
		return Request{}, fmt.Errorf("create service request: %w", err)
	}
	return toRequest(row), nil
}

func (r *Repository) GetRequestByID(ctx context.Context, id int64) (Request, error) {
	row, err := r.queriesFor(ctx).GetServiceRequestByID(ctx, id)
	if err != nil {
		return Request{}, fmt.Errorf("get service request by id: %w", err)
	}
	return toRequest(row), nil
}

func (r *Repository) ListByClient(ctx context.Context, clientID int64, limit, offset int32) ([]Request, error) {
	rows, err := r.queriesFor(ctx).ListRequestsByClient(ctx, db.ListRequestsByClientParams{
		ClientID:   clientID,
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list requests by client: %w", err)
	}
	return toRequests(rows), nil
}

func (r *Repository) CountByClient(ctx context.Context, clientID int64) (int64, error) {
	count, err := r.queriesFor(ctx).CountRequestsByClient(ctx, clientID)
	if err != nil {
		return 0, fmt.Errorf("count requests by client: %w", err)
	}
	return count, nil
}

func (r *Repository) ListByMaster(ctx context.Context, masterID int64, limit, offset int32) ([]Request, error) {
	rows, err := r.queriesFor(ctx).ListRequestsByMaster(ctx, db.ListRequestsByMasterParams{
		MasterID:   pgtype.Int8{Int64: masterID, Valid: true},
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list requests by master: %w", err)
	}
	return toRequests(rows), nil
}

func (r *Repository) CountByMaster(ctx context.Context, masterID int64) (int64, error) {
	count, err := r.queriesFor(ctx).CountRequestsByMaster(ctx, pgtype.Int8{Int64: masterID, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("count requests by master: %w", err)
	}
	return count, nil
}

func (r *Repository) ListOpenForMaster(ctx context.Context, masterUserID int64, limit, offset int32) ([]Request, error) {
	rows, err := r.queriesFor(ctx).ListOpenRequestsForMaster(ctx, db.ListOpenRequestsForMasterParams{
		MasterUserID: masterUserID,
		PageLimit:    limit,
		PageOffset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list open requests for master: %w", err)
	}
	return toRequests(rows), nil
}

func (r *Repository) CountOpenForMaster(ctx context.Context, masterUserID int64) (int64, error) {
	count, err := r.queriesFor(ctx).CountOpenRequestsForMaster(ctx, masterUserID)
	if err != nil {
		return 0, fmt.Errorf("count open requests for master: %w", err)
	}
	return count, nil
}

// Claim is atomic (WHERE status = 'OPEN') as a defense against a race
// between two masters claiming the same request concurrently -- the
// specialization check happens earlier in the service layer, not here.
func (r *Repository) Claim(ctx context.Context, id, masterID int64) (Request, error) {
	row, err := r.queriesFor(ctx).ClaimServiceRequest(ctx, db.ClaimServiceRequestParams{
		MasterID: pgtype.Int8{Int64: masterID, Valid: true},
		ID:       id,
	})
	if err != nil {
		return Request{}, fmt.Errorf("claim service request: %w", err)
	}
	return toRequest(row), nil
}

func (r *Repository) Complete(ctx context.Context, id int64) (Request, error) {
	row, err := r.queriesFor(ctx).CompleteServiceRequest(ctx, id)
	if err != nil {
		return Request{}, fmt.Errorf("complete service request: %w", err)
	}
	return toRequest(row), nil
}

func (r *Repository) Cancel(ctx context.Context, id int64, reason string) (Request, error) {
	row, err := r.queriesFor(ctx).CancelServiceRequest(ctx, db.CancelServiceRequestParams{
		ID:           id,
		CancelReason: pgtype.Text{String: reason, Valid: reason != ""},
	})
	if err != nil {
		return Request{}, fmt.Errorf("cancel service request: %w", err)
	}
	return toRequest(row), nil
}

// --- row -> domain model conversions ---

func toCategory(row db.ServiceCategory) Category {
	return Category{
		ID:        row.ID,
		Name:      row.Name,
		Active:    row.Active,
		CreatedAt: row.CreatedAt.Time,
	}
}

func toOffering(row db.Service) Offering {
	description := ""
	if row.Description.Valid {
		description = row.Description.String
	}
	return Offering{
		ID:          row.ID,
		CategoryID:  row.CategoryID,
		Name:        row.Name,
		Description: description,
		Active:      row.Active,
		CreatedAt:   row.CreatedAt.Time,
	}
}

func toMasterProfile(row db.MasterProfile) MasterProfile {
	city, bio := "", ""
	if row.City.Valid {
		city = row.City.String
	}
	if row.Bio.Valid {
		bio = row.Bio.String
	}
	return MasterProfile{
		UserID:    row.UserID,
		City:      city,
		Bio:       bio,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func toRequest(row db.ServiceRequest) Request {
	var masterID *int64
	if row.MasterID.Valid {
		masterID = &row.MasterID.Int64
	}
	cancelReason := ""
	if row.CancelReason.Valid {
		cancelReason = row.CancelReason.String
	}
	return Request{
		ID:           row.ID,
		ClientID:     row.ClientID,
		OfferingID:   row.ServiceID,
		Description:  row.Description,
		AddressText:  row.AddressText,
		Latitude:     float64OrNil(row.Latitude),
		Longitude:    float64OrNil(row.Longitude),
		Status:       RequestStatus(row.Status),
		MasterID:     masterID,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		ClaimedAt:    timeOrNil(row.ClaimedAt),
		CompletedAt:  timeOrNil(row.CompletedAt),
		CanceledAt:   timeOrNil(row.CanceledAt),
		CancelReason: cancelReason,
	}
}

func toRequests(rows []db.ServiceRequest) []Request {
	out := make([]Request, 0, len(rows))
	for _, row := range rows {
		out = append(out, toRequest(row))
	}
	return out
}

func float8OrNil(f *float64) pgtype.Float8 {
	if f == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *f, Valid: true}
}

func float64OrNil(f pgtype.Float8) *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

func timeOrNil(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
