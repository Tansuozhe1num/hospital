package httpserver

import (
	controllers "hospital-system/controller"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// 患者管理
	patientGroup := router.Group("/api/patients")
	{
		patientGroup.GET("/getPatients", controllers.GetPatients)
		patientGroup.GET("/getPatient", controllers.GetPatient)
		patientGroup.POST("/createPatient", controllers.CreatePatient)
		patientGroup.PUT("/updatePatient", controllers.UpdatePatient)
		patientGroup.DELETE("/deletePatient", controllers.DeletePatient)
	}

	//// 病种管理API
	//diseaseGroup := router.Group("/api/diseases")
	//{
	//	diseaseGroup.GET("/", getDiseases)
	//	diseaseGroup.GET("/:id", getDisease)
	//	diseaseGroup.POST("/", createDisease)
	//	diseaseGroup.PUT("/:id", updateDisease)
	//	diseaseGroup.DELETE("/:id", deleteDisease)
	//}
	//
	//// 医生管理API
	//doctorGroup := router.Group("/api/doctors")
	//{
	//	doctorGroup.GET("/", getDoctors)
	//	doctorGroup.GET("/:id", getDoctor)
	//	doctorGroup.POST("/", createDoctor)
	//	doctorGroup.PUT("/:id", updateDoctor)
	//	doctorGroup.DELETE("/:id", deleteDoctor)
	//	doctorGroup.GET("/:id/diseases", getDoctorDiseases)
	//}
	//
	//// 挂号管理API
	//registrationGroup := router.Group("/api/registrations")
	//{
	//	registrationGroup.GET("/", getRegistrations)
	//	registrationGroup.GET("/:id", getRegistration)
	//	registrationGroup.POST("/", createRegistration)
	//	registrationGroup.PUT("/:id", updateRegistration)
	//	registrationGroup.DELETE("/:id", deleteRegistration)
	//	registrationGroup.GET("/patient/:patientId", getPatientRegistrations)
	//}
}

func SetUpFronted(router *gin.Engine) {
	frontendDir := filepath.Clean("../fronter")
	router.Static("/css", filepath.Join(frontendDir, "css"))
	router.Static("/js", filepath.Join(frontendDir, "js"))
	router.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(frontendDir, "index.html"))
	})
	router.GET("/index.html", func(c *gin.Context) {
		c.File(filepath.Join(frontendDir, "index.html"))
	})
	router.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.File(filepath.Join(frontendDir, "index.html"))
	})
}
