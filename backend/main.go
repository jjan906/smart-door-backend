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
	api.Get("/sensor", handlers.GetSensorData)
	api.Get("/sensor/latest", handlers.GetLatestSensorData)

	api.Get("/access", handlers.GetAccessLogs)
	api.Get("/access/stats", handlers.GetAccessStats)
	api.Delete("/access/all", handlers.DeleteAllAccessLogs)
	api.Delete("/access/:id", handlers.DeleteAccessLog)

	api.Get("/door/status", handlers.GetDoorStatus)
	api.Get("/device/status", handlers.GetDeviceStatus)

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

	api.Post("/buzzer/on", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishBuzzerCommand("on"); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Buzzer dinyalakan"})
	})

	api.Post("/buzzer/off", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishBuzzerCommand("off"); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Buzzer dimatikan"})
	})

	api.Post("/device/reset", func(c *fiber.Ctx) error {
		if err := mqttclient.PublishReset(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Perintah reset terkirim"})
	})

	// ── START SERVER ───────────────────────────────────
	addr := fmt.Sprintf(":%s", constant.SERVER_PORT)
	fmt.Printf("[SERVER] Listening on %s\n", addr)
	log.Fatal(app.Listen(addr))
}