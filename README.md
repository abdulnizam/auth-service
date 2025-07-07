# 🛡️ Auth Service (Go + MariaDB + JWT + Mailjet)

A lightweight authentication microservice built with Go, supporting:

-   ✅ User registration and login
-   📧 Email verification with a secure code (via Mailjet)
-   🔐 JWT-based authentication
-   🐬 MariaDB for persistent user storage
-   🐳 Docker support for easy deployment

---

## 📦 Features

-   Register users with hashed passwords
-   Send email verification codes
-   Prevent login until email is verified
-   Secure JWT token generation on login
-   RESTful API endpoints
-   Dockerized setup for consistent local and cloud deployment

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
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-secret
GOOGLE_REDIRECT_URI=http://localhost:8080/auth/google/callback
SMTP_HOST=in.mailjet.com
SMTP_PORT=587
SMTP_USER=your-mailjet-api-key
SMTP_PASS=your-mailjet-secret-key
EMAIL_FROM=verify@xxx.com
```

> 💡 Ensure your local MariaDB is running with the `auth_service` DB created and user `auth` has correct access.

### 3. Run locally

```bash
go run ./cmd/main.go
```

Server should start on: `http://localhost:8080`

---

## 🐳 Docker Support

### 1. Build & Run (Connects to host's MySQL)

```bash
docker compose up --build auth-service -d
```

> ⚠️ Ensure your `.env` uses `DB_HOST=host.docker.internal` to connect to your Mac's local MySQL from inside Docker.

---

## 🧪 API Endpoints

| Method | Endpoint         | Description                 |
| ------ | ---------------- | --------------------------- |
| POST   | `/auth/register` | Register new user           |
| POST   | `/auth/verify`   | Verify email with code      |
| POST   | `/auth/login`    | Login with email & password |

---

## 🔐 Sample JSON Payloads

### Register

```json
POST /auth/register
{
  "email": "test@example.com",
  "password": "password123"
}
```

### Verify Email

```json
POST /auth/verify
{
  "email": "test@example.com",
  "token": "12345"
}
```

### Login

```json
POST /auth/login
{
  "email": "test@example.com",
  "password": "password123"
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

---

## 👨‍💻 Author

**Abdul Nizam**  
[abdulnizam.com](https://abdulnizam.com)

---

## 📄 License

This project is licensed under the MIT License.
