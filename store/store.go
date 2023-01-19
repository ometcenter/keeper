package store

import (
	"database/sql"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBmain *sql.DB

var GormDB *gorm.DB

var PoolDB *DBPool

type DBPool struct {
	Mutex sync.Mutex
	Pool  map[string]*sql.DB
}

func InitBD(DatabaseURL string) error {
	var err error
	DBmain, err = sql.Open("postgres", DatabaseURL)
	if err != nil {
		return err
	}

	err = DBmain.Ping()
	if err != nil {
		return err
	}

	PoolDB = &DBPool{
		Pool: make(map[string]*sql.DB),
	}

	PoolDB.Mutex.Lock()
	defer PoolDB.Mutex.Unlock()

	PoolDB.Pool[DatabaseURL] = DBmain

	return nil
}

func GetDB(DSN string) (*sql.DB, error) {

	// db, err = sql.Open("postgres", fmt.Sprintf(
	// 	"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
	// 	configStruct[dbAlias].DbHost,
	// 	configStruct[dbAlias].DbPort,
	// 	configStruct[dbAlias].DbUser,
	// 	configStruct[dbAlias].DbName,
	// 	configStruct[dbAlias].DbPass)

	db, ok := PoolDB.Pool[DSN]
	if ok {
		return db, nil
	}

	var err error
	db, err = sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	PoolDB.Mutex.Lock()
	defer PoolDB.Mutex.Unlock()

	PoolDB.Pool[DSN] = db

	return db, nil
}

func InitgORM(DBIn *sql.DB) error {

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: DBIn,
	}), &gorm.Config{})

	// TODO: Извлекать из Gorm существующее подключение необязательно, потому что мы сейчас используем
	// GORM поверх существующего глобального подключения.
	DB, err := gormDB.DB()
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	//gormDB.AutoMigrate(&SettingsJobsTest{})
	// // Add name field
	// // TODO: Почему-то поле не добавляется.
	// gormDB.Migrator().AddColumn(&SettingsJobsTest{}, "name2")
	// if err != nil {
	// 	return err
	// }
	// // Check column exists
	// fmt.Println(gormDB.Migrator().HasColumn(&SettingsJobsTest{}, "name2"))
	// // Returns current using database name
	// fmt.Println(gormDB.Migrator().CurrentDatabase())

	//if config.Conf.GrabPasswordFromMail {

	// DBMainAnalytics, err := shareStore.GetDB(config.Conf.DatabaseURLMainAnalytics)
	// if err != nil {
	// 	return err
	// }

	// gormDBMainAnalytics, err := gorm.Open(postgres.New(postgres.Config{
	// 	Conn: DBMainAnalytics,
	// }), &gorm.Config{})

	// gormDBMainAnalytics.AutoMigrate(&models.EkisAreas{})
	// //gormDBMainAnalytics.Table("dit_ekis_areas").AutoMigrate(&models.EkisAreas{})
	// gormDBMainAnalytics.AutoMigrate(&modelsShare.EkisOrganizationDesctiption{})
	// //gormDBMainAnalytics.Table("dit_ekis_organization_desctiptions").AutoMigrate(&modelsShare.EkisOrganizationDesctiption{})
	// gormDBMainAnalytics.AutoMigrate(&modelsShare.OrganizationRegistrationInformation{})
	// gormDBMainAnalytics.AutoMigrate(&modelsShare.EkisOrganizationAddresses{})
	//}

	//gormDB.AutoMigrate(&models.LkUsers{})

	// err = TestCase(DB)
	// if err != nil {
	// 	return err
	// }

	// err = TestCase2(DB)
	// if err != nil {
	// 	return err
	// }

	// TODO: Переделать красиво глобальную переменную.
	GormDB = gormDB

	return nil

}
