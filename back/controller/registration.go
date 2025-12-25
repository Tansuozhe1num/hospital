package controllers

import (
	"hospital-system/models"
	"hospital-system/resource"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRegistrations(ctx *gin.Context) {
	registrations, err := resource.RegistrationService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, registrations)
}

func GetRegistration(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	registration, err := resource.RegistrationService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, registration)
}

func CreateRegistration(ctx *gin.Context) {
	var registration models.Registration
	if err := ctx.ShouldBindJSON(&registration); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	var registration models.Registration
	if err := ctx.ShouldBindJSON(&registration); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
