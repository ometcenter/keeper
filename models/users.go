package models

import (
	"time"

	"gorm.io/gorm"
)

type LkUsers struct {
	gorm.Model
	FullName     string
	UserID       string // Он же collaborator_id
	Login        string `json:"login"`
	Password     string `json:"Password"`
	HashPassword string
	SecretJWT    string //`gorm:"index:idx_lk_users_jw_ttoken,type:btree"`
	//JWTtoken        string
	//JWTExp          int64
	ExpSec          int64
	Role            string
	InsuranceNumber string
	Email           string
	Status          string //Уволен и т.д
	DateDismissals  time.Time
	Blocked         bool
	Notes           string
}
