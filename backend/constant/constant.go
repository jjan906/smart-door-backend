package constant

import "os"

// getEnv membaca environment variable, fallback ke default jika tidak ada
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

var (
	// MQTT HiveMQ Cloud config
	MQTT_URL      = getEnv("MQTT_URL", "45ed91ba68e640d880f44148cb38bae6.s1.eu.hivemq.cloud")
	MQTT_USERNAME = getEnv("MQTT_USERNAME", "esp32user")
	MQTT_PASSWORD = getEnv("MQTT_PASSWORD", "Esp32@IoT2025")
	MQTT_PORT     = getEnv("MQTT_PORT", "8883")

	// MongoDB config
	MONGO_URI = getEnv("MONGO_URI", "mongodb://localhost:27017")
	MONGO_DB  = getEnv("MONGO_DB", "smartDoorDB")

	// Server port
	SERVER_PORT = getEnv("PORT", "3000")

	// ← TAMBAH INI
	ADMIN_USERNAME = getEnv("ADMIN_USERNAME", "admin")
	ADMIN_PASSWORD = getEnv("ADMIN_PASSWORD", "Admin@IoT2025")
)

// MQTT Topics
const (
	TOPIC_SENSOR    = "smartdoor/sensor"
	TOPIC_ACCESS    = "smartdoor/access"
	TOPIC_STATUS    = "smartdoor/status"
	TOPIC_DOOR      = "smartdoor/door"
	TOPIC_CMD_DOOR  = "smartdoor/cmd/door"
	TOPIC_CMD_RESET = "smartdoor/cmd/reset"
	TOPIC_CMD_BUZZER = "smartdoor/cmd/buzzer" // ← tambah ini
)
