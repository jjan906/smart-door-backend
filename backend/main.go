package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"smart-door-backend/constant"
	"smart-door-backend/database"
	"smart-door-backend/handlers"
	mqttclient "smart-door-backend/mqtt"
)

func main() {
	// =====================================================
	// 1. INIT MONGODB
	// =====================================================
	if err := database.InitMongoDB(constant.MONGO_URI, constant.MONGO_DB); err != nil {
		log.Fatal("[FATAL] MongoDB:", err)
	}

	// =====================================================
	// 2. INIT MQTT (subscribe ke semua topic ESP32)
	// =====================================================
	if err := mqttclient.InitMQTT(); err != nil {
		log.Fatal("[FATAL] MQTT:", err)
	}

	// =====================================================
	// 3. SETUP FIBER WEB SERVER (untuk akses Frontend, publik)
	// =====================================================
	app := fiber.New(fiber.Config{
		AppName: "Smart Door Lock Backend",
	})


	// CORS - izinkan akses dari domain manapun (untuk frontend publik)
	// SESUDAH
	app.Use(cors.New(cors.Config{
    AllowOrigins:     "https://smart-door-frontend.vercel.app",
    AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
    AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
    AllowCredentials: false,
	}))

	app.Use(logger.New())

	// =====================================================
	// ROUTES
	// =====================================================
	api := app.Group("/api")

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "Smart Door Lock Backend",
		})
	})

	// Data sensor INA219
	api.Get("/sensor", handlers.GetSensorData)
	api.Get("/sensor/latest", handlers.GetLatestSensorData)
	api.Get("/access", handlers.GetAccessLogs)
	api.Get("/access/stats", handlers.GetAccessStats)

	// ← TAMBAH INI
	api.Delete("/access/all", handlers.DeleteAllAccessLogs)
	api.Delete("/access/:id", handlers.DeleteAccessLog)

	// Status pintu & perangkat
	api.Get("/door/status", handlers.GetDoorStatus)
	api.Get("/device/status", handlers.GetDeviceStatus)

	// Kontrol pintu dari Frontend (kirim perintah ke ESP32 via MQTT)
	api.Post("/door/open", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishDoorCommand("open"); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Perintah buka pintu terkirim"})
	})

	api.Post("/door/close", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishDoorCommand("close"); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Perintah tutup pintu terkirim"})
	})

	api.Post("/device/reset", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishReset(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Perintah reset terkirim"})
	})

	// =====================================================
	// START SERVER
	// =====================================================
	addr := fmt.Sprintf(":%s", constant.SERVER_PORT)
	fmt.Printf("[SERVER] Listening on %s\n", addr)
	log.Fatal(app.Listen(addr))
}
