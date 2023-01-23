package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// import (
// 	"time"

// 	"github.com/ometcenter/keeper/web"
// 	"gorm.io/datatypes"
// 	"gorm.io/gorm"
// )

// TODO: Used for migration. Need a replacement for web.LkUsers
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
	Source          string
	PersonJSONByte  datatypes.JSON
	//Person                         V1ActiveWorkers `gorm:"-"`
	AdditionalSettingsUserJSONByte datatypes.JSON
	//AdditionalSettingsUser         AdditionalSettingsUser `gorm:"-"`
	Notes string
}
