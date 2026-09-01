package marketplace

import "time"

type CategoryDTO struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type ServiceDTO struct {
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"categoryId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Active      bool   `json:"active"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type CreateOfferingRequest struct {
	CategoryID  int64  `json:"categoryId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MasterProfileDTO struct {
	City              string  `json:"city,omitempty"`
	Bio               string  `json:"bio,omitempty"`
	SpecializationIDs []int64 `json:"specializationIds"`
}

type UpdateMasterProfileRequest struct {
	City              string  `json:"city"`
	Bio               string  `json:"bio"`
	SpecializationIDs []int64 `json:"specializationIds"`
}

type CreateRequestRequest struct {
	ServiceID   int64    `json:"serviceId"`
	Description string   `json:"description"`
	AddressText string   `json:"addressText"`
	// Latitude/Longitude are optional -- populated from the Yandex Suggest
	// response when the client picks a suggested address, omitted for a
	// free-typed address with no matching suggestion.
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
}

type CancelRequestBody struct {
	Reason string `json:"reason"`
}

type RequestDTO struct {
	ID           int64         `json:"id"`
	ServiceID    int64         `json:"serviceId"`
	ServiceName  string        `json:"serviceName,omitempty"`
	CategoryName string        `json:"categoryName,omitempty"`
	Description  string        `json:"description"`
	AddressText  string        `json:"addressText"`
	Latitude     *float64      `json:"latitude,omitempty"`
	Longitude    *float64      `json:"longitude,omitempty"`
	Status       RequestStatus `json:"status"`
	ClientID     int64         `json:"clientId"`
	ClientName   string        `json:"clientName,omitempty"`
	MasterID     *int64        `json:"masterId,omitempty"`
	MasterName   string        `json:"masterName,omitempty"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	ClaimedAt    *time.Time    `json:"claimedAt,omitempty"`
	CompletedAt  *time.Time    `json:"completedAt,omitempty"`
	CanceledAt   *time.Time    `json:"canceledAt,omitempty"`
	CancelReason string        `json:"cancelReason,omitempty"`
}

func ToCategoryDTO(c Category) CategoryDTO {
	return CategoryDTO{ID: c.ID, Name: c.Name, Active: c.Active}
}

func ToServiceDTO(o Offering) ServiceDTO {
	return ServiceDTO{
		ID:          o.ID,
		CategoryID:  o.CategoryID,
		Name:        o.Name,
		Description: o.Description,
		Active:      o.Active,
	}
}

func ToRequestDTO(r Request) RequestDTO {
	return RequestDTO{
		ID:           r.ID,
		ServiceID:    r.OfferingID,
		Description:  r.Description,
		AddressText:  r.AddressText,
		Latitude:     r.Latitude,
		Longitude:    r.Longitude,
		Status:       r.Status,
		ClientID:     r.ClientID,
		MasterID:     r.MasterID,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
		ClaimedAt:    r.ClaimedAt,
		CompletedAt:  r.CompletedAt,
		CanceledAt:   r.CanceledAt,
		CancelReason: r.CancelReason,
	}
}
