package schemas

import (
	"sesasi-backend-app/models"
	"time"

	"github.com/google/uuid"
)

type UserSignUpRequest struct {
    FirstName       string `json:"first_name" validate:"required,max=50"`
    LastName        string `json:"last_name" validate:"required,max=50"`
    Email           string `json:"email" validate:"required,email,max=100"`
    Password        string `json:"password" validate:"required,min=6,max=50"`
    PasswordConfirm string `json:"password_confirmation" validate:"required"`
}


type UserLoginRequest struct {
    Email      string `json:"email" validate:"required,email,max=100"`
    Password   string `json:"password" validate:"required,min=6,max=50"`
    RememberMe bool   `json:"remember_me"`
}


type ResetPasswordRequest struct {
    Email      string `json:"email" validate:"required,email,max=100"`
	NewPassword		string	`json:"new_password" validate:"required"`
}

type UpdatePasswordRequest struct {
    Email      string `json:"email" validate:"required,email,max=100"`
	NewPassword					string	`json:"new_password" validate:"required"`
	NewPasswordConfirmation     string 	`json:"new_password_confirmation" validate:"required"`
}

// Response yang aman
type UserLoginResponse struct {
	ID           uuid.UUID   `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UserName   string `json:"user_name,omitempty"`
	Email        string `json:"email"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	Status 	 bool `json:"status,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Role string `json:"role,omitempty"`
	Permission []string `json:"permission,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserSignUpResponse struct {
	ID           uuid.UUID   `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`	
}

func ToLoginResponse(user models.User, permissions []string) UserLoginResponse {
	userLoginResponse := UserLoginResponse{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		UserName:       user.UserName,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		Phone:          user.Phone,
		Role:           user.Role.Name,
		Status:         user.IsVerified,
		Permission:     permissions,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	return userLoginResponse
}

func ToSignupResponse(user models.User) UserSignUpResponse {
	userSignUpResponse := UserSignUpResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return userSignUpResponse
}