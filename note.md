📦 Repository คือที่จัดการกับฐานข้อมูลโดยตรง
เช่น ดึงข้อมูลจาก PostgreSQL, INSERT, UPDATE, DELETE

ทำก่อนเพราะ:

Service และ Handler จะ “เรียกใช้” มัน

ต้องแน่ใจก่อนว่า query ของเราทำงานได้จริง

ตัวอย่าง: internal/user/repository.go

Repository → Service → Handler
(ชั้นล่าง) (ตรรกะธุรกิจ) (ชั้น API / Fiber)

/\*
NewRepository เป็น constructor function

ใช้สร้าง object ของ Repository
โดยต้องส่ง \*bun.DB เข้ามา (connection ที่สร้างไว้แล้ว)
แล้ว return pointer ของ struct Repository

พูดง่าย ๆ:
→ มันคือ "ฟังก์ชันที่สร้าง repository พร้อมกับเชื่อมต่อฐานข้อมูลให้พร้อมใช้"
\*/

//[ Fiber Handler ] → เรียก → [ Repository ] → ใช้ → [ Bun DB Connection ] → PostgreSQL

```go
1️⃣ เริ่มจาก Repository (ชั้นล่างสุด)

package user

import (
	"context"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.NewSelect().Model(&users).Scan(ctx)
	return users, err
}




2️⃣ ต่อด้วย Service (ชั้นกลาง)

⚙️ Service คือ layer สำหรับ “ตรรกะทางธุรกิจ”
— เช่น validate input, ตรวจสิทธิ์, คำนวณ logic ก่อนหรือหลัง query

ทำต่อเพราะ:

มันต้องเรียก Repository ที่คุณเพิ่งสร้างไว้

Service สามารถรวมหลาย Repository ได้ (เช่น user + order)

ตัวอย่าง: internal/user/service.go
package user

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllUsers(ctx context.Context) ([]User, error) {
	return s.repo.GetAll(ctx)
}




3️⃣ สุดท้าย Handler (ชั้นบนสุด / API Layer)

Handler คือส่วนที่เชื่อมกับ framework (เช่น Fiber)
— รับ request, เรียก service, แล้วส่ง response กลับ client

ทำสุดท้ายเพราะ:

ต้องมี service ก่อนถึงจะ inject เข้า handler ได้

เป็นจุดเชื่อมต่อกับ REST API หรือ WebSocket

ตัวอย่าง: internal/user/handler.go

package user

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

```
