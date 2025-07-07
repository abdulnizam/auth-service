package handler

import (
	"auth-service/internal/model"
	"auth-service/internal/service"
	"auth-service/internal/utils"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"type"` // optional: defaults to "standard" if empty or invalid
}

type VerifyRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func generateSecureCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(90000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%05d", n.Int64()+10000), nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "hash_error", "Error hashing password")
		return
	}

	token, err := generateSecureCode()
	if err != nil {
		log.Println("Failed to generate code:", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "code_generation_error", "Could not generate verification code")
		return
	}

	user := &model.User{
		Email:             req.Email,
		Password:          hashedPassword,
		VerificationToken: token,
		IsVerified:        false,
		IsActive:          true,
		UserType:          "standard",
	}

	if req.UserType == "admin" || req.UserType == "standard" {
		user.UserType = req.UserType
	}

	if err := model.CreateUser(user); err != nil {
		// Check if it's a duplicate entry error (code 1062)
		if strings.Contains(err.Error(), "Duplicate entry") {
			utils.WriteJSONError(w, http.StatusConflict, "user_exists", "A user with that email already exists.")
			return
		}

		utils.WriteJSONError(w, http.StatusInternalServerError, "user_creation_failed", "User creation failed: "+err.Error())
		return
	}

	if err := utils.SendVerificationEmail(user.Email, token); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "email_send_failed", "Failed to send verification email: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered. Check your email to verify your account."})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	user, err := service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "auth_error", "Authentication failed: "+err.Error())
		return
	}

	if !user.IsVerified {
		utils.WriteJSONError(w, http.StatusForbidden, "account_not_verified", "Please verify your email before logging in.")
		return
	}

	token, err := service.GenerateJWT(user.ID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "token_not_generate", "Token generation failed")
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	user, err := model.GetUserByEmail(req.Email)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "user_not_found", "User not found")
		return
	}

	if user.VerificationToken != req.Token {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_token", "Invalid verification token")
		return
	}

	user.IsVerified = true
	user.VerificationToken = ""
	if err := model.UpdateUser(user); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "verification_failed", "Failed to verify user")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email verified successfully"})
}

func ResendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	user, err := model.GetUserByEmail(req.Email)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "user_not_found", "User not found")
		return
	}

	if user.IsVerified {
		utils.WriteJSONError(w, http.StatusBadRequest, "already_verified", "User is already verified")
		return
	}

	// Generate new token
	token, err := generateSecureCode()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "code_generation_failed", "Could not generate verification code")
		return
	}

	user.VerificationToken = token
	if err := model.UpdateUser(user); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "update_failed", "Could not update user")
		return
	}

	if err := utils.SendVerificationEmail(user.Email, token); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "email_send_failed", "Failed to resend verification email")
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Verification code resent successfully"})
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := model.GetAllUsers()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "db_error", "Failed to fetch users")
		return
	}

	// Optionally hide password and token fields
	type PublicUser struct {
		ID         uint   `json:"id"`
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
		IsActive   bool   `json:"is_active"`
		UserType   string `json:"user_type"`
		CreatedAt  string `json:"created_at"`
	}

	var safeUsers []PublicUser
	for _, u := range users {
		safeUsers = append(safeUsers, PublicUser{
			ID:         u.ID,
			Email:      u.Email,
			IsVerified: u.IsVerified,
			IsActive:   u.IsActive,
			UserType:   u.UserType,
			CreatedAt:  u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(safeUsers)
}
