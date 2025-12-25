package models

type Admin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // 实际应用中应该加密存储
	Role     string `json:"role"`     // admin, operator等
}
