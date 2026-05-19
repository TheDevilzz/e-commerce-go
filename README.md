# ShopHub E-Commerce / ระบบร้านค้าออนไลน์ ShopHub

ShopHub is a full-stack e-commerce project with a React + TypeScript frontend and a Go + Gin + GORM backend API. The system includes product/category management, cart, wishlist, orders, payments, promotions, dashboard statistics, authentication, refresh tokens, logout revocation, and role-based admin protection.

ShopHub คือโปรเจกต์ร้านค้าออนไลน์แบบ full-stack ที่ใช้ React + TypeScript สำหรับ frontend และ Go + Gin + GORM สำหรับ backend API รองรับสินค้า หมวดหมู่ ตะกร้า wishlist คำสั่งซื้อ การชำระเงิน โปรโมชั่น dashboard ระบบ login/register, refresh token, logout revoke token และการป้องกันหน้า admin ด้วย role

## Project Structure / โครงสร้างโปรเจกต์

```text
.
├── src/                 # Frontend React source
├── public/              # Static frontend assets
├── backend-go/          # Go backend API
│   ├── cmd/             # Backend entrypoint
│   ├── config/          # Environment config
│   ├── database/        # MySQL connection
│   └── internal/        # Handlers, middleware, models, repositories, services, router
├── package.json         # Frontend scripts and dependencies
└── README.md
```

## Features / ฟีเจอร์หลัก

- Product browsing, category filters, product detail, cart, wishlist, checkout flow
- Admin dashboard for products, categories, orders, customers, promotions, inventory, and stats
- Login/register with bcrypt password hashing
- Access Token lifetime: 15 minutes
- Refresh Token lifetime: 7-30 days, default 30 days
- Refresh Token stored in database as SHA-256 hash
- Refresh API with token rotation
- Logout API revokes refresh token
- Role middleware for admin/user protected routes
- Secure config via environment variables, no hardcoded secrets
- HTTPS enforcement middleware for production, with localhost allowed for development

- ดูสินค้า กรองหมวดหมู่ รายละเอียดสินค้า ตะกร้า wishlist และ checkout
- หน้า admin สำหรับจัดการสินค้า หมวดหมู่ ออเดอร์ ลูกค้า โปรโมชั่น inventory และ dashboard
- Login/Register พร้อม hash password ด้วย bcrypt
- Access Token อายุ 15 นาที
- Refresh Token อายุ 7-30 วัน ค่าเริ่มต้น 30 วัน
- เก็บ Refresh Token ในฐานข้อมูลแบบ SHA-256 hash
- มี Refresh API พร้อม rotation token
- Logout API สำหรับ revoke refresh token
- มี Role middleware แยกสิทธิ์ admin/user
- ใช้ environment variables สำหรับ secret ไม่มีการ hardcode
- มี middleware บังคับ HTTPS สำหรับ production และอนุญาต localhost สำหรับ development

## Requirements / สิ่งที่ต้องติดตั้ง

- Node.js 20+
- Go 1.24+
- MySQL 8+
- Git

## Backend Setup / วิธีรัน Backend

```bash
cd backend-go
cp .env.example .env
```

Edit `.env` and set real values:

แก้ไฟล์ `.env` แล้วใส่ค่าจริง:

```env
PORT=3001
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=shop_db
TOKEN_SECRET=your_long_random_secret
REFRESH_TOKEN_DAYS=30
ENFORCE_HTTPS=true
ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

Run backend:

รัน backend:

```bash
go mod download
go run ./cmd
```

Backend runs on:

```text
http://localhost:3001
```

## Frontend Setup / วิธีรัน Frontend

```bash
npm install
npm run dev
```

Frontend runs on:

```text
http://localhost:5173
```

Production build:

```bash
npm run build
```

## Authentication Flow / การทำงานของระบบ Login

1. User logs in via `POST /api/login`.
2. Backend returns an access token and refresh token.
3. Frontend stores:
   - `shophub_token`
   - `shophub_refresh_token`
   - `shophub_user`
4. Protected requests send `Authorization: Bearer <access_token>`.
5. If access token expires, frontend calls `POST /api/refresh`.
6. Backend revokes the old refresh token and returns a new token pair.
7. Logout calls `POST /api/logout` to revoke the refresh token.

1. ผู้ใช้ login ผ่าน `POST /api/login`
2. Backend ส่ง access token และ refresh token กลับมา
3. Frontend เก็บข้อมูล:
   - `shophub_token`
   - `shophub_refresh_token`
   - `shophub_user`
4. API ที่ต้อง login ส่ง `Authorization: Bearer <access_token>`
5. ถ้า access token หมดอายุ frontend จะเรียก `POST /api/refresh`
6. Backend revoke refresh token เก่าและออก token คู่ใหม่
7. Logout เรียก `POST /api/logout` เพื่อ revoke refresh token

## Important API Endpoints / API สำคัญ

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| POST | `/api/register` | No | Register user / สมัครสมาชิก |
| POST | `/api/login` | No | Login and receive tokens / เข้าสู่ระบบ |
| POST | `/api/refresh` | No | Rotate refresh token / ต่ออายุ token |
| POST | `/api/logout` | No | Revoke refresh token / ออกจากระบบ |
| GET | `/api/me` | User | Current profile / ข้อมูลผู้ใช้ปัจจุบัน |
| PUT | `/api/me` | User | Update profile / แก้ไขข้อมูลผู้ใช้ |
| GET | `/api/products` | No | Product list / รายการสินค้า |
| GET | `/api/products/:id` | No | Product detail / รายละเอียดสินค้า |
| POST | `/api/products` | Admin | Create product / เพิ่มสินค้า |
| PUT | `/api/products/:id` | Admin | Update product / แก้ไขสินค้า |
| DELETE | `/api/products/:id` | Admin | Delete product / ลบสินค้า |
| GET | `/api/categories` | No | Category list / รายการหมวดหมู่ |
| POST | `/api/categories` | Admin | Create category / เพิ่มหมวดหมู่ |
| PUT | `/api/categories/:id` | Admin | Update category / แก้ไขหมวดหมู่ |
| DELETE | `/api/categories/:id` | Admin | Delete category / ลบหมวดหมู่ |
| GET | `/api/cart` | User | Cart items / ตะกร้าสินค้า |
| POST | `/api/cart` | User | Add cart item / เพิ่มสินค้าในตะกร้า |
| PUT | `/api/cart/:id` | User | Update cart item / แก้ไขตะกร้า |
| DELETE | `/api/cart/:id` | User | Remove cart item / ลบสินค้าออกจากตะกร้า |
| GET | `/api/orders/my` | User | My orders / ออเดอร์ของฉัน |
| GET | `/api/orders` | Admin | All orders / ออเดอร์ทั้งหมด |
| GET | `/api/dashboard/stats` | Admin | Admin stats / สถิติ dashboard |

## Test Commands / คำสั่งทดสอบ

Backend:

```bash
cd backend-go
go test ./...
```

Frontend:

```bash
npm run build
```

## Security Notes / หมายเหตุด้านความปลอดภัย

- Do not commit `.env`.
- Use a long random `TOKEN_SECRET`.
- Keep `ENFORCE_HTTPS=true` in production.
- Set `ALLOWED_ORIGINS` to the real frontend domain in production.
- Refresh tokens are stored only as hashes in the database.
- Logout revokes refresh tokens, so revoked tokens cannot be reused.
- Admin endpoints are protected by role middleware.

- ห้าม commit ไฟล์ `.env`
- ใช้ `TOKEN_SECRET` ที่ยาวและสุ่มจริง
- production ควรใช้ `ENFORCE_HTTPS=true`
- production ต้องตั้ง `ALLOWED_ORIGINS` ให้ตรงกับ domain frontend จริง
- Refresh token ถูกเก็บในฐานข้อมูลเป็น hash เท่านั้น
- Logout จะ revoke refresh token ทำให้ token เดิมใช้ซ้ำไม่ได้
- API ของ admin ถูกป้องกันด้วย role middleware

## Deployment / การ Deploy

For production, build the frontend and deploy the backend with environment variables from your hosting platform. If the backend runs behind a reverse proxy, forward `X-Forwarded-Proto: https` so HTTPS middleware can verify secure traffic.

สำหรับ production ให้ build frontend และตั้ง environment variables ใน hosting/backend server ถ้า backend อยู่หลัง reverse proxy ให้ส่ง header `X-Forwarded-Proto: https` เพื่อให้ HTTPS middleware ตรวจได้ว่า request มาจาก HTTPS จริง
