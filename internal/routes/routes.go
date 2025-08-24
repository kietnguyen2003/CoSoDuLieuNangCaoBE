package routes

import (
	"database/sql"

	"clinic-management/internal/handlers"
	"clinic-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	api := router.Group("/api/v1")

	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)
	clinicHandler := handlers.NewClinicHandler(db)
	appointmentHandler := handlers.NewAppointmentHandler(db)
	medicalRecordHandler := handlers.NewMedicalRecordHandler(db)

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		users := protected.Group("/users")
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.PUT("/password", userHandler.ChangePassword)
		}

		clinics := protected.Group("/clinics")
		{
			clinics.GET("", clinicHandler.GetClinics)
			clinics.GET("/:id", clinicHandler.GetClinic)
			clinics.GET("/:id/doctors", clinicHandler.GetDoctors)
			clinics.GET("/:id/schedules", clinicHandler.GetSchedules)
		}

		appointments := protected.Group("/appointments")
		{
			appointments.GET("", appointmentHandler.GetAppointments)
			appointments.POST("", appointmentHandler.CreateAppointment)
			appointments.GET("/:id", appointmentHandler.GetAppointment)
			appointments.PUT("/:id", appointmentHandler.UpdateAppointment)
			appointments.DELETE("/:id", appointmentHandler.CancelAppointment)
		}

		medicalRecords := protected.Group("/medical-records")
		{
			medicalRecords.GET("", medicalRecordHandler.GetMedicalRecords)
			medicalRecords.GET("/:id", medicalRecordHandler.GetMedicalRecord)
			medicalRecords.POST("", medicalRecordHandler.CreateMedicalRecord)
			medicalRecords.PUT("/:id", medicalRecordHandler.UpdateMedicalRecord)
		}
	}
}