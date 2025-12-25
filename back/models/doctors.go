package models

type Doctor struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Department   string         `json:"department"` // 科室
	Title        string         `json:"title"`      // 职称
	Introduction string         `json:"introduction"`
	Photo        string         `json:"photo"`    // 照片URL
	Diseases     []string       `json:"diseases"` // 管理的病种ID列表 (1-3个)
	WorkSchedule []WorkSchedule `json:"workSchedule"`
	MaxPatients  int            `json:"maxPatients"` // 每日最大接诊数
	Fee          float64        `json:"fee"`         // 挂号费
}

type WorkSchedule struct {
	DayOfWeek   string `json:"dayOfWeek"` // 周一至周日
	StartTime   string `json:"startTime"` // 09:00
	EndTime     string `json:"endTime"`   // 17:00
	IsAvailable bool   `json:"isAvailable"`
}
