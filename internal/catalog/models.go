package catalog

import "time"

type Category struct {
	ID               int64
	ParentCategoryID *int64
	Name             string
	Active           bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Offering is a concrete service a client can order within a Category.
// Named Offering rather than Service to avoid colliding with this package's
// own business-logic Service struct -- "service" the architectural layer
// and "service" the marketplace noun are unrelated concepts that happen to
// share an English word.
type Offering struct {
	ID          int64
	CategoryID  int64
	Name        string
	Description string
	PriceFrom   *float64
	PriceTo     *float64
	Unit        string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
