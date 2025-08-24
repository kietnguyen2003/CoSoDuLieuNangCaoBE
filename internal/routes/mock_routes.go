package routes

import (
	"clinic-management/internal/handlers"
	"clinic-management/internal/middleware"
	"clinic-management/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupMockRoutes(router *gin.Engine, mockService *services.MockDataService) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	api := router.Group("/api/v1")

	authHandler := handlers.NewMockAuthHandler(mockService)
	userHandler := handlers.NewMockUserHandler(mockService)
	clinicHandler := handlers.NewMockClinicHandler(mockService)
	appointmentHandler := handlers.NewMockAppointmentHandler(mockService)
	medicalRecordHandler := handlers.NewMockMedicalRecordHandler(mockService)
	scheduleHandler := handlers.NewMockScheduleHandler(mockService)

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

		schedules := protected.Group("/schedules")
		{
			schedules.GET("", scheduleHandler.GetSchedules)
			schedules.POST("", scheduleHandler.CreateSchedule)
			schedules.GET("/:id", scheduleHandler.GetSchedule)
			schedules.PUT("/:id", scheduleHandler.UpdateSchedule)
			schedules.DELETE("/:id", scheduleHandler.DeleteSchedule)
			schedules.POST("/assign", scheduleHandler.AssignDoctor)
			schedules.PUT("/:id/assign", scheduleHandler.ReassignDoctor)
			schedules.GET("/conflicts", scheduleHandler.CheckConflicts)
		}

		doctors := protected.Group("/doctors")
		{
			doctors.GET("/:id/schedules", scheduleHandler.GetDoctorSchedules)
		}
	}
}