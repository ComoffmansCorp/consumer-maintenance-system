package organization

type TenantDTO struct {
	ID     int64      `json:"id"`
	Name   string     `json:"name"`
	Code   string     `json:"code"`
	Plan   TenantPlan `json:"plan"`
	Active bool       `json:"active"`
}

type CreateTenantRequest struct {
	Name string     `json:"name"`
	Code string     `json:"code"`
	Plan TenantPlan `json:"plan"`
}

type UpdateTenantPlanRequest struct {
	Plan TenantPlan `json:"plan"`
}
