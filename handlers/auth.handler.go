package handlers

import (
	"errors"
	"fmt"
	"os"
	"sesasi-backend-app/models"
	"sesasi-backend-app/schemas"
	"sesasi-backend-app/utils"

	"time"

	// "github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(ctx *fiber.Ctx, h *Handler) error {
	//get request body
	var req schemas.UserSignUpRequest

	// Parse body ke struct
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":  err.Error(),
		})
	}

	// Validasi
	if err := h.Validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error": err.Error(),
		})
	}

	var existingUser models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		// Data ditemukan
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email Already Registered",
		})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Error lain (misalnya koneksi DB)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error checking existing email",
			"error":   err.Error(),
		})
	}

	// Cek apakah password sesuai dengan konfirmasi password
	if req.Password != req.PasswordConfirm {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password confirmation does not match",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
			"error":  err.Error(),
		})
	}

	defaultAvatar := "" //tambah default avatar

	//ambil role user
	var role models.Role 
	if err := h.DB.Where("name = ?", models.RoleUser).First(&role).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to Get Role User",
			"error":  err.Error(),
		})
	}

	// Simpan user baru ke database
	user := models.User{
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Password:    string(hashedPassword),
		Role: role,
		ProfilePicture: defaultAvatar,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to Create User",
			"error":  err.Error(),
		})
	}

	//mapping ke Signup Response
	userSignUpResponse := schemas.ToSignupResponse(user)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User Successfully Created",
		"user": userSignUpResponse,
	})
	
}

func Login(ctx *fiber.Ctx, h *Handler) error {
	var req schemas.UserLoginRequest
	var user models.User

	// Parse body
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":  err.Error(),
		})
	}

	// Validasi request
	if err := h.Validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error": err.Error(),
		})
	}

	// Cari user by email, preload Role dan Permissions
	if err := h.DB.Preload("Role.Permissions").Where("email = ?", req.Email).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Email Not Found in Record",
			"error":  err.Error(),
		})
	}

	// Cek password aakah sama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Password is Not Matched",
			"error":  err.Error(),
		})
	}

	// Ambil secret key JWT
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "JWT_SECRET_KEY is not set in environment variables or Missing",
		})
	}

	// Generate token
	accessToken, accessTokenErr := utils.GenerateToken(user, time.Hour*24)
	refreshToken, refreshTokenErr := utils.GenerateToken(user, time.Hour*24*14)
	if accessTokenErr != nil || refreshTokenErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
			"error":   accessTokenErr.Error(),
		})
	}

	// Convert permissions ke []string
	var permissions []string
	
	for _, p := range user.Role.Permissions {
		permissions = append(permissions, p.Name)
	
	}

	//buat PAT
	privateAccessToken := models.PrivateAccessToken{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserID:          user.ID,
		IPAddress:       ctx.IP(),
		AccessTokenExp:  time.Now().Add(time.Hour * 2),
		RefreshTokenExp: time.Now().Add(time.Hour * 24 * 7),
		Revoked:         false,
	}

	if err := h.DB.Create(&privateAccessToken).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to Stored Private Access Token",
		})
	}

	//mapping ke Login Response
	userLoginResponse := schemas.ToLoginResponse(user, permissions)

	// Response JSON
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    userLoginResponse,
		"access_token":  accessToken, 
		"refresh_token": refreshToken,
	})
}


func Logout(ctx *fiber.Ctx, h *Handler) error {
	// Ambil user ID DARI BODY
	userId := ctx.Locals("user_id").(string)

	userIdParse, errParse := uuid.Parse(userId)

	if errParse != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse User Id",
			"error":  errParse.Error(),
		})
	}

	// Hapus token dari database
	if err := h.DB.Model(&models.PrivateAccessToken{}).Where("user_id = ?", userIdParse).Delete(&models.PrivateAccessToken{}).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete token",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully Logout",
	})

}


//return refresh token
func Refresh(ctx *fiber.Ctx, h *Handler) error {

	refreshToken := ctx.Get("Authorization")

	if refreshToken == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Refresh token is not Found",
		})
	}

	// Verifikasi dan parse token
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	token, errParse := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if errParse != nil || !token.Valid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Token is not valid",
			"error":  errParse.Error(),
		})
	}

	userIDStr := token.Claims.(jwt.MapClaims)["user_id"].(string)

	userIdParse, errParse := uuid.Parse(userIDStr)
	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User ID in token is not valid",
			"error":  errParse.Error(),
		})
}


	var user models.User
	if err := h.DB.First(&user, userIdParse).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
			"error":  err.Error(),
		})
	}

	newAccessToken, err := utils.GenerateToken(user, time.Hour*2)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create new access token",
			"error":  err.Error(),
		})
	}
	

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully generated new access token",
		"access_token": newAccessToken,
	})
}


func ResetPassword(ctx *fiber.Ctx, h *Handler) error {
	var resetPasswordReq schemas.ResetPasswordRequest
	var userTarget models.User

	if err := ctx.BodyParser(&resetPasswordReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read request body",
			"error": err.Error(),
		})
	}

	if err := h.Validator.Struct(resetPasswordReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request validation failed",
			"error": err.Error(),
		})
	}

	//cek ke database ada ga emailnya
	if err := h.DB.Where("email = ?", resetPasswordReq.Email).First(&userTarget).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Email not found",
			"error":  err.Error(),
		})

	}

	// Hash the new password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(resetPasswordReq.NewPassword), bcrypt.DefaultCost)
	if errHash != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash new password",
			"error":  errHash.Error(),
		})
	}

	// Update the user's password in the database
	if err := h.DB.Model(&userTarget).Update("password", string(hashedPassword)).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update password",
			"error":  err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset request successful",
	})
}

func UpdatePassword(ctx *fiber.Ctx, h *Handler) error {
	var updatePasswordReq schemas.UpdatePasswordRequest

	if err := ctx.BodyParser(&updatePasswordReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read request body",
			"error":   err.Error(),
		})
	}

	if err := h.Validator.Struct(updatePasswordReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request validation failed",
			"error":   err.Error(),
		})
	}

	// Get user ID from local context (RAM)
	userId := ctx.Locals("user_id").(string)
	userIdParse, errParse := uuid.Parse(userId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse user ID",
			"error":   errParse.Error(),
		})
	}

	// Find user by ID
	var user models.User
	if err := h.DB.Where("id = ?", userIdParse).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Account not found",
			"error":   err.Error(),
		})
	}

	// Check if password and confirmation match
	if updatePasswordReq.NewPassword != updatePasswordReq.NewPasswordConfirmation {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password confirmation does not match",
		})
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatePasswordReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash new password",
			"error":   err.Error(),
		})
	}

	// Update user's password
	if err := h.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update password",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

