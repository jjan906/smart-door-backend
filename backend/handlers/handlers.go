package handlers

import (
	"context"
	"encoding/base64"   // ← tambah
	"fmt"               // ← tambah
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"smart-door-backend/constant"   // ← tambah
	"smart-door-backend/database"
	data "smart-door-backend/Types"
)

const queryTimeout = 10 * time.Second

// =====================================================
// SAVE HANDLERS (dipanggil dari MQTT message handler)
// =====================================================

func SaveSensorData(sensor data.SensorData) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	sensor.ReceivedAt = time.Now().Format(time.RFC3339)
	collection := database.GetCollection("sensor_data")
	_, err := collection.InsertOne(ctx, sensor)
	return err
}

func SaveAccessLog(access data.AccessLog) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	access.ReceivedAt = time.Now().Format(time.RFC3339)
	collection := database.GetCollection("access_logs")
	_, err := collection.InsertOne(ctx, access)
	return err
}

func SaveDoorStatus(door data.DoorStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	door.ReceivedAt = time.Now().Format(time.RFC3339)
	collection := database.GetCollection("door_status")
	_, err := collection.InsertOne(ctx, door)
	return err
}

func SaveDeviceStatus(status data.DeviceStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	status.ReceivedAt = time.Now().Format(time.RFC3339)
	collection := database.GetCollection("device_status")
	_, err := collection.InsertOne(ctx, status)
	return err
}

// =====================================================
// HTTP GET HANDLERS (untuk Frontend)
// =====================================================

// GET /api/sensor?limit=50
func GetSensorData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	limit := c.QueryInt("limit", 50)
	collection := database.GetCollection("sensor_data")

	opts := options.Find().SetSort(bson.D{{Key: "received_at", Value: -1}}).SetLimit(int64(limit))
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer cursor.Close(ctx)

	var results []data.SensorData
	if err := cursor.All(ctx, &results); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// GET /api/sensor/latest - data sensor terbaru saja
func GetLatestSensorData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	collection := database.GetCollection("sensor_data")
	opts := options.FindOne().SetSort(bson.D{{Key: "received_at", Value: -1}})

	var result data.SensorData
	err := collection.FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Belum ada data sensor"})
	}

	return c.JSON(result)
}

// GET /api/access?limit=50
func GetAccessLogs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	limit := c.QueryInt("limit", 50)
	collection := database.GetCollection("access_logs")

	opts := options.Find().SetSort(bson.D{{Key: "received_at", Value: -1}}).SetLimit(int64(limit))
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer cursor.Close(ctx)

	var results []data.AccessLog
	if err := cursor.All(ctx, &results); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// GET /api/door/status - status pintu terkini
func GetDoorStatus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	collection := database.GetCollection("door_status")
	opts := options.FindOne().SetSort(bson.D{{Key: "received_at", Value: -1}})

	var result data.DoorStatus
	err := collection.FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Belum ada data status pintu"})
	}

	return c.JSON(result)
}

// GET /api/device/status - status perangkat ESP32 terkini
func GetDeviceStatus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	collection := database.GetCollection("device_status")
	opts := options.FindOne().SetSort(bson.D{{Key: "received_at", Value: -1}})

	var result data.DeviceStatus
	err := collection.FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Belum ada data status perangkat"})
	}

	return c.JSON(result)
}

// GET /api/access/stats - statistik granted vs denied
func GetAccessStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	collection := database.GetCollection("access_logs")

	granted, _ := collection.CountDocuments(ctx, bson.M{"status": "granted"})
	denied, _ := collection.CountDocuments(ctx, bson.M{"status": "denied"})

	return c.JSON(fiber.Map{
		"granted": granted,
		"denied":  denied,
		"total":   granted + denied,
	})
}

// DELETE /api/access/all - hapus semua log akses
func DeleteAllAccessLogs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	collection := database.GetCollection("access_logs")
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Semua log berhasil dihapus"})
}

// DELETE /api/access/:id - hapus satu log akses
func DeleteAccessLog(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	collection := database.GetCollection("access_logs")
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Log berhasil dihapus"})
}

// =====================================================
// AUTH HANDLER
// =====================================================

// POST /api/login
// Body: { "username": "...", "password": "..." }
func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}

	// Credentials hardcoded — ganti sesuai kebutuhan
	validUser := constant.ADMIN_USERNAME
	validPass := constant.ADMIN_PASSWORD

	if req.Username != validUser || req.Password != validPass {
		return c.Status(401).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	// Token sederhana: base64(username:timestamp)
	// Cukup untuk demo — bukan JWT produksi
	raw := fmt.Sprintf("%s:%d", req.Username, time.Now().Unix())
	token := base64.StdEncoding.EncodeToString([]byte(raw))

	return c.JSON(fiber.Map{
		"token":   token,
		"message": "Login berhasil",
	})
}