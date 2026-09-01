package request

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	platformdb "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/db"
	db "github.com/myurbondarchuk/consumer-maintenance-system/internal/request/db"
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

// --- requests ---

func (r *Repository) CreateRequest(ctx context.Context, clientID, serviceID int64, description, addressText string, latitude, longitude *float64) (ServiceRequest, error) {
	row, err := r.queriesFor(ctx).CreateServiceRequest(ctx, db.CreateServiceRequestParams{
		ClientID:    clientID,
		ServiceID:   serviceID,
		Description: description,
		AddressText: addressText,
		Latitude:    float8OrNil(latitude),
		Longitude:   float8OrNil(longitude),
	})
	if err != nil {
		return ServiceRequest{}, fmt.Errorf("create service request: %w", err)
	}
	return toRequest(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (ServiceRequest, error) {
	row, err := r.queriesFor(ctx).GetServiceRequestByID(ctx, id)
	if err != nil {
		return ServiceRequest{}, fmt.Errorf("get service request by id: %w", err)
	}
	return toRequest(row), nil
}

func (r *Repository) ListByClient(ctx context.Context, clientID int64, limit, offset int32) ([]ServiceRequest, error) {
	rows, err := r.queriesFor(ctx).ListRequestsByClient(ctx, db.ListRequestsByClientParams{
		ClientID: clientID, PageLimit: limit, PageOffset: offset,
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

func (r *Repository) ListByMaster(ctx context.Context, masterID int64, limit, offset int32) ([]ServiceRequest, error) {
	rows, err := r.queriesFor(ctx).ListRequestsByMaster(ctx, db.ListRequestsByMasterParams{
		MasterID: pgtype.Int8{Int64: masterID, Valid: true}, PageLimit: limit, PageOffset: offset,
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

func (r *Repository) ListOpenForMaster(ctx context.Context, masterUserID int64, limit, offset int32) ([]ServiceRequest, error) {
	rows, err := r.queriesFor(ctx).ListOpenRequestsForMaster(ctx, db.ListOpenRequestsForMasterParams{
		MasterUserID: masterUserID, PageLimit: limit, PageOffset: offset,
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

func (r *Repository) ListAdmin(ctx context.Context, status string, limit, offset int32) ([]ServiceRequest, error) {
	rows, err := r.queriesFor(ctx).ListRequestsAdmin(ctx, db.ListRequestsAdminParams{
		Status: status, PageLimit: limit, PageOffset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list requests admin: %w", err)
	}
	return toRequests(rows), nil
}

func (r *Repository) CountAdmin(ctx context.Context, status string) (int64, error) {
	count, err := r.queriesFor(ctx).CountRequestsAdmin(ctx, status)
	if err != nil {
		return 0, fmt.Errorf("count requests admin: %w", err)
	}
	return count, nil
}

// Assign is atomic (WHERE status = 'OPEN') as a defense against a race
// between two offer-accept calls landing concurrently.
func (r *Repository) Assign(ctx context.Context, id, masterID int64, agreedPrice float64) (ServiceRequest, error) {
	price := agreedPrice
	row, err := r.queriesFor(ctx).AssignServiceRequest(ctx, db.AssignServiceRequestParams{
		MasterID:    pgtype.Int8{Int64: masterID, Valid: true},
		AgreedPrice: &price,
		ID:          id,
	})
	if err != nil {
		return ServiceRequest{}, fmt.Errorf("assign service request: %w", err)
	}
	return toRequest(row), nil
}

func (r *Repository) Complete(ctx context.Context, id int64) (ServiceRequest, error) {
	row, err := r.queriesFor(ctx).CompleteServiceRequest(ctx, id)
	if err != nil {
		return ServiceRequest{}, fmt.Errorf("complete service request: %w", err)
	}
	return toRequest(row), nil
}

func (r *Repository) Cancel(ctx context.Context, id int64, reason string) (ServiceRequest, error) {
	row, err := r.queriesFor(ctx).CancelServiceRequest(ctx, db.CancelServiceRequestParams{
		ID:           id,
		CancelReason: pgtype.Text{String: reason, Valid: reason != ""},
	})
	if err != nil {
		return ServiceRequest{}, fmt.Errorf("cancel service request: %w", err)
	}
	return toRequest(row), nil
}

// --- offers ---

func (r *Repository) CreateOffer(ctx context.Context, requestID, masterID int64, price float64, comment string) (Offer, error) {
	row, err := r.queriesFor(ctx).CreateOffer(ctx, db.CreateOfferParams{
		RequestID: requestID,
		MasterID:  masterID,
		Price:     price,
		Comment:   pgtype.Text{String: comment, Valid: comment != ""},
	})
	if err != nil {
		return Offer{}, fmt.Errorf("create offer: %w", err)
	}
	return toOffer(row), nil
}

func (r *Repository) GetOfferByID(ctx context.Context, id int64) (Offer, error) {
	row, err := r.queriesFor(ctx).GetOfferByID(ctx, id)
	if err != nil {
		return Offer{}, fmt.Errorf("get offer by id: %w", err)
	}
	return toOffer(row), nil
}

func (r *Repository) ListOffersByRequest(ctx context.Context, requestID int64) ([]Offer, error) {
	rows, err := r.queriesFor(ctx).ListOffersByRequest(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("list offers by request: %w", err)
	}
	out := make([]Offer, 0, len(rows))
	for _, row := range rows {
		out = append(out, toOffer(row))
	}
	return out, nil
}

func (r *Repository) AcceptOffer(ctx context.Context, id int64) (Offer, error) {
	row, err := r.queriesFor(ctx).AcceptOffer(ctx, id)
	if err != nil {
		return Offer{}, fmt.Errorf("accept offer: %w", err)
	}
	return toOffer(row), nil
}

func (r *Repository) RejectOtherOffers(ctx context.Context, requestID, acceptedOfferID int64) error {
	if err := r.queriesFor(ctx).RejectOtherOffers(ctx, db.RejectOtherOffersParams{
		RequestID: requestID, AcceptedOfferID: acceptedOfferID,
	}); err != nil {
		return fmt.Errorf("reject other offers: %w", err)
	}
	return nil
}

// --- status history ---

func (r *Repository) AddStatusHistory(ctx context.Context, requestID int64, fromStatus, toStatus string, changedBy int64, comment string) error {
	if err := r.queriesFor(ctx).CreateStatusHistory(ctx, db.CreateStatusHistoryParams{
		RequestID:  requestID,
		FromStatus: pgtype.Text{String: fromStatus, Valid: fromStatus != ""},
		ToStatus:   toStatus,
		ChangedBy:  changedBy,
		Comment:    pgtype.Text{String: comment, Valid: comment != ""},
	}); err != nil {
		return fmt.Errorf("create status history: %w", err)
	}
	return nil
}

func (r *Repository) ListStatusHistory(ctx context.Context, requestID int64) ([]StatusHistoryEntry, error) {
	rows, err := r.queriesFor(ctx).ListStatusHistoryByRequest(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("list status history: %w", err)
	}
	out := make([]StatusHistoryEntry, 0, len(rows))
	for _, row := range rows {
		fromStatus := ""
		if row.FromStatus.Valid {
			fromStatus = row.FromStatus.String
		}
		comment := ""
		if row.Comment.Valid {
			comment = row.Comment.String
		}
		out = append(out, StatusHistoryEntry{
			ID: row.ID, RequestID: row.RequestID, FromStatus: fromStatus, ToStatus: row.ToStatus,
			ChangedBy: row.ChangedBy, Comment: comment, CreatedAt: row.CreatedAt.Time,
		})
	}
	return out, nil
}

// --- favorites ---

func (r *Repository) AddFavorite(ctx context.Context, clientID, masterID int64) error {
	if err := r.queriesFor(ctx).AddFavorite(ctx, db.AddFavoriteParams{ClientID: clientID, MasterID: masterID}); err != nil {
		return fmt.Errorf("add favorite: %w", err)
	}
	return nil
}

func (r *Repository) RemoveFavorite(ctx context.Context, clientID, masterID int64) error {
	if err := r.queriesFor(ctx).RemoveFavorite(ctx, db.RemoveFavoriteParams{ClientID: clientID, MasterID: masterID}); err != nil {
		return fmt.Errorf("remove favorite: %w", err)
	}
	return nil
}

func (r *Repository) ListFavoritesByClient(ctx context.Context, clientID int64) ([]Favorite, error) {
	rows, err := r.queriesFor(ctx).ListFavoritesByClient(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("list favorites by client: %w", err)
	}
	out := make([]Favorite, 0, len(rows))
	for _, row := range rows {
		out = append(out, Favorite{ClientID: row.ClientID, MasterID: row.MasterID, CreatedAt: row.CreatedAt.Time})
	}
	return out, nil
}

// --- row -> domain model conversions ---

func toRequest(row db.ServiceRequest) ServiceRequest {
	var masterID *int64
	if row.MasterID.Valid {
		masterID = &row.MasterID.Int64
	}
	cancelReason := ""
	if row.CancelReason.Valid {
		cancelReason = row.CancelReason.String
	}
	return ServiceRequest{
		ID:           row.ID,
		ClientID:     row.ClientID,
		ServiceID:    row.ServiceID,
		Description:  row.Description,
		AddressText:  row.AddressText,
		Latitude:     float64OrNil(row.Latitude),
		Longitude:    float64OrNil(row.Longitude),
		Status:       Status(row.Status),
		MasterID:     masterID,
		AgreedPrice:  row.AgreedPrice,
		CancelReason: cancelReason,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}

func toRequests(rows []db.ServiceRequest) []ServiceRequest {
	out := make([]ServiceRequest, 0, len(rows))
	for _, row := range rows {
		out = append(out, toRequest(row))
	}
	return out
}

func toOffer(row db.RequestOffer) Offer {
	comment := ""
	if row.Comment.Valid {
		comment = row.Comment.String
	}
	return Offer{
		ID:        row.ID,
		RequestID: row.RequestID,
		MasterID:  row.MasterID,
		Price:     row.Price,
		Comment:   comment,
		Status:    OfferStatus(row.Status),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
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
