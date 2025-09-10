package utils

import (
	"crypto/rand"
	"encoding/hex"
	"sesasi-backend-app/models"
	"strings"
	"time"

	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthTokenJWT(ctx *fiber.Ctx) (*jwt.Token, error) {
	// Ambil token dari header Authorization
	accessToken := ctx.Get("Authorization")
	if accessToken == "" {
		return nil, fmt.Errorf("authorization token tidak ditemukan")
	}

	// Jika pakai "Bearer <token>", hapus prefix
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")


	// Parse dan verifikasi token
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak valid: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token tidak valid atau expired")
	}

	return token, nil
}

func GetUserID(ctx *fiber.Ctx) (uuid.UUID, error) {
	// Validasi dan ambil token
	token, err := AuthTokenJWT(ctx)
	if err != nil {
		return uuid.UUID{},  ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Error Retrieved Token From Bearer",
			"error": err.Error(),
		})
	}

	// Ambil user_id dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.UUID{}, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Failed Parse claim token",
		})
	}

	// Ambil user_id dan parse ke uuid.UUID
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.UUID{}, ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No found User ID in token",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "UUID format is not valid",
			"error": err.Error(),
		})
	}
	// userID := uint(userIDFloat)

	return userID, nil
}


func GenerateRandomToken() (string, error) {
    bytes := make([]byte, 32) // 256 bit = 32 byte
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func GenerateToken(user models.User, duration time.Duration) (string, error) {

	//buat claim data
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name": user.FirstName + " " + user.LastName,
		"role": user.Role.Name,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	}

	//kunci jwt
	jwtSecret := os.Getenv("JWT_SECRET_KEY")

	//metode claim dan signed jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))

	return signedToken, err
}