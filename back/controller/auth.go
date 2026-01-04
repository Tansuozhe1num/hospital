package controllers

import (
	"net/http"
	"time"

	"hospital-system/auth"
	"hospital-system/models"
	"hospital-system/resource"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type assignDoctorAccountRequest struct {
	DoctorID string `json:"doctorId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type upsertMyPatientProfileRequest struct {
	Name             string `json:"name"`
	Gender           string `json:"gender"`
	Age              int    `json:"age"`
	Phone            string `json:"phone"`
	IDCard           string `json:"idCard"`
	Address          string `json:"address"`
	EmergencyContact string `json:"emergencyContact"`
	EmergencyPhone   string `json:"emergencyPhone"`
}

func LoginOrRegister(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, isNew, err := resource.AccountService.LoginOrRegister(ctx, req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.IssueToken(account.ID, account.Role, 72*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"role":  account.Role,
		"id":    account.ID,
		"isNew": isNew,
	})
}

func GetMe(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, gin.H{
		"id":       account.ID,
		"username": account.Username,
		"role":     account.Role,
		"linkedId": account.LinkedID,
	})
}

func AssignDoctorAccount(ctx *gin.Context) {
	var req assignDoctorAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DoctorID == "" || req.Username == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "doctorId/username/password cannot be empty"})
		return
	}

	if _, err := resource.DoctorService.GetByID(ctx, req.DoctorID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "doctor not found"})
		return
	}

	account, err := resource.AccountService.UpsertDoctorAccount(ctx, req.DoctorID, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       account.ID,
		"username": account.Username,
		"role":     account.Role,
		"linkedId": account.LinkedID,
	})
}

func UpsertMyPatientProfile(ctx *gin.Context) {
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
	if account.Role != "patient" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req upsertMyPatientProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patient := models.Patient{
		Name:             req.Name,
		Gender:           req.Gender,
		Age:              req.Age,
		Phone:            req.Phone,
		IDCard:           req.IDCard,
		Address:          req.Address,
		EmergencyContact: req.EmergencyContact,
		EmergencyPhone:   req.EmergencyPhone,
	}

	if account.LinkedID != "" {
		if err := resource.PatientService.Update(ctx, account.LinkedID, &patient); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"patientId": account.LinkedID,
			"patient":   patient,
		})
		return
	}

	if err := resource.PatientService.Create(ctx, &patient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedAccount, err := resource.AccountService.SetLinkedID(ctx, account.ID, patient.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"patientId": updatedAccount.LinkedID,
		"patient":   patient,
	})
}
