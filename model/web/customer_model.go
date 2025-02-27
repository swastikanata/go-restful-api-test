package web

type CustomerCreateRequest struct {
	Name       string `validate:"required,min=1,max=100" json:"name"`
	Email      string `validate:"required" json:"column:email"`
	Phone      string `validate:"required,min=1,max=100" json:"column:phone"`
	Address    string `validate:"required,min=1,max=100" json:"column:address"`
	LoyaltyPts int    `validate:"required,min=0" json:"column:loyalty_points"`
}

type CustomerUpdateRequest struct {
	Id         uint64 `validate:"required"`
	Name       string `validate:"required,min=1,max=100" json:"name"`
	Email      string `validate:"required" json:"column:email"`
	Phone      string `validate:"required,min=1,max=100" json:"column:phone"`
	Address    string `validate:"required,min=1,max=100" json:"column:address"`
	LoyaltyPts int    `validate:"required,min=0" json:"column:loyalty_points"`
}

type CustomerResponse struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	LoyaltyPts int    `json:"loyalty_points"`
}
