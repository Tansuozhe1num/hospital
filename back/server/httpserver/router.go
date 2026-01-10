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
		patientGroup.GET("/getPatients", auth.GinAuthMiddleware("admin", "doctor"), controllers.GetPatients)
		patientGroup.GET("/getPatient", auth.GinAuthMiddleware("admin", "doctor"), controllers.GetPatient)
		patientGroup.POST("/createPatient", auth.GinAuthMiddleware("admin"), controllers.CreatePatient)
		patientGroup.PUT("/updatePatient", auth.GinAuthMiddleware("admin"), controllers.UpdatePatient)
		patientGroup.DELETE("/deletePatient", auth.GinAuthMiddleware("admin"), controllers.DeletePatient)
	}

	diseaseGroup := router.Group("/api/diseases")
	{
		diseaseGroup.GET("/getDiseases", auth.GinAuthMiddleware("admin", "doctor", "patient"), controllers.GetDiseases)
		diseaseGroup.GET("/getDisease", auth.GinAuthMiddleware("admin", "doctor", "patient"), controllers.GetDisease)
		diseaseGroup.POST("/createDisease", auth.GinAuthMiddleware("admin", "doctor"), controllers.CreateDisease)
		diseaseGroup.PUT("/updateDisease", auth.GinAuthMiddleware("admin", "doctor"), controllers.UpdateDisease)
		diseaseGroup.DELETE("/deleteDisease", auth.GinAuthMiddleware("admin", "doctor"), controllers.DeleteDisease)
	}

	doctorGroup := router.Group("/api/doctors")
	{
		doctorGroup.GET("/getDoctors", auth.GinAuthMiddleware("admin", "doctor", "patient"), controllers.GetDoctors)
		doctorGroup.GET("/getDoctor", auth.GinAuthMiddleware("admin", "doctor", "patient"), controllers.GetDoctor)
		doctorGroup.POST("/createDoctor", auth.GinAuthMiddleware("admin"), controllers.CreateDoctor)
		doctorGroup.PUT("/updateDoctor", auth.GinAuthMiddleware("admin"), controllers.UpdateDoctor)
		doctorGroup.DELETE("/deleteDoctor", auth.GinAuthMiddleware("admin"), controllers.DeleteDoctor)
	}

	departmentGroup := router.Group("/api/departments")
	{
		departmentGroup.GET("/getDepartments", auth.GinAuthMiddleware("admin", "doctor", "patient"), controllers.GetDepartments)
		departmentGroup.GET("/getDepartment", auth.GinAuthMiddleware("admin", "doctor", "patient"), controllers.GetDepartment)
		departmentGroup.POST("/createDepartment", auth.GinAuthMiddleware("admin"), controllers.CreateDepartment)
		departmentGroup.PUT("/updateDepartment", auth.GinAuthMiddleware("admin"), controllers.UpdateDepartment)
		departmentGroup.DELETE("/deleteDepartment", auth.GinAuthMiddleware("admin"), controllers.DeleteDepartment)
	}

	registrationGroup := router.Group("/api/registrations")
	{
		registrationGroup.GET("/getRegistrations", auth.GinAuthMiddleware("admin", "doctor", "patient"), controllers.GetRegistrations)
		registrationGroup.GET("/getRegistration", auth.GinAuthMiddleware("admin", "doctor", "patient"), controllers.GetRegistration)
		registrationGroup.POST("/createRegistration", auth.GinAuthMiddleware("admin", "patient"), controllers.CreateRegistration)
		registrationGroup.PUT("/updateRegistration", auth.GinAuthMiddleware("admin", "doctor"), controllers.UpdateRegistration)
		registrationGroup.DELETE("/deleteRegistration", auth.GinAuthMiddleware("admin"), controllers.DeleteRegistration)
	}

	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", controllers.LoginOrRegister)
		authGroup.GET("/me", auth.GinAuthMiddleware(), controllers.GetMe)
		authGroup.POST("/assignDoctorAccount", auth.GinAuthMiddleware("admin"), controllers.AssignDoctorAccount)
		authGroup.POST("/upsertMyPatientProfile", auth.GinAuthMiddleware("patient"), controllers.UpsertMyPatientProfile)
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
