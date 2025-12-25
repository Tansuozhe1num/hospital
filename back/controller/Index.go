package controllers

import (
	"github.com/gin-gonic/gin"
	"hospital-system/server/baseinfo"
)

func Index(ctx *gin.Context) {
	info, err := baseinfo.GetIndexInfo(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, info)
}
