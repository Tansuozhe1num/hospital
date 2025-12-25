package controllers

import (
	"hospital-system/models"
	"hospital-system/resource"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDiseases(ctx *gin.Context) {
	diseases, err := resource.DiseaseService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, diseases)
}

func GetDisease(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	disease, err := resource.DiseaseService.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, disease)
}

func CreateDisease(ctx *gin.Context) {
	var disease models.Disease
	if err := ctx.ShouldBindJSON(&disease); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := resource.DiseaseService.Create(ctx, &disease); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, disease)
}

func UpdateDisease(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	var disease models.Disease
	if err := ctx.ShouldBindJSON(&disease); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := resource.DiseaseService.Update(ctx, id, &disease); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, disease)
}

func DeleteDisease(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		id = ctx.Query("id")
	}

	if err := resource.DiseaseService.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Disease deleted successfully"})
}
