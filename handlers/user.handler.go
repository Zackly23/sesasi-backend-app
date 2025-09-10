package handlers

import (
	"sesasi-backend-app/models"
	"sesasi-backend-app/schemas"

	// "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	// "gorm.io/gorm"
)


func GetAllUsers(ctx *fiber.Ctx, h *Handler) error  {
	var users []models.User 

	if err := h.DB.Preload("Role.Permissions").Find(&users).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve users",
			"error": err.Error(),
		})
	}

	//mapping ke user list repsonse
	userListResp := schemas.ToUserListResponses(users)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"users": userListResp,
	})
}

func GetUserVerified(ctx *fiber.Ctx, h *Handler) error {
    statusQuery := ctx.Query("verified", "true") // default "true"

    var isVerified bool
    switch statusQuery {
	case "true":
		isVerified = true
	case "false":
		isVerified = false
	default:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status filter, must be 'true' or 'false'",
		})		
	}

	var users []models.User
	if err := h.DB.Preload("Role.Permissions").
		Joins("JOIN roles r ON r.id = users.role_id").
		Where("r.name = ? AND users.is_verified = ?", "user", isVerified).
		Find(&users).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve users",
			})
	}

	//mapping ke user list repsonse
	userListResp := schemas.ToUserListResponses(users)

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Users retrieved successfully",
        "users": userListResp,
    })
}


func ChangeUserRole(ctx *fiber.Ctx, h *Handler) error {
	// user id yang ingin diubah
	userTargetId := ctx.Params("userId")

	//parse
	userTargetIdParsed, err := uuid.Parse(userTargetId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var userTarget models.User 
	if err := h.DB.First(&userTarget, "id = ?", userTargetIdParsed).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var verifikatorRole models.Role
	if err := h.DB.Where("name = ?", models.RoleVerifikator).First(&verifikatorRole).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Role verifikator not found",
		})
	}

	userTarget.RoleID = verifikatorRole.ID
	userTarget.IsVerified = true 

	if err := h.DB.Save(&userTarget).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to change user role",
		})
	}


	if err := h.DB.Save(&userTarget).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to change user role",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users Role successfully Changed to Verifikator",
		
	})
}

func CreateVerifikator(ctx *fiber.Ctx, h *Handler) error {	
	
	var req schemas.CreateVerifikatorRequest

	// Parse body ke struct
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi
	if err := h.Validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

// Cek apakah email sudah terdaftar
	if err := h.DB.Where("email = ?", req.Email).First(&models.User{}).Error; err == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is Already Registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash new password",
			"error": err.Error(),
		})
	}

	defaultAvatar := "" //tambah default avatar

	//ambil role user
	var role models.Role 
	if err := h.DB.Where("name = ?", models.RoleVerifikator).First(&role).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mendapatkan role user",
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
		IsVerified: true,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan user",
		})
	}

	//mapping ke user response
	userResp := schemas.ToUserResponse(user)
	
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Verifikator berhasil dibuat",
			"user": userResp,
		})
}

func VerifyUser(ctx *fiber.Ctx, h *Handler) error {
	// user id yang ingin diubah
	userTargetId := ctx.Params("userId")

	//parse
	userTargetIdParsed, err := uuid.Parse(userTargetId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var userTarget models.User
	if err := h.DB.First(&userTarget, "id = ?", userTargetIdParsed).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Cek apakah user sudah terverifikasi
	if userTarget.IsVerified {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User already verified",
		})
	}

	userTarget.IsVerified = true

	if err := h.DB.Save(&userTarget).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User successfully verified",
	})
}
