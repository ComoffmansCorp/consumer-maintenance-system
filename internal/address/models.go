package address

import "time"

type Address struct {
	ID         int64
	Street     string
	House      string
	Building   string
	Apartment  string
	TenantID   int64
	ConsumerID *int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
