package models

type UserStatus string
type UserRole string
type TypePengajuanIzin string 
type PengajuanIzinStatus string

const (
	// status user
	StatusVerified   UserStatus = "verified"
	StatusUnverified UserStatus = "unverified"
	
	// role user
	RoleAdmin       UserRole = "admin"
	RoleVerifikator UserRole = "verifikator"
	RoleUser        UserRole = "user"

	// status pengajuan
	IzinDiajukan    PengajuanIzinStatus = "diajukan"
	IzinDitolak PengajuanIzinStatus = "ditolak"
	IzinDiterima  PengajuanIzinStatus = "diterima"
	IzinDiBatalkan PengajuanIzinStatus = "dibatalkan"
	IzinDirevisi  PengajuanIzinStatus = "revisi"

	// jenis pengajuan izin
	IzinSakit TypePengajuanIzin = "sakit"
	IzinCuti  TypePengajuanIzin = "cuti"
	IzinLibur TypePengajuanIzin = "libur"
	IzinLainnya TypePengajuanIzin = "lainnya"
)
