package act

import "time"

type InspectionType string

const (
	InspectionScheduled   InspectionType = "SCHEDULED"
	InspectionUnscheduled InspectionType = "UNSCHEDULED"
)

func (t InspectionType) Valid() bool {
	switch t {
	case InspectionScheduled, InspectionUnscheduled:
		return true
	default:
		return false
	}
}

type InspectionAct struct {
	ID             int64
	TaskID         int64
	TenantID       int64
	AddressID      int64
	InspectionDate *time.Time
	ConsumerID     *int64
	InspectionType InspectionType
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ReplacementAct struct {
	ID               int64
	TaskID           int64
	TenantID         int64
	AddressID        int64
	AccountNumber    string
	InstallationDate *time.Time
	OldBrand         string
	OldSerialNumber  string
	OldReadings      *float64
	NewBrand         string
	NewSerialNumber  string
	NewReadings      *float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
