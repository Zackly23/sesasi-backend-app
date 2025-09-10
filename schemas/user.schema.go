package schemas

import (
	"sesasi-backend-app/models"
	"time"

	"github.com/google/uuid"
)

type CreateVerifikatorRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID             uuid.UUID  `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	UserName       string    `json:"user_name,omitempty"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profile_picture,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	IsVerified     bool    `json:"is_verified"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToUserListResponse(user models.User) UserResponse {
	userLoginResponse := UserResponse{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		UserName:       user.UserName,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		Phone:          user.Phone,
		IsVerified:     user.IsVerified,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	return userLoginResponse
}

func ToUserListResponses(users []models.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, ToUserListResponse(user))
	}
	return responses
}

func ToUserResponse(user models.User) UserResponse {
	userSignUpResponse := UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		IsVerified: user.IsVerified,
	}

	return userSignUpResponse
}