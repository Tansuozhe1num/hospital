package controllers

import (
	"net/http"
	"time"

	"hospital-system/auth"
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
