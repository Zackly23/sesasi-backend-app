package seeders

import (
	"math/rand"
	"fmt"
	"sesasi-backend-app/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) error {
	// daftar user default yang mau dibuat
	defaultUsers := []struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
		RoleName  models.UserRole
	}{
		{"Admin", "User", "admin123@gmail.com", "admin123", models.RoleAdmin},

		{"David", "Dwi Nugroho", "daviddwinugraha2@gmail.com", "david123", models.RoleVerifikator},
		{"Ani", "Wijaya", "aniwijaya@gmail.com", "ani123", models.RoleVerifikator},
		{"Rudi", "Hartono", "rudihartono@gmail.com", "rudi123", models.RoleVerifikator},

		// User biasa (4 orang)
		{"Budi", "Santoso", "budisantoso@gmail.com", "budi123", models.RoleUser},
		{"Siti", "Aisyah", "sitiaisyah@gmail.com", "siti123", models.RoleUser},
		{"Andi", "Saputra", "andisaputra@gmail.com", "andi123", models.RoleUser},
		{"Maya", "Putri", "mayaputri@gmail.com", "maya123", models.RoleUser},
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	verifiedOpts := []bool{
		true, false,
	}

	for _, u := range defaultUsers {
		var count int64
		db.Model(&models.User{}).Where("email = ?", u.Email).Count(&count)
		if count > 0 {
			fmt.Printf("User %s already exists, skipping...\n", u.Email)
			continue
		}

		// cek / buat role
		var role models.Role
		if err := db.Where("name = ?", u.RoleName).First(&role).Error; err != nil {
			role = models.Role{
				ID:   uuid.New(),
				Name: string(u.RoleName),
			}
			db.Create(&role)
		}

		// hash password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

		var verif bool
		if u.FirstName == "Budi" {
			verif = false
		} else if u.RoleName == "user" {
			verif = verifiedOpts[r.Intn(len(verifiedOpts))] // random true/false
		} else {
			verif = true
		}
			

		var userId uuid.UUID		
		if u.FirstName == "Budi" {
			userId = uuid.MustParse("47ee3220-ed79-4fce-a5ed-806f7813d302")
		} else {
			userId = uuid.New()
		}		
		
		// buat user
		user := models.User{
			ID:        userId,
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			Email:      u.Email,
			Password:   string(hashedPassword),
			IsVerified: verif,
			RoleID:     role.ID,
		}

		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("Error creating user %s: %v\n", u.Email, err)
		} else {
			fmt.Printf("User %s created successfully\n", u.Email)
		}
	}

	return nil
}
