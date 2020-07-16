package cmcserverinfo

import "time"

// AvailableProduct login account for many products
type AvailableProduct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// LoginResult service return json format
type LoginResult struct {
	Username          string             `json:"username"`
	FirstName         string             `json:"first_name"`
	LastName          string             `json:"last_name"`
	Companygroup      interface{}        `json:"companygroup"`
	LastPwdReset      time.Time          `json:"last_pwd_reset"`
	Company           int                `json:"company"`
	IsActive          bool               `json:"is_active"`
	Site              int                `json:"site"`
	Managedsites      []interface{}      `json:"managedsites"`
	Email             string             `json:"email"`
	IsSuperuser       bool               `json:"is_superuser"`
	IsStaff           bool               `json:"is_staff"`
	LastLogin         time.Time          `json:"last_login"`
	AvailableProducts []AvailableProduct `json:"available_products"`
	Groups            []int              `json:"groups"`
	Crosssites        []interface{}      `json:"crosssites"`
	Managedcompanys   []interface{}      `json:"managedcompanys"`
	ID                int                `json:"id"`
	ManageCredit      bool               `json:"manageCredit"`
	DateJoined        time.Time          `json:"date_joined"`
}
