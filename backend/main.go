package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"smart-door-backend/constant"
	"smart-door-backend/database"
	"smart-door-backend/handlers"
	mqttclient "smart-door-backend/mqtt"
)

// Ambil username dari Bearer token
func getActorFromToken(c *fiber.Ctx) string {
	auth := c.Get("Authorization")
	if len(auth) < 8 { return "unknown" }
	decoded, err := base64.StdEncoding.DecodeString(auth[7:])
	if err != nil { return "unknown" }
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) == 0 { return "unknown" }
	return parts[0]
}

func main() {
	// 1. INIT MONGODB
	if err := database.InitMongoDB(constant.MONGO_URI, constant.MONGO_DB); err != nil {
		log.Fatal("[FATAL] MongoDB:", err)
	}

	// 2. INIT MQTT
	if err := mqttclient.InitMQTT(); err != nil {
		log.Fatal("[FATAL] MQTT:", err)
	}

	// 3. SETUP FIBER
	app := fiber.New(fiber.Config{
		AppName: "Smart Door Lock Backend",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://smart-door-frontend.vercel.app",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: false,
	}))

	app.Use(logger.New())

	// ── HEALTH CHECK (publik) ──────────────────────────
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "Smart Door Lock Backend",
		})
	})

	api := app.Group("/api")

	// ── LOGIN (publik) ─────────────────────────────────
	api.Post("/login", handlers.Login)

	// ── AUTH MIDDLEWARE ────────────────────────────────
	api.Use(func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if len(auth) < 8 || auth[:7] != "Bearer " {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak ada atau tidak valid"})
		}
		token := auth[7:]
		decoded, err := base64.StdEncoding.DecodeString(token)
		if err != nil || !strings.HasPrefix(string(decoded), constant.ADMIN_USERNAME+":") {
			return c.Status(401).JSON(fiber.Map{"error": "Token tidak valid"})
		}
		return c.Next()
	})

	// ── PROTECTED ROUTES ───────────────────────────────

	// Auth
	api.Post("/auth/change-password", handlers.ChangePassword)

	// Sensor
	api.Get("/sensor", handlers.GetSensorData)
	api.Get("/sensor/latest", handlers.GetLatestSensorData)

	// Access logs
	api.Get("/access", handlers.GetAccessLogs)
	api.Get("/access/stats", handlers.GetAccessStats)
	api.Delete("/access/all", handlers.DeleteAllAccessLogs)
	api.Delete("/access/:id", handlers.DeleteAccessLog)

	// Door & device status
	api.Get("/door/status", handlers.GetDoorStatus)
	api.Get("/device/status", handlers.GetDeviceStatus)

	// Device management
	api.Get("/devices", handlers.GetDevices)
	api.Delete("/devices/:id", handlers.DeleteDevice)

	// Event logs
	api.Get("/logs", handlers.GetEventLogs)

	// Door control
	api.Post("/door/open", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishDoorCommand("open"); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		go handlers.SaveEventLog("door_open", getActorFromToken(c), "Perintah buka pintu", c.IP(), "success")
		return c.JSON(fiber.Map{"message": "Perintah buka pintu terkirim"})
	})

	api.Post("/door/close", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishDoorCommand("close"); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		go handlers.SaveEventLog("door_close", getActorFromToken(c), "Perintah tutup pintu", c.IP(), "success")
		return c.JSON(fiber.Map{"message": "Perintah tutup pintu terkirim"})
	})

	// Buzzer control
	api.Post("/buzzer/on", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishBuzzerCommand("on"); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		go handlers.SaveEventLog("buzzer_on", getActorFromToken(c), "Buzzer dinyalakan", c.IP(), "success")
		return c.JSON(fiber.Map{"message": "Buzzer dinyalakan"})
	})

	api.Post("/buzzer/off", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishBuzzerCommand("off"); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		go handlers.SaveEventLog("buzzer_off", getActorFromToken(c), "Buzzer dimatikan", c.IP(), "success")
		return c.JSON(fiber.Map{"message": "Buzzer dimatikan"})
	})

	// Device reset
	api.Post("/device/reset", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishReset(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		go handlers.SaveEventLog("device_reset", getActorFromToken(c), "Perintah reset ESP32", c.IP(), "success")
		return c.JSON(fiber.Map{"message": "Perintah reset terkirim"})
	})

	// ── START SERVER ───────────────────────────────────
	addr := fmt.Sprintf(":%s", constant.SERVER_PORT)
	fmt.Printf("[SERVER] Listening on %s\n", addr)
	log.Fatal(app.Listen(addr))
}