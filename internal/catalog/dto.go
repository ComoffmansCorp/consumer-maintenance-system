package catalog

type CategoryDTO struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Active        bool          `json:"active"`
	Subcategories []CategoryDTO `json:"subcategories,omitempty"`
}

type ServiceDTO struct {
	ID          int64    `json:"id"`
	CategoryID  int64    `json:"categoryId"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	PriceFrom   *float64 `json:"priceFrom,omitempty"`
	PriceTo     *float64 `json:"priceTo,omitempty"`
	Unit        string   `json:"unit,omitempty"`
	Active      bool     `json:"active"`
}

type CreateCategoryRequest struct {
	Name             string `json:"name"`
	ParentCategoryID *int64 `json:"parentCategoryId,omitempty"`
}

type UpdateCategoryRequest struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type CreateServiceRequest struct {
	CategoryID  int64    `json:"categoryId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	PriceFrom   *float64 `json:"priceFrom,omitempty"`
	PriceTo     *float64 `json:"priceTo,omitempty"`
	Unit        string   `json:"unit,omitempty"`
}

type UpdateServiceRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	PriceFrom   *float64 `json:"priceFrom,omitempty"`
	PriceTo     *float64 `json:"priceTo,omitempty"`
	Unit        string   `json:"unit,omitempty"`
	Active      bool     `json:"active"`
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
		PriceFrom:   o.PriceFrom,
		PriceTo:     o.PriceTo,
		Unit:        o.Unit,
		Active:      o.Active,
	}
}
