# Go Fiber Test

โปรเจกต์ตัวอย่าง REST API ด้วย Go และ Fiber พร้อมฐานข้อมูล MySQL

## ภาพรวม

แอปพลิเคชันนี้ใช้:

- Go
- Fiber (Web framework)
- GORM (ORM)
- MySQL

แบ่งงานตามหลัก SRP ให้ชัดเจน:

- `main.go` สำหรับเริ่มต้นแอปและลงทะเบียน route
- `database/` สำหรับเชื่อมต่อและ migration
- `controllers/` สำหรับจัดการ request logic
- `models/` สำหรับโครงสร้างข้อมูลและ scopes
- `routes/` สำหรับจัดกลุ่ม route และ middleware

## การติดตั้ง

1. ติดตั้ง Go
2. เปิดเทอร์มินัลที่โฟลเดอร์โปรเจกต์
3. ติดตั้ง dependency

```powershell
cd c:\Users\praew_\OneDrive\Documents\go-fiber-test
go mod download
```

4. ตั้งค่า MySQL

ใช้ค่าเริ่มต้นใน `database/database.go`:

- user: `root`
- password: `` (ว่าง)
- host: `127.0.0.1`
- port: `3306`
- database: `golang_test`

## การรันโปรเจกต์

```powershell
go run main.go
```

จากนั้นเข้าถึง API ได้ที่ `http://localhost:3000`

## Route สำคัญ

### Auth

- `GET /api/v1/` - ต้องใช้ Basic Auth: `gofiber` / `21022566`
- `POST /api/v1/` - ต้องใช้ Basic Auth
- `GET /api/v1/user-params/:name` - ต้องใช้ Basic Auth
- `GET /api/v1/fact/:number` - ต้องใช้ Basic Auth
- `POST /api/v1/inet` - ต้องใช้ Basic Auth

### Dogs

- `GET /api/v1/dog` - ดึงข้อมูลสุนัขทั้งหมด
- `GET /api/v1/dog/filter?search=<id>` - หา dog ตาม dog_id
- `GET /api/v1/dog/json` - ผลลัพธ์ dog ในรูปแบบ JSON
- `GET /api/v1/dog/json-v2` - ผลลัพธ์ dog พร้อมสรุปสี
- `POST /api/v1/dog/seed` - สร้างข้อมูล dummy ของสุนัข
- `GET /api/v1/dog/deleted` - ดูสุนัขที่ถูกลบแล้ว
- `GET /api/v1/dog/range` - ดูสุนัขตามช่วง dog_id
- `POST /api/v1/dog` - เพิ่มสุนัขใหม่
- `PUT /api/v1/dog/:id` - อัปเดตสุนัข
- `DELETE /api/v1/dog/:id` - ลบสุนัข

### Users

- `GET /api/v1/user` - ดึงข้อมูลโปรไฟล์ผู้ใช้งาน
- `GET /api/v1/user/generations` - สรุปผู้ใช้งานตาม generation
- `GET /api/v1/user/search?search=<keyword>` - ค้นหาผู้ใช้งาน
- `POST /api/v1/user/seed` - สร้างข้อมูลโปรไฟล์ผู้ใช้งานแบบ dummy (Basic Auth: `testgo` / `23012023`)
- `POST /api/v1/user` - เพิ่มผู้ใช้งานใหม่ (Basic Auth)
- `PUT /api/v1/user/:id` - อัปเดตผู้ใช้งาน (Basic Auth)
- `DELETE /api/v1/user/:id` - ลบผู้ใช้งาน (Basic Auth)

### Company

- `GET /api/v1/company/` - ดึงข้อมูลบริษัททั้งหมด
- `GET /api/v1/company/:id` - ดึงข้อมูลบริษัทตาม ID
- `POST /api/v1/company/` - เพิ่มบริษัท
- `PUT /api/v1/company/:id` - อัปเดตบริษัท
- `DELETE /api/v1/company/:id` - ลบบริษัท

### API เวอร์ชันอื่น

- `GET /api/v2/` - เทสเวอร์ชัน 2
- `GET /api/v3/:name?tax_id=<tax_id>` - แปลง tax_id เป็น ASCII codes

## โครงสร้างโฟลเดอร์

- `main.go`
- `controllers/`
- `database/`
- `models/`
- `routes/`

## หมายเหตุ

ไฟล์นี้ช่วยให้เข้าใจ flow ของโปรเจกต์ และสะดวกสำหรับการทดสอบ API ด้วย Postman หรือ curl
