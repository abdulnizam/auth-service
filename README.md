# 🛡️ Auth Service (Go + MariaDB + JWT + Mailjet)

A lightweight authentication microservice built with Go, supporting:

-   ✅ User registration and login
-   📧 Email verification with a secure code (via Mailjet)
-   📧 Email verification with a secure link (via Mailjet)
-   🔐 JWT-based authentication
-   🐬 MariaDB for persistent user storage
-   🛠️ Admin dashboard (Next.js)
-   🐳 Docker support for easy deployment

---

## 📦 Features

-   Register users with hashed passwords
-   Send email verification codes
-   Send email verification **links**
-   Prevent login until email is verified
-   JWT token generation on login
-   Admin dashboard to manage users
    -   Activate/deactivate accounts
    -   Change user type (admin/standard)
-   RESTful API endpoints
-   Dockerized setup for local/cloud deployment

---

## 🏗️ Project Structure

```
auth-service/
├── cmd/                # Entry point (main.go)
├── config/             # Env config loader
├── internal/
│   ├── db/             # DB init (MariaDB)
│   ├── model/          # User model
│   ├── handler/        # HTTP handlers
│   ├── service/        # Business logic
│   └── utils/          # Mailjet, hash, jwt helpers
├── Dockerfile
├── docker-compose.yml
├── go.mod / go.sum
└── README.md
```

---

## 🚀 Getting Started Locally

### 1. Clone the repo

```bash
git clone https://github.com/abdulnizam/auth-service-go.git
cd auth-service
```

### 2. Set up your `.env`

```env
PORT=8080
DB_USER=xxx
DB_PASS=xxx
DB_HOST=localhost
DB_PORT=3306
DB_NAME=auth_service
JWT_SECRET=your-super-secret
MJ_API_KEY=your-mailjet-api-key
MJ_SECRET_KEY=your-mailjet-secret
EMAIL_FROM=verify@xxx.com
```

> 💡 Ensure MariaDB is running with the `auth_service` DB created and your user has privileges.

---

### 3. Run locally

```bash
go run ./cmd/main.go
```

Server should start at: `http://localhost:8080`

---

## 🐳 Docker Support

### 1. Build & Run (uses local MariaDB)

```bash
docker compose up --build auth-service -d
```

> Use `DB_HOST=host.docker.internal` if MariaDB runs on your host machine.

---

## 🧪 API Endpoints

| Method | Endpoint           | Description                       |
| ------ | ------------------ | --------------------------------- |
| POST   | `/auth/register`   | Register new user                 |
| POST   | `/auth/verify`     | Verify email with secure link     |
| POST   | `/auth/login`      | Login with email & password       |
| POST   | `/auth/resend`     | Resend verification email         |
| GET    | `/users`           | Get all users (for dashboard)     |
| POST   | `/admin/users`     | Admin creates user & sends link   |
| PUT    | `/admin/users/:id` | Update user type or active status |

---

## 🧪 Sample Payloads

### Register

```json
POST /auth/register
{
  "email": "test@example.com",
  "password": "password123"
}
```

### Verify (from email link)

```http
GET /verify?token=abc123&email=test@example.com
```

### Admin Create

```json
POST /admin/users
{
  "email": "user@example.com",
  "password": "secure",
  "type": "standard"  // optional
}
```

### Admin Update

```json
PUT /admin/users/1
{
  "is_active": true,
  "user_type": "admin"
}
```

---

## 🛠️ Built With

-   [Go](https://golang.org/)
-   [GORM](https://gorm.io/)
-   [MariaDB](https://mariadb.org/)
-   [Mailjet](https://www.mailjet.com/)
-   [Docker](https://www.docker.com/)
-   [JWT](https://jwt.io/)
-   [Next.js (Dashboard)](https://nextjs.org/)

---

## 👨‍💻 Author

**Abdul Nizam**  
[abdulnizam.com](https://abdulnizam.com)

---

## 📄 License

MIT License
