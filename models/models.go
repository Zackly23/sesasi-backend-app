package models

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Permission{},
		&Role{},
		&RolePermission{},
		&PrivateAccessToken{},
		&PengajuanIzin{},
		&PengajuanIzinComment{},
	}
}