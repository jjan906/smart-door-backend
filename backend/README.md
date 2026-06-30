# Smart Door Lock - Backend

Backend Go untuk sistem Smart Door Lock IoT, menghubungkan ESP32 (via MQTT HiveMQ Cloud) dengan Frontend melalui REST API. Arsitektur: **Backend-Centric**.

## Arsitektur

```
ESP32 --(MQTT)--> HiveMQ Cloud --(MQTT)--> Backend Go --(REST API)--> Frontend
                                                |
                                                v
                                          MongoDB Atlas
```

## Struktur Proyek

```
├── main.go              # Entry point, setup Fiber + routes
├── Types/sensor.go       # Struct data (Sensor, Access, Door, Status)
├── constant/constant.go  # Konfigurasi (env variable)
├── database/database.go  # Koneksi MongoDB
├── handlers/handlers.go  # Simpan & ambil data (CRUD)
├── mqtt/mqtt.go           # MQTT subscribe & publish
└── Dockerfile
```

## Teknologi

- Go + Fiber (web framework)
- Eclipse Paho MQTT
- MongoDB Atlas

## Menjalankan Lokal

1. Install dependencies:
```
go mod download
```

2. Copy `.env.example` ke `.env` dan isi kredensial:
```
cp .env.example .env
```

3. Jalankan:
```
go run main.go
```

Server berjalan di `http://localhost:3000`

## Endpoint API

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/` | Health check |
| GET | `/api/sensor?limit=50` | Histori data sensor INA219 |
| GET | `/api/sensor/latest` | Data sensor terbaru |
| GET | `/api/access?limit=50` | Histori log akses RFID |
| GET | `/api/access/stats` | Statistik granted/denied |
| GET | `/api/door/status` | Status pintu terkini |
| GET | `/api/device/status` | Status perangkat ESP32 |
| POST | `/api/door/open` | Buka pintu remote |
| POST | `/api/door/close` | Tutup pintu remote |
| POST | `/api/device/reset` | Reset ESP32 |

## Deploy ke Railway (Publik)

Lihat panduan lengkap di chat / dokumen terpisah.
