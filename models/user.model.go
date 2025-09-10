package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission
type Permission struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string         `json:"name" gorm:"type:varchar(100);uniqueIndex"`
	Type      string         `json:"type,omitempty" gorm:"type:varchar(50)"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Role
type Role struct {
	ID          uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string        `json:"name" gorm:"type:varchar(50);uniqueIndex"` // admin, verifikator, user
	Description string        `json:"description,omitempty"`
	Permissions []Permission  `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Pivot role_permissions (opsional kalau butuh timestamps)
type RolePermission struct {
    RoleID       uuid.UUID `gorm:"uniqueIndex:idx_role_perm"`
    PermissionID uuid.UUID `gorm:"uniqueIndex:idx_role_perm"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// User (ubah: simpan RoleID & preload Role)
type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	FirstName      string         `json:"first_name" gorm:"not null"`
	LastName       string         `json:"last_name" gorm:"not null"`
	UserName       string         `json:"user_name,omitempty"`
	Email          string         `json:"email" gorm:"uniqueIndex;not null"`
	Password       string         `json:"password" gorm:"not null"`
	Phone          *string        `json:"phone,omitempty" gorm:"uniqueIndex"`
	IsVerified         bool     `json:"is_verified,omitempty" gorm:"type:boolean;default:false"`
	ProfilePicture string         `json:"profile_picture,omitempty"`
	RoleID         uuid.UUID     `gorm:"type:uuid;index" json:"role_id,omitempty"`
	Role           Role          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
