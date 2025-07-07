package handler

import (
	"auth-service/internal/model"
	"auth-service/internal/service"
	"auth-service/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AdminCreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AdminCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req AdminCreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request payload")
		return
	}

	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "hash_error", "Could not hash password")
		return
	}

	token, err := generateSecureCode()
	if err != nil {
		log.Println("Failed to generate code:", err)
		utils.WriteJSONError(w, http.StatusInternalServerError, "code_generation_error", "Failed to generate verification code")
		return
	}

	user := &model.User{
		Email:             req.Email,
		Password:          hashedPassword,
		VerificationToken: token,
		IsVerified:        false,
	}

	if err := model.CreateUser(user); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			utils.WriteJSONError(w, http.StatusConflict, "user_exists", "User already exists")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "user_creation_failed", "Could not create user")
		return
	}

	if err := utils.SendVerificationLinkEmail(user.Email, token); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "email_send_failed", "Failed to send verification email")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created and verification email sent."})
}

type UpdateUserRequest struct {
	UserType string `json:"user_type"`
	Active   *bool  `json:"is_active"`
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/admin/users/")
	idUint64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid user ID")
		return
	}
	id := uint(idUint64)
	user, err := model.GetUserByID(id)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "user_not_found", "User not found")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_payload", "Invalid request payload")
		return
	}

	if req.UserType != "" {
		user.UserType = req.UserType
	}
	if req.Active != nil {
		user.IsActive = *req.Active
	}

	if err := model.UpdateUser(user); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "update_failed", "Failed to update user")
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}
