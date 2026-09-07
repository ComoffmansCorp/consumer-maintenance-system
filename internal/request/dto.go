package request

import "time"

type CreateRequestRequest struct {
	ServiceID   int64  `json:"serviceId"`
	Description string `json:"description"`
	AddressText string `json:"addressText"`
	// Latitude/Longitude are optional -- populated from the Yandex Suggest
	// response when the client picks a suggested address, omitted for a
	// free-typed address with no matching suggestion.
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type CancelRequestBody struct {
	Reason string `json:"reason"`
}

type SubmitOfferRequest struct {
	Price   float64 `json:"price"`
	Comment string  `json:"comment"`
}

type StatusHistoryEntryDTO struct {
	FromStatus string    `json:"fromStatus,omitempty"`
	ToStatus   string    `json:"toStatus"`
	ChangedBy  int64     `json:"changedBy"`
	Comment    string    `json:"comment,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type RequestDTO struct {
	ID           int64                   `json:"id"`
	ServiceID    int64                   `json:"serviceId"`
	ServiceName  string                  `json:"serviceName,omitempty"`
	Description  string                  `json:"description"`
	AddressText  string                  `json:"addressText"`
	Latitude     *float64                `json:"latitude,omitempty"`
	Longitude    *float64                `json:"longitude,omitempty"`
	Status       Status                  `json:"status"`
	ClientID     int64                   `json:"clientId"`
	MasterID     *int64                  `json:"masterId,omitempty"`
	AgreedPrice  *float64                `json:"agreedPrice,omitempty"`
	CancelReason string                  `json:"cancelReason,omitempty"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
	History      []StatusHistoryEntryDTO `json:"history,omitempty"`
}

type OfferDTO struct {
	ID        int64       `json:"id"`
	RequestID int64       `json:"requestId"`
	MasterID  int64       `json:"masterId"`
	Price     float64     `json:"price"`
	Comment   string      `json:"comment,omitempty"`
	Status    OfferStatus `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	// MasterAvatarURL is enrichment, not stored on the offer itself -- see
	// Service.ListOffers, which fills it in via SpecializationPort the same
	// way RequestDTO.ServiceName is filled in via CatalogPort.
	MasterAvatarURL *string `json:"masterAvatarUrl,omitempty"`
}

type FavoriteDTO struct {
	MasterID  int64     `json:"masterId"`
	CreatedAt time.Time `json:"createdAt"`
}

func ToRequestDTO(r ServiceRequest) RequestDTO {
	return RequestDTO{
		ID:           r.ID,
		ServiceID:    r.ServiceID,
		Description:  r.Description,
		AddressText:  r.AddressText,
		Latitude:     r.Latitude,
		Longitude:    r.Longitude,
		Status:       r.Status,
		ClientID:     r.ClientID,
		MasterID:     r.MasterID,
		AgreedPrice:  r.AgreedPrice,
		CancelReason: r.CancelReason,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func ToOfferDTO(o Offer) OfferDTO {
	return OfferDTO{
		ID:        o.ID,
		RequestID: o.RequestID,
		MasterID:  o.MasterID,
		Price:     o.Price,
		Comment:   o.Comment,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func ToStatusHistoryDTO(h StatusHistoryEntry) StatusHistoryEntryDTO {
	return StatusHistoryEntryDTO{
		FromStatus: h.FromStatus,
		ToStatus:   h.ToStatus,
		ChangedBy:  h.ChangedBy,
		Comment:    h.Comment,
		CreatedAt:  h.CreatedAt,
	}
}

func ToFavoriteDTO(f Favorite) FavoriteDTO {
	return FavoriteDTO{MasterID: f.MasterID, CreatedAt: f.CreatedAt}
}
