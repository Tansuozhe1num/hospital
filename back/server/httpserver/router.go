package httpserver

import (
	controllers "hospital-system/controller"
	"net/http"
	"path/filepath"

	"hospital-system/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// 首页展示信息
	index := router.Group("/index")
	{
		index.GET("/", controllers.Index)
	}

	// 患者管理
	patientGroup := router.Group("/api/patients")
	{
		patientGroup.GET("/getPatients", controllers.GetPatients)
		patientGroup.GET("/getPatient", controllers.GetPatient)
		patientGroup.POST("/createPatient", controllers.CreatePatient)
		patientGroup.PUT("/updatePatient", controllers.UpdatePatient)
		patientGroup.DELETE("/deletePatient", controllers.DeletePatient)
	}

	diseaseGroup := router.Group("/api/diseases")
	{
		diseaseGroup.GET("/getDiseases", controllers.GetDiseases)
		diseaseGroup.GET("/getDisease", controllers.GetDisease)
		diseaseGroup.POST("/createDisease", controllers.CreateDisease)
		diseaseGroup.PUT("/updateDisease", controllers.UpdateDisease)
		diseaseGroup.DELETE("/deleteDisease", controllers.DeleteDisease)
	}

	doctorGroup := router.Group("/api/doctors")
	{
		doctorGroup.GET("/getDoctors", controllers.GetDoctors)
		doctorGroup.GET("/getDoctor", controllers.GetDoctor)
		doctorGroup.POST("/createDoctor", controllers.CreateDoctor)
		doctorGroup.PUT("/updateDoctor", controllers.UpdateDoctor)
		doctorGroup.DELETE("/deleteDoctor", controllers.DeleteDoctor)
	}

	registrationGroup := router.Group("/api/registrations")
	{
		registrationGroup.GET("/getRegistrations", controllers.GetRegistrations)
		registrationGroup.GET("/getRegistration", controllers.GetRegistration)
		registrationGroup.POST("/createRegistration", controllers.CreateRegistration)
		registrationGroup.PUT("/updateRegistration", controllers.UpdateRegistration)
		registrationGroup.DELETE("/deleteRegistration", controllers.DeleteRegistration)
	}

	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", controllers.LoginOrRegister)
		authGroup.GET("/me", auth.GinAuthMiddleware(), controllers.GetMe)
		authGroup.POST("/assignDoctorAccount", auth.GinAuthMiddleware("admin"), controllers.AssignDoctorAccount)
	}
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
