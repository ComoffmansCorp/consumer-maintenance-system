package consumer

import "time"

type Type string

const (
	TypeCommercial Type = "COMMERCIAL"
	TypeGovernment Type = "GOVERNMENT"
	TypeResidential Type = "RESIDENTIAL"
)

func (t Type) Valid() bool {
	switch t {
	case TypeCommercial, TypeGovernment, TypeResidential:
		return true
	default:
		return false
	}
}

type Consumer struct {
	ID          int64
	Name        string
	Type        Type
	Description string
	TenantID    int64
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
