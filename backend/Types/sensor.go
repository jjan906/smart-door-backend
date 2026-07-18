package types

// SensorData - data dari INA219 (voltage, current, power)
type SensorData struct {
	DeviceID   string  `json:"device_id" bson:"device_id"`
	VoltageV   float64 `json:"voltage_V" bson:"voltage_V"`
	CurrentMA  float64 `json:"current_mA" bson:"current_mA"`
	PowerMW    float64 `json:"power_mW" bson:"power_mW"`
	ShuntMV    float64 `json:"shunt_mV" bson:"shunt_mV"`
	Timestamp  int64   `json:"timestamp" bson:"timestamp"`
	ReceivedAt string  `json:"received_at,omitempty" bson:"received_at,omitempty"`
}

// AccessLog - log akses RFID / button / remote
type AccessLog struct {
	DeviceID   string `json:"device_id" bson:"device_id"`
	UID        string `json:"uid" bson:"uid"`
	Name       string `json:"name" bson:"name"`
	Status     string `json:"status" bson:"status"` // granted / denied
	DoorOpen   bool   `json:"door_open" bson:"door_open"`
	Timestamp  int64  `json:"timestamp" bson:"timestamp"`
	ReceivedAt string `json:"received_at,omitempty" bson:"received_at,omitempty"`
}

// DoorStatus - status pintu locked/unlocked
type DoorStatus struct {
	DeviceID   string `json:"device_id" bson:"device_id"`
	Door       string `json:"door" bson:"door"` // locked / unlocked
	Timestamp  int64  `json:"timestamp" bson:"timestamp"`
	ReceivedAt string `json:"received_at,omitempty" bson:"received_at,omitempty"`
}

// DeviceStatus - status perangkat ESP32 (online/offline)
type DeviceStatus struct {
	DeviceID   string `json:"device_id" bson:"device_id"`
	State      string `json:"state" bson:"state"`
	DoorStatus string `json:"door_status" bson:"door_status"`
	WifiRSSI   int    `json:"wifi_rssi" bson:"wifi_rssi"`
	UptimeMS   int64  `json:"uptime_ms" bson:"uptime_ms"`
	IP         string `json:"ip" bson:"ip"`
	Timestamp  int64  `json:"timestamp" bson:"timestamp"`
	ReceivedAt string `json:"received_at,omitempty" bson:"received_at,omitempty"`
}

// DoorCommand - perintah buka/tutup pintu dari server ke ESP32
type DoorCommand struct {
	Command string `json:"command"` // "open" atau "close"
}

// Device - perangkat ESP32 yang terdaftar
type Device struct {
	DeviceID   string `json:"device_id" bson:"device_id"`
	Name       string `json:"name" bson:"name"`
	IP         string `json:"ip" bson:"ip"`
	State      string `json:"state" bson:"state"`       // online / offline
	DoorStatus string `json:"door_status" bson:"door_status"`
	WifiRSSI   int    `json:"wifi_rssi" bson:"wifi_rssi"`
	UptimeMS   int64  `json:"uptime_ms" bson:"uptime_ms"`
	FirstSeen  string `json:"first_seen" bson:"first_seen"`
	LastSeen   string `json:"last_seen" bson:"last_seen"`
}

// EventLog - pencatatan aktivitas sistem
type EventLog struct {
	Action     string `json:"action" bson:"action"`         // login, door_open, buzzer_on, dst
	Actor      string `json:"actor" bson:"actor"`           // username / device_id
	Detail     string `json:"detail" bson:"detail"`         // keterangan tambahan
	IP         string `json:"ip" bson:"ip"`                 // IP peminta
	Status     string `json:"status" bson:"status"`         // success / failed
	ReceivedAt string `json:"received_at" bson:"received_at"`
}