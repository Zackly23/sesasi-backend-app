package seeders

import (
	"fmt"
	"math/rand"
	"sesasi-backend-app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func PengajuanIzinSeed(db *gorm.DB) error {
	// Ambil semua user dengan role User
	var users []models.User
	if err := db.Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", models.RoleUser).
		Find(&users).Error; err != nil {
		return err
	}

	jenisIzinOptions := []string{"cuti", "libur", "sakit", "lainnya"}
	alasanOptions := []string{
		"Menghadiri pernikahan kerabat dekat",
		"Menghadiri acara keluarga besar di luar kota",
		"Sedang kurang sehat dan butuh istirahat",
		"Mengurus administrasi penting di kantor pemerintah",
		"Perlu menghadiri wisuda saudara kandung",
	}
	keteranganOptions := []string{
		"Mohon izin diberikan karena ini acara yang tidak bisa ditinggalkan.",
		"Rencana izin hanya sementara, aktivitas akan segera dilanjutkan setelah selesai.",
		"Dokumen pendukung bisa saya lampirkan jika diperlukan.",
		"Harap pengajuan ini dapat segera diproses, terima kasih.",
		"Izin ini sangat penting untuk urusan keluarga.",
	}

	statusOptions := []string{
		"diajukan",
		"dibatalkan",
		"direvisi",
		"diterima",
		"ditolak",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for index, user := range users {
		// tanggal mulai random dari hari ini s/d +7
		tanggalMulai := time.Now().AddDate(0, 0, r.Intn(8))
		// tanggal selesai max +4 hari dari tanggal mulai
		tanggalSelesai := tanggalMulai.AddDate(0, 0, r.Intn(5)+1)

		var pengajuanId uuid.UUID
		var status string
		if index == 1 {
			pengajuanId = uuid.MustParse("a57e7c4b-b3a1-4374-b82a-7996045bedd1")
			status = string(models.IzinDiajukan)
		} else {
			pengajuanId = uuid.New()
			status = statusOptions[r.Intn(len(statusOptions))]
		}

		izin := models.PengajuanIzin{
			ID:             pengajuanId,
			UserID:         user.ID,
			TanggalMulai:   tanggalMulai,
			TanggalSelesai: tanggalSelesai,
			JenisIzin:      models.TypePengajuanIzin(jenisIzinOptions[r.Intn(len(jenisIzinOptions))]),
			AlasanIzin:     alasanOptions[r.Intn(len(alasanOptions))],
			KeteranganIzin: keteranganOptions[r.Intn(len(keteranganOptions))],
			Status: models.PengajuanIzinStatus(status),
		}

		if err := db.Create(&izin).Error; err != nil {
			fmt.Printf("Error creating Pengajuan Izin %s: %v\n", izin.ID, err)
		} else {
			fmt.Printf("Pengajuan Izin %s created successfully\n", izin.ID)
		}
	}

	return nil
}

func PengajuanIzinCommentSeed(db *gorm.DB) error {
	// Ambil user yang rolenya verifikator
	var verifikators []models.User
	if err := db.Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", models.RoleVerifikator).
		Find(&verifikators).Error; err != nil {
		return err
	}

	// Ambil semua pengajuan izin yang statusnya diterima / ditolak
	var pengajuanIzins []models.PengajuanIzin
	if err := db.Where("status IN ?", []string{string(models.IzinDiterima), string(models.IzinDitolak)}).
		Find(&pengajuanIzins).Error; err != nil {
		return err
	}

	if len(verifikators) == 0 || len(pengajuanIzins) == 0 {
		return nil // tidak ada data untuk di-seed
	}

	// Komentar dummy
	comments := []string{
		"Pengajuan sudah lengkap dan valid, saya setujui.",
		"Dokumen pendukung kurang jelas, mohon diperbaiki di pengajuan berikutnya.",
		"Izin tidak dapat diterima karena tidak memenuhi syarat administrasi.",
		"Alasan cukup kuat, pengajuan saya terima.",
		"Perlu koordinasi lebih lanjut, namun sementara disetujui.",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, izin := range pengajuanIzins {
		randomUser := verifikators[r.Intn(len(verifikators))]
		comment := models.PengajuanIzinComment{
			ID:              uuid.New(),
			PengajuanIzinID: izin.ID,
			UserID:          randomUser.ID,
			Status:          izin.Status,
			Comment:         comments[r.Intn(len(comments))],
			CreatedAt:       time.Now(),
		}

		if err := db.Create(&comment).Error; err != nil {
			fmt.Printf("Error creating comment %s: %v\n", comment.ID, err)
		} else {
			fmt.Printf("Comment %s created successfully\n", comment.ID)
		}
	}

	return nil
}