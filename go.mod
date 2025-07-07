module auth-service

go 1.21

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.21.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.7
)

require github.com/rs/cors v1.11.1 // indirect

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/google/uuid v1.6.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mailjet/mailjet-apiv3-go/v4 v4.0.7 // indirect
)
