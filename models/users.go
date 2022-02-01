package models

import "gorm.io/gorm"

type LkUsers struct {
	gorm.Model
	FullName        string
	UserID          string
	Login           string `json:"login"`
	Password        string `json:"Password"`
	HashPassword    string
	SecretJWT       string `gorm:"index:idx_lk_users_jw_ttoken,type:btree"`
	JWTtoken        string
	JWTExp          int64
	Role            string
	InsuranceNumber string
	Email           string
	Notes           string
}
