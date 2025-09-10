package utils

import "sesasi-backend-app/models"

func IsValidPengajuanStatus(status string) bool {
	switch models.PengajuanIzinStatus(status) {
	case models.IzinDiajukan, models.IzinDitolak, models.IzinDiterima, models.IzinDiBatalkan, models.IzinDirevisi:
		return true
	}
	return false
}

func IsValidTypePengajuanIzin(jenis string) bool {
	switch models.TypePengajuanIzin(jenis) {
	case models.IzinCuti, models.IzinLibur, models.IzinSakit, models.IzinLainnya:
		return true
	}
	return false
}

func IsValidUserRole(role string) bool {
	switch models.UserRole(role) {
	case models.RoleUser, models.RoleAdmin, models.RoleVerifikator:
		return true
	}
	return false
}

func IsValidUserStatus(status string) bool {
	switch models.UserStatus(status) {
	case models.StatusVerified,models.StatusUnverified:
		return true
	}
	return false
}