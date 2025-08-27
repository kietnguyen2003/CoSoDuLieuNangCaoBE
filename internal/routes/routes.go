package routes

import (
	"database/sql"

	"clinic-management/internal/handlers"
	"clinic-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB, jwtSecret string) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	api := router.Group("/api/v1")

	authHandler := handlers.NewAuthHandler(db, jwtSecret)
	userHandler := handlers.NewUserHandler(db)
	clinicHandler := handlers.NewClinicHandler(db)
	appointmentHandler := handlers.NewAppointmentHandler(db)
	medicalRecordHandler := handlers.NewMedicalRecordHandler(db)
	prescriptionHandler := handlers.NewPrescriptionHandler(db)
	customerHandler := handlers.NewCustomerHandler(db)
	// labTestHandler := handlers.NewLabTestHandler(db)
	scheduleHandler := handlers.NewScheduleHandler(db)

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)       //check
		auth.POST("/register", authHandler.Register) //check
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
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
			clinics.GET("/specialties", clinicHandler.GetSpecialties)
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

		prescriptions := protected.Group("/prescriptions")
		{
			prescriptions.GET("", prescriptionHandler.GetPrescriptions)
			prescriptions.GET("/:id", prescriptionHandler.GetPrescription)
			prescriptions.POST("", prescriptionHandler.CreatePrescription)
			prescriptions.PUT("/:id", prescriptionHandler.UpdatePrescription)
		}

		medications := protected.Group("/medications")
		{
			medications.GET("", prescriptionHandler.GetMedications)
		}

		customers := protected.Group("/customers")
		{
			customers.GET("", customerHandler.GetCustomers)
			customers.GET("/:id", customerHandler.GetCustomer)
			customers.POST("", customerHandler.CreateCustomer)
		}

		// labTests := protected.Group("/lab-tests")
		// {
		// 	labTests.GET("", labTestHandler.GetLabTests)
		// 	labTests.GET("/:id", labTestHandler.GetLabTest)
		// 	labTests.POST("", labTestHandler.CreateLabOrder)
		// 	labTests.PUT("/:id", labTestHandler.UpdateLabTest)
		// 	labTests.DELETE("/:id", labTestHandler.DeleteLabTest)
		// }

		// labTestTypes := protected.Group("/lab-test-types")
		// {
		// 	labTestTypes.GET("", labTestHandler.GetLabTestTypes)
		// }

		schedules := protected.Group("/schedules")
		{
			schedules.GET("", scheduleHandler.GetSchedules)
			schedules.GET("/:id", scheduleHandler.GetSchedule)
			schedules.POST("", scheduleHandler.CreateSchedule)
			schedules.PUT("/:id", scheduleHandler.UpdateSchedule)
			schedules.DELETE("/:id", scheduleHandler.DeleteSchedule)
		}
	}
}
