package controllers

import (
	"hospital-system/models"
	"hospital-system/resource"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPatients(ctx *gin.Context) {
	patients, err := resource.PatientService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, patients)
}

func GetPatient(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	patient, err := resource.PatientService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, patient)
}

func CreatePatient(ctx *gin.Context) {
	var patient models.Patient
	if err := ctx.ShouldBindJSON(&patient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := resource.PatientService.Create(ctx, &patient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, patient)
}

func UpdatePatient(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	var patient models.Patient
	if err := ctx.ShouldBindJSON(&patient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := resource.PatientService.Update(ctx, id, &patient); err != nil {
		if err.Error() == "patient not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, patient)
}

func DeletePatient(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	if err := resource.PatientService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}
