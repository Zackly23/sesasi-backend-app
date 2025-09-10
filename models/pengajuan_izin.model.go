package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PengajuanIzin struct {
	ID        uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID     `gorm:"type:uuid;index" json:"user_id"`
	User	  User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	TanggalMulai time.Time `json:"tanggal_mulai"`
	TanggalSelesai time.Time `json:"tanggal_selesai"`
	AlasanIzin string       `gorm:"type:text" json:"alasan_izin"`
	KeteranganIzin string  `gorm:"type:text" json:"keterangan_izin"`
	JenisIzin	 TypePengajuanIzin  `gorm:"type:varchar(50);default:'sakit'" json:"type"`
	Status    PengajuanIzinStatus `gorm:"type:varchar(50);default:'diajukan'" json:"status"`
	Comments   []*PengajuanIzinComment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comment,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type PengajuanIzinComment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	PengajuanIzinID uuid.UUID `gorm:"type:uuid;index" json:"pengajuan_izin_id"`
	PengajuanIzin PengajuanIzin `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"pengajuan_izin,omitempty"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	Comment   string    `json:"comment"`	
	Status   PengajuanIzinStatus `gorm:"type:varchar(50);default:'diajukan'" json:"status"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}