package organization

import "time"

type TenantPlan string

const (
	PlanFree       TenantPlan = "FREE"
	PlanBusiness   TenantPlan = "BUSINESS"
	PlanEnterprise TenantPlan = "ENTERPRISE"
)

func (p TenantPlan) UserLimit() int {
	switch p {
	case PlanFree:
		return 5
	case PlanBusiness:
		return 50
	case PlanEnterprise:
		return int(^uint(0) >> 1)
	default:
		return 5
	}
}

type Tenant struct {
	ID        int64
	Name      string
	Code      string
	Plan      TenantPlan
	Active    bool
	CreatedAt time.Time
}
