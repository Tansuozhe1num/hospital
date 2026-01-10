package controllers

import (
	"errors"
	"hospital-system/auth"
	"hospital-system/models"
	"hospital-system/resource"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRegistrations(ctx *gin.Context) {
	claims, ok := auth.GetClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	account, err := resource.AccountService.GetByID(ctx, claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account not found"})
		return
	}

	registrations, err := resource.RegistrationService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if account.Role == "admin" {
		ctx.JSON(http.StatusOK, registrations)
		return
	}

	linkedID := account.LinkedID
	if linkedID == "" {
		ctx.JSON(http.StatusOK, []models.Registration{})
		return
	}

	filtered := make([]models.Registration, 0, len(registrations))
	switch account.Role {
	case "doctor":
		for _, r := range registrations {
			if r.DoctorID == linkedID {
				filtered = append(filtered, r)
			}
		}
	case "patient":
		for _, r := range registrations {
			if r.PatientID == linkedID {
				filtered = append(filtered, r)
			}
		}
	default:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	ctx.JSON(http.StatusOK, filtered)
}

func GetRegistration(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	claims, ok := auth.GetClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	account, err := resource.AccountService.GetByID(ctx, claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account not found"})
		return
	}

	registration, err := resource.RegistrationService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if account.Role == "admin" {
		ctx.JSON(http.StatusOK, registration)
		return
	}
	if account.LinkedID == "" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if account.Role == "doctor" && registration.DoctorID != account.LinkedID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	if account.Role == "patient" && registration.PatientID != account.LinkedID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	ctx.JSON(http.StatusOK, registration)
}

func CreateRegistration(ctx *gin.Context) {
	claims, ok := auth.GetClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	account, err := resource.AccountService.GetByID(ctx, claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account not found"})
		return
	}

	var registration models.Registration
	if err := ctx.ShouldBindJSON(&registration); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch account.Role {
	case "admin":
	case "patient":
		if account.LinkedID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "patient profile not linked"})
			return
		}
		registration.PatientID = account.LinkedID
		registration.Status = "pending"
	default:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if _, err := resource.PatientService.GetByID(ctx, registration.PatientID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "patient not found"})
		return
	}
	doctor, err := resource.DoctorService.GetByID(ctx, registration.DoctorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "doctor not found"})
		return
	}

	if len(registration.Departments) == 0 && registration.Department != "" {
		registration.Departments = []string{registration.Department}
	}
	if registration.Department == "" && len(registration.Departments) > 0 {
		registration.Department = registration.Departments[0]
	}

	if len(registration.Departments) == 0 {
		if doctor.Department == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "departments cannot be empty"})
			return
		}
		registration.Departments = []string{doctor.Department}
		registration.Department = doctor.Department
	}

	if err := validateDepartmentsExist(ctx, registration.Departments); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := resource.RegistrationService.Create(ctx, &registration); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, registration)
}

func UpdateRegistration(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	claims, ok := auth.GetClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	account, err := resource.AccountService.GetByID(ctx, claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account not found"})
		return
	}

	existing, err := resource.RegistrationService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "registration not found"})
		return
	}

	var registration models.Registration
	if err := ctx.ShouldBindJSON(&registration); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch account.Role {
	case "admin":
	case "doctor":
		if account.LinkedID == "" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if existing.DoctorID != account.LinkedID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		registration.PatientID = existing.PatientID
		registration.DoctorID = existing.DoctorID
		registration.Department = existing.Department
		registration.Departments = existing.Departments
		registration.RegistrationDate = existing.RegistrationDate
		registration.VisitDate = existing.VisitDate
		registration.TimeSlot = existing.TimeSlot
		registration.Symptoms = existing.Symptoms
		if strings.TrimSpace(registration.Status) == "" {
			registration.Status = existing.Status
		}
		if !isAllowedDoctorRegistrationStatusTransition(existing.Status, registration.Status) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status transition"})
			return
		}
	default:
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if _, err := resource.PatientService.GetByID(ctx, registration.PatientID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "patient not found"})
		return
	}
	doctor, err := resource.DoctorService.GetByID(ctx, registration.DoctorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "doctor not found"})
		return
	}

	if len(registration.Departments) == 0 && registration.Department != "" {
		registration.Departments = []string{registration.Department}
	}
	if registration.Department == "" && len(registration.Departments) > 0 {
		registration.Department = registration.Departments[0]
	}

	if len(registration.Departments) == 0 {
		if doctor.Department == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "departments cannot be empty"})
			return
		}
		registration.Departments = []string{doctor.Department}
		registration.Department = doctor.Department
	}

	if err := validateDepartmentsExist(ctx, registration.Departments); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := resource.RegistrationService.Update(ctx, id, &registration); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, registration)
}

func DeleteRegistration(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	if err := resource.RegistrationService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Registration deleted successfully"})
}

func isAllowedDoctorRegistrationStatusTransition(from string, to string) bool {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	if from == "" || to == "" {
		return false
	}
	if from == to {
		return true
	}
	switch from {
	case "pending":
		return to == "confirmed" || to == "cancelled"
	case "confirmed":
		return to == "completed" || to == "cancelled"
	case "completed", "cancelled":
		return false
	default:
		return false
	}
}

func validateDepartmentsExist(ctx *gin.Context, deptNames []string) error {
	if len(deptNames) == 0 {
		return nil
	}
	departments, err := resource.DepartmentService.GetAll(ctx)
	if err != nil {
		return err
	}
	existing := make(map[string]struct{}, len(departments))
	for _, d := range departments {
		existing[d.Name] = struct{}{}
	}
	for _, name := range deptNames {
		if _, ok := existing[strings.TrimSpace(name)]; !ok {
			return errors.New("department not found: " + name)
		}
	}
	return nil
}
