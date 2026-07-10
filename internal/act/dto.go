package act

import (
	"time"

	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

type InspectionDTO struct {
	ID             int64          `json:"id"`
	TaskID         int64          `json:"taskId"`
	AddressID      int64          `json:"addressId"`
	AddressLabel   string         `json:"addressLabel,omitempty"`
	InspectionDate *httpx.Date    `json:"inspectionDate,omitempty"`
	ConsumerID     *int64         `json:"consumerId,omitempty"`
	ConsumerName   string         `json:"consumerName,omitempty"`
	InspectionType InspectionType `json:"inspectionType"`
	Notes          string         `json:"notes,omitempty"`
	MeterCount     int            `json:"meterCount"`
	PhotoCount     int            `json:"photoCount"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type CreateInspectionRequest struct {
	TaskID         int64          `json:"taskId"`
	InspectionDate *httpx.Date    `json:"inspectionDate"`
	ConsumerID     *int64         `json:"consumerId"`
	InspectionType InspectionType `json:"inspectionType"`
	Notes          string         `json:"notes"`
}

type UpdateInspectionRequest struct {
	InspectionDate *httpx.Date    `json:"inspectionDate"`
	ConsumerID     *int64         `json:"consumerId"`
	InspectionType InspectionType `json:"inspectionType"`
	Notes          string         `json:"notes"`
}

type ReplacementDTO struct {
	ID               int64       `json:"id"`
	TaskID           int64       `json:"taskId"`
	AddressID        int64       `json:"addressId"`
	AddressLabel     string      `json:"addressLabel,omitempty"`
	AccountNumber    string      `json:"accountNumber"`
	InstallationDate *httpx.Date `json:"installationDate,omitempty"`
	OldBrand         string      `json:"oldBrand,omitempty"`
	OldSerialNumber  string      `json:"oldSerialNumber,omitempty"`
	OldReadings      *float64    `json:"oldReadings,omitempty"`
	NewBrand         string      `json:"newBrand,omitempty"`
	NewSerialNumber  string      `json:"newSerialNumber,omitempty"`
	NewReadings      *float64    `json:"newReadings,omitempty"`
	PhotoCount       int         `json:"photoCount"`
	CreatedAt        time.Time   `json:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt"`
}

type CreateReplacementRequest struct {
	TaskID           int64       `json:"taskId"`
	AccountNumber    string      `json:"accountNumber"`
	InstallationDate *httpx.Date `json:"installationDate"`
	OldBrand         string      `json:"oldBrand"`
	OldSerialNumber  string      `json:"oldSerialNumber"`
	OldReadings      *float64    `json:"oldReadings"`
	NewBrand         string      `json:"newBrand"`
	NewSerialNumber  string      `json:"newSerialNumber"`
	NewReadings      *float64    `json:"newReadings"`
}

type UpdateReplacementRequest struct {
	AccountNumber    string      `json:"accountNumber"`
	InstallationDate *httpx.Date `json:"installationDate"`
	OldBrand         string      `json:"oldBrand"`
	OldSerialNumber  string      `json:"oldSerialNumber"`
	OldReadings      *float64    `json:"oldReadings"`
	NewBrand         string      `json:"newBrand"`
	NewSerialNumber  string      `json:"newSerialNumber"`
	NewReadings      *float64    `json:"newReadings"`
}

func toInspectionDTO(a InspectionAct) InspectionDTO {
	return InspectionDTO{
		ID:             a.ID,
		TaskID:         a.TaskID,
		AddressID:      a.AddressID,
		InspectionDate: httpx.TimeToDate(a.InspectionDate),
		ConsumerID:     a.ConsumerID,
		InspectionType: a.InspectionType,
		Notes:          a.Notes,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

func toReplacementDTO(a ReplacementAct) ReplacementDTO {
	return ReplacementDTO{
		ID:               a.ID,
		TaskID:           a.TaskID,
		AddressID:        a.AddressID,
		AccountNumber:    a.AccountNumber,
		InstallationDate: httpx.TimeToDate(a.InstallationDate),
		OldBrand:         a.OldBrand,
		OldSerialNumber:  a.OldSerialNumber,
		OldReadings:      a.OldReadings,
		NewBrand:         a.NewBrand,
		NewSerialNumber:  a.NewSerialNumber,
		NewReadings:      a.NewReadings,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
	}
}
