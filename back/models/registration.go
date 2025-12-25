package models

import "time"

type Registration struct {
	ID               string    `json:"id"`
	PatientID        string    `json:"patientId"`
	DoctorID         string    `json:"doctorId"`
	Department       string    `json:"department"` // 科室
	Departments      []string  `json:"departments,omitempty"`
	RegistrationDate time.Time `json:"registrationDate"`
	VisitDate        time.Time `json:"visitDate"`
	TimeSlot         string    `json:"timeSlot"` // 时间段
	Status           string    `json:"status"`   // pending, confirmed, completed, cancelled
	Symptoms         string    `json:"symptoms"` // 症状描述
	Notes            string    `json:"notes"`    // 备注
	CreatedAt        time.Time `json:"createdAt"`
}
