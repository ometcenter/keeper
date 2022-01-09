package models

import "gorm.io/gorm"

type LkUsers struct {
	gorm.Model
	FullName        string
	UserID          string
	Login           string `json:"login"`
	Password        string `json:"Password"`
	HashPassword    string
	SecretJWT       string
	JWTtoken        string
	JWTExp          int64
	Role            string
	InsuranceNumber string
	Email           string
	Notes           string
}
