package controllers

import (
	"errors"
	"hospital-system/models"
	"hospital-system/resource"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDoctors(ctx *gin.Context) {
	doctors, err := resource.DoctorService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, doctors)
}

func GetDoctor(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	doctor, err := resource.DoctorService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, doctor)
}

func CreateDoctor(ctx *gin.Context) {
	var doctor models.Doctor
	if err := ctx.ShouldBindJSON(&doctor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateDoctorDiseases(ctx, doctor.Diseases); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := resource.DoctorService.Create(ctx, &doctor); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, doctor)
}

func UpdateDoctor(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	var doctor models.Doctor
	if err := ctx.ShouldBindJSON(&doctor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateDoctorDiseases(ctx, doctor.Diseases); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := resource.DoctorService.Update(ctx, id, &doctor); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, doctor)
}

func DeleteDoctor(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	if err := resource.DoctorService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Doctor deleted successfully"})
}

func validateDoctorDiseases(ctx *gin.Context, diseaseIDs []string) error {
	if len(diseaseIDs) < 1 || len(diseaseIDs) > 3 {
		return errors.New("diseases length must be between 1 and 3")
	}

	diseases, err := resource.DiseaseService.GetAll(ctx)
	if err != nil {
		return err
	}

	existing := make(map[string]struct{}, len(diseases))
	for _, d := range diseases {
		existing[d.ID] = struct{}{}
	}

	for _, id := range diseaseIDs {
		if _, ok := existing[id]; !ok {
			return errors.New("disease id not found: " + id)
		}
	}

	return nil
}
