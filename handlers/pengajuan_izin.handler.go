package handlers

import (
	"errors"
	"sesasi-backend-app/models"
	"sesasi-backend-app/schemas"
	"sesasi-backend-app/utils"

	// "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


func GetAllPengajuanIzin(ctx *fiber.Ctx, h *Handler) error {
	var pengajuanIzins []models.PengajuanIzin

	if err := h.DB.Preload("User.Role").Preload("Comments.User.Role").Find(&pengajuanIzins).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve pengajuan izin",
			"error":   err.Error(),
		})
	}

	// Jika record kosong
	if len(pengajuanIzins) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No pengajuan izin found",
			"izin":    []models.PengajuanIzin{},
		})
	}

	// mAPPING ke Pengujian Izin Response 
	pengajuanIzinResp := schemas.ToPengajuanIzinResponses(pengajuanIzins)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pengajuan Izin retrieved successfully",
		"pengajuan_izins":    pengajuanIzinResp,
	}) 	
}

func GetFilteredPengajuanIzin(ctx *fiber.Ctx, h *Handler) error {
    var pengajuanIzins []models.PengajuanIzin

	// filter status pengajuan izin
    status := ctx.Query("status", "all")

	//cek apakah status query ada di list enum status pengajuan
	isValid := utils.IsValidPengajuanStatus(status)

	if !isValid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status filter",
		})
	}

    query := h.DB.Preload("User.Role").Preload("Comments.User.Role")

	// kalau query status bukan "all" (semua) maka filter bedasarkan statusya
    if status != "all" {
        query = query.Where("status = ?", models.PengajuanIzinStatus(status))
    }

	// ambil record
    if err := query.Find(&pengajuanIzins).Error; err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to retrieve Pengajuan izin",
			"error": err.Error(),
        })
    }

	//mapping ke Pengajuan Izin Responses
	pengajuanIzinResp := schemas.ToPengajuanIzinResponses(pengajuanIzins)

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Pengajuan Izin retrieved successfully",
        "pengajuan_izins":    pengajuanIzinResp,
    })
}


func UpdateStatusPengajuanIzin(ctx *fiber.Ctx, h *Handler) error {
	// ambil user_id dari locals ram
	userId := ctx.Locals("user_id").(string)
	userIdParse, errParse := uuid.Parse(userId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse User Id to UUID",
		})
	}

	//ambil izinId dari params
	izinId := ctx.Params("izinId")
	izinIdParse, errParse := uuid.Parse(izinId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse Izin Id to UUID",
			"error":   errParse.Error(),
		})
	}

	var pengajuanIzinReq schemas.VerifikasiPengajuanIzinRequest

	if err := ctx.BodyParser(&pengajuanIzinReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	//validasi request body
	if err := h.Validator.Struct(pengajuanIzinReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error": err.Error(),
		})
	}

	var pengajuanIzin models.PengajuanIzin
	//ambil record pengajuan izin berdasarkan izinId
	if err := h.DB.Find(&pengajuanIzin, "id = ?", izinIdParse).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve pengajuan izin",
			"error":   err.Error(),
		})
	}

	//Buat data comment baru 
	comment := models.PengajuanIzinComment{
		PengajuanIzinID: pengajuanIzin.ID,
		UserID: userIdParse,
		Comment: pengajuanIzinReq.Comment,
		Status: models.PengajuanIzinStatus(pengajuanIzinReq.Status),
	}

	//simpan record comment
	if err := h.DB.Create(&comment).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create comment",
		})
	}

	// ubah status pengajuan izin sesuai status request
	pengajuanIzin.Status = models.PengajuanIzinStatus(pengajuanIzinReq.Status)

	//simpen perubahan pengajunan izIn
	if err := h.DB.Save(&pengajuanIzin).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update status pengajuan izin",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pengajuan Izin Successfully Updated",		
	}) 	
}

func CreatePengajuanIzin(ctx *fiber.Ctx, h *Handler) error {
	// ambil user_id dari locals ram
	userId := ctx.Locals("user_id").(string)
	userIdParse, errParse := uuid.Parse(userId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse User Id to UUID",
			"error":   errParse.Error(),
		})
	}

	var pengajuanIzinReq schemas.PengajuanIzinRequest

	if err := ctx.BodyParser(&pengajuanIzinReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	//validasi request body
	if err := h.Validator.Struct(pengajuanIzinReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error": err.Error(),
		})
	}

	pengajuanIzin := models.PengajuanIzin{
		UserID: userIdParse,
		AlasanIzin: pengajuanIzinReq.AlasanIzin,
		JenisIzin: models.TypePengajuanIzin(pengajuanIzinReq.JenisIzin),
		KeteranganIzin: pengajuanIzinReq.KeteranganIzin,
		Status: models.IzinDiajukan,
		TanggalMulai: pengajuanIzinReq.TanggalMulai,
		TanggalSelesai: pengajuanIzinReq.TanggalSelesai,
	}

	//s impan record pengajuan izin
	if err := h.DB.Create(&pengajuanIzin).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to create pengajuan izin",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pengajuan Izin Successfully Created",	
		"pengajuan_izin_id" : pengajuanIzin.ID,	
	}) 

}

func GetUserPengajuanIzin(ctx *fiber.Ctx, h *Handler) error {

	userId := ctx.Locals("user_id").(string)
	userIdParse, errParse := uuid.Parse(userId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse User Id to UUID",
			"error":   errParse.Error(),
		})
	}

	//ambilr record pengajuan izin berdasarkan user yang login
	var pengajuanIzins []models.PengajuanIzin
	if err := h.DB.Preload("User").Where("user_id = ?", userIdParse).Find(&pengajuanIzins).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to retrieve pengajuan izin",
			"error":   err.Error(),
		})	
	}

	//maping ke Pengajuan Izin Responses 
	pengajuanIzinResp := schemas.ToPengajuanIzinResponses(pengajuanIzins)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pengajuan Izin retrieved successfully",	
		"pengajuan_izins" : pengajuanIzinResp,	
	}) 
}

func GetDetailPengajuanIzin(ctx *fiber.Ctx, h *Handler) error {
	//ambil izinid dari params
	pengajuanIzinId := ctx.Params("izinId")
	pengajuanIzinIdParse, errParse := uuid.Parse(pengajuanIzinId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse User Id to UUID",
			"error":   errParse.Error(),
		})
	}

	//ambil user_id dari locals ram
	userId := ctx.Locals("user_id").(string)
	userIdParse, errParse := uuid.Parse(userId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
			"error":   errParse.Error(),
		})
	}

	//ambil record pengajuan izin berdasarkan izinId dan user yang login
	var pengajuanIzin models.PengajuanIzin
	if err := h.DB.Preload("User").Preload("Comments.User.Role").Where("id = ? AND user_id = ?", pengajuanIzinIdParse,  userIdParse).First(&pengajuanIzin).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Pengajuan Izin Failed Retrieved",
			"error":   err.Error(),
		})	
	}

	//mapping ke Pengajuan Izin Detail Response
	pengajuanIzinResp := schemas.ToPengajuanIzinDetailResponse(pengajuanIzin)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pengajuan Izin retrieved successfully",
		"pengajuan_izin" : pengajuanIzinResp,		
	})
}

func UpdateDetailPengajuanIzin(ctx *fiber.Ctx, h *Handler) error {
	//ambil izin id dari params
	pengajuanIzinId := ctx.Params("izinId")
	pengajuanIzinIdParse, errParse := uuid.Parse(pengajuanIzinId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse Izin Id to UUID",
			"error":   errParse.Error(),
		})
	}

	// ambil user_id dari locals RAM
	userId := ctx.Locals("user_id").(string)
	userIdParse, errParse := uuid.Parse(userId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse User Id to UUID",
			"error":   errParse.Error(),
		})
	}

	var pengajuanIzinReq schemas.PengajuanIzinRequest

	if err := ctx.BodyParser(&pengajuanIzinReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	//valdaasi request body
	if err := h.Validator.Struct(pengajuanIzinReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error": err.Error(),
		})
	}

	//ambil record pengajuan izin berdasarkan izinId dan user yang login
	var pengajuanIzin models.PengajuanIzin
	if err := h.DB.Preload("User").Where("id = ? AND user_id = ?", pengajuanIzinIdParse,  userIdParse).First(&pengajuanIzin).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Pengajuan Izin Failed Retrieved",
			"error":   err.Error(),
		})	
	}	
	
	//update field pengajuan izin
	pengajuanIzin.TanggalMulai = pengajuanIzinReq.TanggalMulai
	pengajuanIzin.TanggalSelesai = pengajuanIzinReq.TanggalSelesai
	pengajuanIzin.KeteranganIzin = pengajuanIzinReq.KeteranganIzin
	pengajuanIzin.AlasanIzin = pengajuanIzinReq.AlasanIzin
	pengajuanIzin.JenisIzin = models.TypePengajuanIzin(pengajuanIzinReq.JenisIzin)

	if err := h.DB.Save(&pengajuanIzin).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to update pengajuan izin",
			"error":   err.Error(),
		})	
	}	

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pengajuan Izin Successfully Updated",
	})

}

func DeletePengajuanIzin(ctx *fiber.Ctx, h *Handler) error {
    // Ambil pengajuan izin ID dari parameter
    pengajuanIzinId := ctx.Params("izinId")
    pengajuanIzinIdParse, errParse := uuid.Parse(pengajuanIzinId)
    if errParse != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Failed to parse Izin Id to UUID",
			"error": errParse.Error(),
        })
    }

    // Ambil user ID dari ram (login)
    userId := ctx.Locals("user_id").(string)
    userIdParse, errParse := uuid.Parse(userId)
    if errParse != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Failed to parse User Id to UUID",
			"error": errParse.Error(),
        })
    }

    // Hapus pengajuan izin dengan filter user (biar user hanya bisa menghapus izinnya sendiri)
    res := h.DB.Where("id = ? AND user_id = ?", pengajuanIzinIdParse, userIdParse).Delete(&models.PengajuanIzin{})
    if res.Error != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to delete pengajuan izin",
			"error": res.Error.Error(),
        })
    }

	//kalau tidak ada record yang dihapus
    if res.RowsAffected == 0 {
        return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "No pengajuan izin found to delete",
			"error": res.Error.Error(),
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Pengajuan izin successfully deleted",
    })
}

func CancelPengajuanIzin(ctx *fiber.Ctx, h *Handler) error {
    // Ambil pengajuan izin ID dari parameter
    pengajuanIzinId := ctx.Params("izinId")
    pengajuanIzinIdParse, errParse := uuid.Parse(pengajuanIzinId)
    if errParse != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Failed to parse Izin Id to UUID",
			"error": errParse.Error(),
        })
    }

    // Ambil user ID dari context
    userId := ctx.Locals("user_id").(string)
    userIdParse, errParse := uuid.Parse(userId)
    if errParse != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Failed to parse User Id to UUID",
			"error": errParse.Error(),
        })
    }

    // Ambil pengajuan izin milik user
    var pengajuan models.PengajuanIzin
    res := h.DB.Where("id = ? AND user_id = ?", pengajuanIzinIdParse, userIdParse).First(&pengajuan)
    if res.Error != nil {
        if errors.Is(res.Error, gorm.ErrRecordNotFound) {
            return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "message": "No pengajuan izin found",
				"error": res.Error.Error(),
            })
        }

        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to retrieve pengajuan izin",
			"error": res.Error.Error(),
        })
    }

    // Update status menjadi dibatalkan
    pengajuan.Status = models.IzinDiBatalkan
    if err := h.DB.Save(&pengajuan).Error; err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to cancel pengajuan izin",
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Pengajuan izin successfully canceled",
    })
}

func GetStatusPengajuanIzin(ctx *fiber.Ctx, h *Handler) error {
	//ambil izinid dari params
	pengajuanIzinId := ctx.Params("izinId")
	pengajuanIzinIdParse, errParse := uuid.Parse(pengajuanIzinId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse User Id to UUID",
			"error":   errParse.Error(),
		})
	}

	//ambil user_id dari locals ram
	userId := ctx.Locals("user_id").(string)
	userIdParse, errParse := uuid.Parse(userId)

	if errParse != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
			"error":   errParse.Error(),
		})
	}

	//ambil record pengajuan izin berdasarkan izinId dan user yang login
	var pengajuanIzin models.PengajuanIzin
	if err := h.DB.Preload("User").Preload("Comments.User.Role").Where("id = ? AND user_id = ?", pengajuanIzinIdParse,  userIdParse).First(&pengajuanIzin).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Pengajuan Izin Failed Retrieved",
			"error":   err.Error(),
		})	
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Status Pengajuan Izin retrieved successfully",
		"pengajuan_izin_status" : pengajuanIzin.Status,
    })
}