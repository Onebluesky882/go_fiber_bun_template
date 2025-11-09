```go
Constructor คืออะไร

ใน Go constructor เป็นฟังก์ชันที่สร้าง instance ของ struct พร้อมกำหนดค่าเริ่มต้นให้ โดย convention จะใช้

ประโยชน์:

ทำให้คุณไม่ต้องสร้าง struct แล้วตั้งค่า field ด้วยตัวเองทุกครั้ง

สามารถใช้ dependency injection กับ database connection ได้ง่าย




```

## Repository

\*\* สรุป หน้าที่ติดต่อdatabase

หน้าที่หลัก: จัดการเรื่อง data access หรือ database operations

ทำงานใกล้กับฐานข้อมูล เช่น query, insert, update, delete

ไม่รู้จัก HTTP หรือ API

มักมี constructor เพื่อรับ dependency (เช่น \*bun.DB)

## 2️ Handler

หน้าที่หลัก: รับผิดชอบเรื่อง HTTP / API layer

ทำงานใกล้กับ framework (Fiber ในกรณีนี้)

เรียก service หรือ repository เพื่อดึงข้อมูล แล้วแปลงเป็น JSON หรือ HTTP response

## Service (เชื่อมกลาง)

หน้าที่หลัก: business logic, ประมวลผลข้อมูล, เรียก repository

ทำให้ Handler ไม่ต้องรู้รายละเอียด database

```go

 HTTP Request (Fiber)
      ↓
Handler (รับ HTTP, แปลง request → call service, return JSON)
      ↓
Service (business logic, ประมวลผลข้อมูล)
      ↓
Repository (query DB)
      ↓
Database

```
