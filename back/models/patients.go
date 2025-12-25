package models

import (
	"time"
)

type Patient struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Gender           string    `json:"gender"` // 男, 女
	Age              int       `json:"age"`
	Phone            string    `json:"phone"`
	IDCard           string    `json:"idCard"` // 身份证号
	Address          string    `json:"address"`
	EmergencyContact string    `json:"emergencyContact"`
	EmergencyPhone   string    `json:"emergencyPhone"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
