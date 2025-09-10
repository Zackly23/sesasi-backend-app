package seeders 

import (
	"fmt"
	"log"
	"sesasi-backend-app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func PermissionSeeder(db *gorm.DB) error {
	var count int64
	db.Model(&models.Permission{}).Count(&count)
	if count > 0 {
		log.Println("Permission data already exists, skipping seeder.")
		return nil
	}

	permissions := []models.Permission{
		{ID: uuid.New(), Type: "read", Name: "view_profile"},
		{ID: uuid.New(), Type: "write", Name: "edit_profile"},
		{ID: uuid.New(), Type: "verify", Name: "verify_data"},
		{ID: uuid.New(), Type: "delete", Name: "delete_user"},
	}

	if err := db.Create(&permissions).Error; err != nil {
		return fmt.Errorf("failed seeding permissions: %v", err)
	}

	log.Println("Permission seeding success")
	return nil
}

func RolePermissionSeeder(db *gorm.DB) error {
    roles := []models.Role{
        {ID: uuid.New(), Name: string(models.RoleAdmin), Description: "Administrator role"},
        {ID: uuid.New(), Name: string(models.RoleVerifikator), Description: "Verifikator role"},
        {ID: uuid.New(), Name: string(models.RoleUser), Description: "Regular user role"},
    }

    // Insert roles jika belum ada
    if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error; err != nil {
        return fmt.Errorf("failed seeding roles: %v", err)
    }

    // Ambil ulang role dari DB
    var dbRoles []models.Role
    if err := db.Find(&dbRoles).Error; err != nil {
        return fmt.Errorf("failed fetch roles: %v", err)
    }

    // Ambil semua permission
    var permissions []models.Permission
    if err := db.Find(&permissions).Error; err != nil {
        return fmt.Errorf("failed fetch permissions: %v", err)
    }

    // Mapping
    var rolePermissions []models.RolePermission
    for _, role := range dbRoles {
        for _, perm := range permissions {
            switch role.Name {
            case string(models.RoleAdmin):
                rolePermissions = append(rolePermissions, models.RolePermission{
                    RoleID: role.ID, PermissionID: perm.ID, CreatedAt: time.Now(),
                })
            case string(models.RoleVerifikator):
                if perm.Name == "verify_data" {
                    rolePermissions = append(rolePermissions, models.RolePermission{
                        RoleID: role.ID, PermissionID: perm.ID, CreatedAt: time.Now(),
                    })
                }
            case string(models.RoleUser):
                if perm.Name == "view_profile" {
                    rolePermissions = append(rolePermissions, models.RolePermission{
                        RoleID: role.ID, PermissionID: perm.ID, CreatedAt: time.Now(),
                    })
                }
            }
        }
    }

    // Insert tanpa duplikat (butuh unique constraint di DB)
    if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rolePermissions).Error; err != nil {
        return fmt.Errorf("failed seeding role_permissions: %v", err)
    }

    log.Println("RolePermission seeding success")
    return nil
}

