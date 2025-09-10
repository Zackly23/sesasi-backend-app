package schemas

import (
	"sesasi-backend-app/models"
	"time"

	"github.com/google/uuid"
)

// Data Transfer Object (DTO)

type VerifikasiPengajuanIzinRequest struct {
	Status  string `json:"status" validate:"required,oneof=diterima ditolak revisi dibatalkan diajukan"`
	Comment string `json:"comment" validate:"max=500"` // comment optional, max 500 karakter
}

type PengajuanIzinRequest struct {
	TanggalMulai   time.Time `json:"tanggal_mulai" validate:"required"`
	TanggalSelesai time.Time `json:"tanggal_selesai" validate:"required,gtfield=TanggalMulai"`
	JenisIzin      string    `json:"jenis_izin" validate:"required,oneof=cuti libur sakit lainnya"`
	AlasanIzin     string    `json:"alasan_izin" validate:"required,min=5,max=500"`
	KeteranganIzin string    `json:"keterangan_izin" validate:"max=1000"` // opsional
}

type UserPengajuanIzin struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

type KomentatorPengajuanIzin struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
}

type UserCommentPengajuanIzin struct {
	ID              uuid.UUID               `json:"id"`
	Comment         string                  `json:"comment"`
	StatusPengajuan string                  `json:"status_pengajuan"`
	Komentator      KomentatorPengajuanIzin `json:"commentator"`
}

type PengajuanIzinResponse struct {
	ID             uuid.UUID         `json:"id"`
	TanggalMulai   string            `json:"tanggal_mulai"`
	TanggalSelesai string            `json:"tanggal_selesai"`
	JenisIzin      string            `json:"jenis_izin"`
	Status         string            `json:"status"`
	PengajuanUser  UserPengajuanIzin `json:"pengajuan_user"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      *time.Time        `json:"deleted_at,omitempty"`
}

type PengajuanIzinDetailResponse struct {
	ID             uuid.UUID                  `json:"id"`
	TanggalMulai   string                     `json:"tanggal_mulai"`
	TanggalSelesai string                     `json:"tanggal_selesai"`
	JenisIzin      string                     `json:"jenis_izin"`
	AlasanIzin     string                     `json:"alasan_izin"`
	KeteranganIzin string                     `json:"keterangan_izin"`
	Status         string                     `json:"status"`
	PengajuanUser  UserPengajuanIzin          `json:"pengajuan_user"`
	Comments       []UserCommentPengajuanIzin `json:"comments"`
	CreatedAt      time.Time       	           `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
	DeletedAt      *time.Time                 `json:"deleted_at,omitempty"`
}


// Mapper
func ToPengajuanIzinResponse(izin models.PengajuanIzin) PengajuanIzinResponse {
	return PengajuanIzinResponse{
		ID:             izin.ID,
		TanggalMulai:   izin.TanggalMulai.Format("02 January 2006"),
		TanggalSelesai: izin.TanggalSelesai.Format("02 January 2006"),
		JenisIzin:      string(izin.JenisIzin),
		Status:         string(izin.Status),
		PengajuanUser: UserPengajuanIzin{
			ID:        izin.User.ID,
			FirstName: izin.User.FirstName,
			LastName:  izin.User.LastName,
			ProfilePicture: izin.User.ProfilePicture,
		},
		CreatedAt:  izin.CreatedAt,
		UpdatedAt:  izin.UpdatedAt,
		DeletedAt:  &izin.DeletedAt.Time,
	}
}

func ToPengajuanIzinResponses(izins []models.PengajuanIzin) []PengajuanIzinResponse {
	responses := make([]PengajuanIzinResponse, 0, len(izins))
	for _, izin := range izins {
		responses = append(responses, ToPengajuanIzinResponse(izin))
	}
	return responses
}



func ToPengajuanIzinDetailResponse(izin models.PengajuanIzin) PengajuanIzinDetailResponse {
	comments := make([]UserCommentPengajuanIzin, 0, len(izin.Comments))
	for _, comment := range izin.Comments {
		comments = append(comments, UserCommentPengajuanIzin{
			ID: comment.ID,
			Komentator: KomentatorPengajuanIzin{
				ID:        comment.UserID,
				FirstName: comment.User.FirstName,
				LastName:  comment.User.LastName,
				Role:      comment.User.Role.Name,
			},
			Comment:        comment.Comment,
			StatusPengajuan: string(comment.Status),
		})
	}

	return PengajuanIzinDetailResponse{
		ID:             izin.ID,
		TanggalMulai:   izin.TanggalMulai.Format("02 January 2006"),
		TanggalSelesai: izin.TanggalSelesai.Format("02 January 2006"),
		JenisIzin:      string(izin.JenisIzin),
		AlasanIzin:     izin.AlasanIzin,
		KeteranganIzin: izin.KeteranganIzin,
		Status:         string(izin.Status),
		PengajuanUser: UserPengajuanIzin{
			ID:        izin.User.ID,
			FirstName: izin.User.FirstName,
			LastName:  izin.User.LastName,
			Email:     izin.User.Email,
			ProfilePicture: izin.User.ProfilePicture,
		},
		Comments:  comments,
		CreatedAt: izin.CreatedAt,
		UpdatedAt: izin.UpdatedAt,
		DeletedAt: &izin.DeletedAt.Time,
	}
}