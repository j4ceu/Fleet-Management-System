# 🚍 Fleet Management System

Fleet Management System adalah aplikasi backend yang memantau pergerakan kendaraan, menyimpan histori lokasi, dan mengirim notifikasi ketika kendaraan memasuki area geofence tertentu. Sistem ini dirancang untuk scalable dan real-time.

---

## ⚙️ Teknologi & Tools

| Komponen           | Teknologi                            |
|--------------------|---------------------------------------|
| Bahasa             | Go (Golang)                           |
| Framework Web      | Gin                                   |
| DB                 | PostgreSQL + pgx driver               |
| Message Broker     | MQTT (Eclipse Mosquitto), RabbitMQ    |
| Architecture       | Clean Architecture                    |
| Containerization   | Docker + Docker Compose               |
| Logging            | Logrus                                |

---

## 🔧 Cara Menjalankan 

### ✅ 1. Siapkan Environment
- Install **Docker Desktop** (pastikan Docker Engine menyala)
- Clone project ini
- Jalankan command:

```bash
git clone <repo>
cd Fleet-Management-System
docker-compose up -d --build
```

> Pertama kali build akan otomatis:
> - Setup DB PostgreSQL
> - Setup MQTT Mosquitto
> - Setup RabbitMQ
> - Build binary Go app, publisher, worker

### 🧪 2. Cek Service
- RabbitMQ UI: [http://localhost:15672](http://localhost:15672) (user: `guest`, pass: `guest`)
- API Endpoint: `http://localhost:8080/api/v1/vehicles/:id/location`

### 🛰️ 3. Kirim Simulasi Lokasi
```bash
docker logs -f mqtt-publisher
```
> Publisher akan mengirim data GPS ke MQTT tiap 5 detik.

### 📥 4. Cek Worker Alert
```bash
docker logs -f geofence-alert-worker
```
> Akan muncul alert jika kendaraan masuk geofence.

---

---

## 📄 Postman Collection

Tersedia Postman Collection di file `postman_collection.json` untuk mencoba seluruh endpoint.

### 🔧 Cara Menjalankan:
1. Buka Postman
2. Import file `postman_collection.json`
3. Pastikan environment base URL adalah `http://localhost:8080`
4. Jalankan koleksi seperti `Get LAtest Location`, atau `Get Location History`
5. Untuk vehicle_id yg tersedia yaitu B1234XYZ, D5678DGV, F9876CMP, T1111CBC, B2123KCM, A5230KCM

---

## 🗺️ Alur Awal (Soal)
- Kendaraan publish koordinat via MQTT
- Data disimpan ke DB
- Jika kendaraan masuk geofence → kirim notifikasi ke RabbitMQ

---

## 🔄 Alur Yang Dibangun (Final Flow)

```text
[MQTT Publisher]
     |
     v
[MQTT Subscriber]
     |
     v
[Publish to Exchange: fleet.events]
     |
     v
[Queue: save_db]
     | 
     | 
     +--> [Check Geofence]   <-- langsung saat terima
     |         |
     |         v
     |  [If inside geofence]
     |         |
     |  [Publish to Queue: geofence_alerts]
     |         | 
     |  [Worker Consumer : Print Alert]
     |
     v
[Worker: Consume from save_db]
     |
     v
[Store to DB]
```


✅ Semua komunikasi antar service menggunakan message queue untuk performa & decoupling.

---

## 📍 Geofencing & Queue Strategy

Untuk mengecek apakah kendaraan berada dalam area geofence, digunakan pendekatan lingkaran berbasis jarak. Fungsi `Haversine` digunakan untuk menghitung jarak dua titik berdasarkan koordinat GPS. Untuk dummynya routenya ada pada `route.csv` yang saya generate berdasarkan route dari Summarecon Mal Bekasi ke La Terazza

---

## 📥 Kenapa Pakai Queue Sebelum Masuk DB?

Sebelum data disimpan ke database, sistem ini menggunakan **RabbitMQ queue `save_db`** sebagai perantara.

### ✅ Keuntungan:
- **Menghindari ledakan request** langsung ke database dari ratusan kendaraan
- Bisa kasih **buffer** kalau DB lambat atau down sementara
- Mudah untuk **scale out** worker penyimpan ke DB
- Bisa **retry** saat terjadi error

### ⚠️ Risiko / Hal yang Perlu Diperhatikan:
- Butuh monitoring queue agar tidak overload
- Harus pastikan data di queue diproses tepat waktu

Pattern ini disebut **event-driven ingestion** — maksudnya, data disimpan setelah lewat antrian (queue), bukan langsung ke DB.

Tujuannya supaya sistem tetap stabil walaupun data masuk terus-menerus.

---

## 📌 Fitur Utama

- 🔄 Ingest data lokasi via MQTT
- 🛢️ Simpan lokasi ke PostgreSQL
- 📍 Deteksi geofence menggunakan Haversine
- 📡 Publish notifikasi ke RabbitMQ 
- 🔔 Worker untuk konsumsi alert & log
- 🧪 REST API untuk:
  - `GET /api/v1/vehicles/:vehicle_id/location`
  - `GET /api/v1//vehicles/:vehicle_id/history?start=&end=`

---

## 🧠 Apa yang Dipelajari

Sebagai backend engineer:

- ✅ **Pertama kali belajar MQTT** dan cara publish/subscribe menggunakan Eclipse Mosquitto
- ✅ Paham konsep **exchange di RabbitMQ**
- ✅ Belajar bagaimana cara **memisahkan flow data** ke ingestion dan alert secara paralel
- ✅ Ternyata **Docker Compose bisa langsung jalanin SQL** untuk inisialisasi database, cukup dengan menaruh file `.sql` di folder `initdb` dan mount ke container

---

## 😅 Tantangan Selama Pengerjaan

- 🔁 Sinkronisasi MQTT → RabbitMQ → DB agar tidak overload
- 🐳 Error permission Docker saat mount config `mosquitto.conf`
- 📦 Menyesuaikan container agar tidak campur log saat development
- 🧱 Harus build end-to-end pipeline (backend, broker, ingestion, publisher, worker)

---

## 🚀 Pengembangan Selanjutnya

- 📱 Integrasi notifikasi Telegram / Email saat kendaraan masuk geofence
- 🗺️ Visualisasi live GPS & geofence dengan frontend (React + Leaflet)
- 📊 Dashboard riwayat per kendaraan
- ⚠️ Retry / dead-letter queue untuk error handling
- 🔐 Auth & otorisasi API

---

## 👨‍💻 Penulis

Made with ❤️ by [j4ceu] – Backend engineer who keep learning
