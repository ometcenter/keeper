package store

import (
	"database/sql"
	"sync"
)

var DBmain *sql.DB

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
