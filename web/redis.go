package web

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	libraryGoRedis "github.com/go-redis/redis/v8"
	log "github.com/ometcenter/keeper/logging"
	shareRedis "github.com/ometcenter/keeper/redis"
	utilityShare "github.com/ometcenter/keeper/utility"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////// Заполнение кэша Redis //////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

func FillDataToRedisSalary(RedisDB int, DB *sql.DB, RedisClient *libraryGoRedis.Client) error {

	BeginTime := time.Now()

	var argsquery1 []interface{}
	queryAllColumns := `select
	collaborators_posle.collaborator_id
from
	collaborators_posle as collaborators_posle
where 
	status <> 'Увольнение'`
	//where
	//area = '6083'
	//limit 100`

	rows1, err := DB.Query(queryAllColumns, argsquery1...)
	if err != nil {
		return err
	}

	defer rows1.Close()

	var collaborator_idSlice []string
	for rows1.Next() {
		var collaborator_id string
		err = rows1.Scan(&collaborator_id)
		if err != nil {
			//t.Fatalf("Scan: %v", err)
			//AnswerVlasov := AnswerVlasov{false, nil, &ErrorVlasov{http.StatusInternalServerError, err.Error()}}
			return err
		}
		collaborator_idSlice = append(collaborator_idSlice, collaborator_id)
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//InsuranceNumber := "" //c.Params.ByName("worker_id")

	UseYearFilter := false
	yearFilter := "2022"
	if yearFilter != "" {
		UseYearFilter = true
	}

	// DB, err := shareStore.GetDB(config.Conf.DatabaseURLMainAnalytics)
	// if err != nil {
	// 	//AnswerVlasov := AnswerVlasov{false, nil, &ErrorVlasovNew{http.StatusInternalServerError, err.Error()}}
	// 	return err
	// }

	err = shareRedis.FlushdbLibraryGoRedis(RedisClient, RedisDB)
	if err != nil {
		return err
	}

	for _, item := range collaborator_idSlice {

		workerID := item

		var BudgetStat interface{}
		BudgetStat, err = V1BudgetStatGeneral(workerID, UseYearFilter, yearFilter, RedisClient)
		if err != nil {
			BudgetStat = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		}

		//panic occurred in main: redis: can't marshal store.GetPersonalInfoResponds (implement encoding.BinaryMarshaler)
		byteResult, err := json.Marshal(BudgetStat)
		if err != nil {
			return err
		}

		//err = RedisClient.Set(ctxRedis, r.InsuranceNumber, r, 0).Err()
		err = shareRedis.SetLibraryGoRedis(RedisClient, item+yearFilter, byteResult, RedisDB, 0)
		if err != nil {
			return err
		}

		// err = RedisClient.Set(ctxRedis, item+yearFilter, byteResult, 0).Err()
		// if err != nil {
		// 	return err
		// }

	}

	// var AnswerVlasov AnswerVlasov
	// AnswerVlasov.Status = true
	// //AnswerVlasov.Data = ColumnsStructSlice
	// AnswerVlasov.Data = V1BudgetStatGroupRespondsSlice
	// AnswerVlasov.Error = nil
	// c.JSON(http.StatusOK, AnswerVlasov)

	endTime := time.Now()
	elapsedSave := endTime.Sub(BeginTime)

	log.Impl.Infof("Время заполнения кэша Redis по области %d = Старт: %s Конец: %s Продолжительность: %-8v (-) %v\n", RedisDB,
		BeginTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), elapsedSave, utilityShare.ShortDur(elapsedSave))

	return nil

}

func FillDataToRedisVacation(RedisDB int, DB *sql.DB, RedisClient *libraryGoRedis.Client) error {

	BeginTime := time.Now()

	var argsquery1 []interface{}
	queryAllColumns := `select
	collaborators_posle.collaborator_id
from
	collaborators_posle as collaborators_posle
where 
	status <> 'Увольнение'`
	//where
	//area = '6083'
	//limit 100`

	rows1, err := DB.Query(queryAllColumns, argsquery1...)
	if err != nil {
		return err
	}

	defer rows1.Close()

	var collaborator_idSlice []string
	for rows1.Next() {
		var collaborator_id string
		err = rows1.Scan(&collaborator_id)
		if err != nil {
			//t.Fatalf("Scan: %v", err)
			//AnswerVlasov := AnswerVlasov{false, nil, &ErrorVlasov{http.StatusInternalServerError, err.Error()}}
			return err
		}
		collaborator_idSlice = append(collaborator_idSlice, collaborator_id)
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//InsuranceNumber := "" //c.Params.ByName("worker_id")

	UseYearFilter := false
	// yearFilter := c.Query("year")
	// if yearFilter != "" {
	// 	UseYearFilter = true
	// }

	yearFilterFrom := "2022"
	yearFilterTo := "2022"
	if yearFilterFrom != "" && yearFilterTo != "" {
		UseYearFilter = true
	}

	// DB, err := shareStore.GetDB(config.Conf.DatabaseURLMainAnalytics)
	// if err != nil {
	// 	//AnswerVlasov := AnswerVlasov{false, nil, &ErrorVlasovNew{http.StatusInternalServerError, err.Error()}}
	// 	return err
	// }

	err = shareRedis.FlushdbLibraryGoRedis(RedisClient, RedisDB)
	if err != nil {
		return err
	}

	for _, item := range collaborator_idSlice {

		workerID := item
		//?from=2020&to=2023

		var HolidayStat interface{}
		HolidayStat, err = V1HolidayStatGeneral(workerID, UseYearFilter, yearFilterFrom, yearFilterTo, RedisClient)
		if err != nil {
			HolidayStat = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		}

		//panic occurred in main: redis: can't marshal store.GetPersonalInfoResponds (implement encoding.BinaryMarshaler)
		byteResult, err := json.Marshal(HolidayStat)
		if err != nil {
			return err
		}

		err = shareRedis.SetLibraryGoRedis(RedisClient, item+yearFilterFrom+yearFilterTo, byteResult, RedisDB, 0)
		if err != nil {
			return err
		}

	}

	endTime := time.Now()
	elapsedSave := endTime.Sub(BeginTime)

	log.Impl.Infof("Время заполнения кэша Redis по области %d = Старт: %s Конец: %s Продолжительность: %-8v (-) %v\n", RedisDB,
		BeginTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), elapsedSave, utilityShare.ShortDur(elapsedSave))

	return nil

}

func FillDataToRedisAllInformation(RedisDB int, DB *sql.DB, RedisClient *libraryGoRedis.Client) error {

	BeginTime := time.Now()

	var argsquery1 []interface{}
	queryAllColumns := `select
	collaborators_posle.collaborator_id
from
	collaborators_posle as collaborators_posle
where 
	status <> 'Увольнение' and
	area = '6083'`
	//limit 100`

	rows1, err := DB.Query(queryAllColumns, argsquery1...)
	if err != nil {
		return err
	}

	defer rows1.Close()

	var collaborator_idSlice []string
	for rows1.Next() {
		var collaborator_id string
		err = rows1.Scan(&collaborator_id)
		if err != nil {
			//t.Fatalf("Scan: %v", err)
			//AnswerVlasov := AnswerVlasov{false, nil, &ErrorVlasov{http.StatusInternalServerError, err.Error()}}
			return err
		}
		collaborator_idSlice = append(collaborator_idSlice, collaborator_id)
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//InsuranceNumber := "" //c.Params.ByName("worker_id")

	UseYearFilter := false
	// yearFilter := c.Query("year")
	// if yearFilter != "" {
	// 	UseYearFilter = true
	// }

	yearFilterFrom := "2022"
	yearFilterTo := "2022"
	yearFilter := "2022"
	if yearFilterFrom != "" && yearFilterTo != "" {
		UseYearFilter = true
	}

	// DB, err := shareStore.GetDB(config.Conf.DatabaseURLMainAnalytics)
	// if err != nil {
	// 	//AnswerVlasov := AnswerVlasov{false, nil, &ErrorVlasovNew{http.StatusInternalServerError, err.Error()}}
	// 	return err
	// }

	err = shareRedis.FlushdbLibraryGoRedis(RedisClient, RedisDB)
	if err != nil {
		return err
	}

	for _, item := range collaborator_idSlice {

		workerID := item
		//?from=2020&to=2023

		var AllInformationV1 interface{}
		AllInformationV1, err = AllInformationV1General(workerID, UseYearFilter, yearFilter, yearFilterFrom, yearFilterTo, shareRedis.RedisClientGoRedisLibrary)
		if err != nil {
			AllInformationV1 = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		}

		//panic occurred in main: redis: can't marshal store.GetPersonalInfoResponds (implement encoding.BinaryMarshaler)
		byteResult, err := json.Marshal(AllInformationV1)
		if err != nil {
			return err
		}

		// TODO: Внтури сборки частных частей кеша происходит изменения на другую область, поэтому тут меняем.
		err = shareRedis.SelectLibraryGoRedis(RedisClient, RedisDB)
		if err != nil {
			return err
		}

		err = shareRedis.SetLibraryGoRedis(RedisClient, item+yearFilterFrom+yearFilterTo, byteResult, RedisDB, 0)
		if err != nil {
			return err
		}

	}

	endTime := time.Now()
	elapsedSave := endTime.Sub(BeginTime)

	log.Impl.Infof("Время заполнения кэша Redis по области %d = Старт: %s Конец: %s Продолжительность: %-8v (-) %v\n", RedisDB,
		BeginTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), elapsedSave, utilityShare.ShortDur(elapsedSave))

	return nil

}

func FillDataToRedisJobPlace(RedisDB int, DB *sql.DB, RedisClient *libraryGoRedis.Client) error {

	BeginTime := time.Now()

	var argsquery1 []interface{}
	queryAllColumns := `select
	collaborators_posle.collaborator_id
from
	collaborators_posle as collaborators_posle
where 
	status <> 'Увольнение'`
	//where
	//area = '6083'
	//limit 100`

	rows1, err := DB.Query(queryAllColumns, argsquery1...)
	if err != nil {
		return err
	}

	defer rows1.Close()

	var collaborator_idSlice []string
	for rows1.Next() {
		var collaborator_id string
		err = rows1.Scan(&collaborator_id)
		if err != nil {
			//t.Fatalf("Scan: %v", err)
			//AnswerVlasov := AnswerVlasov{false, nil, &ErrorVlasov{http.StatusInternalServerError, err.Error()}}
			return err
		}
		collaborator_idSlice = append(collaborator_idSlice, collaborator_id)
	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	err = shareRedis.FlushdbLibraryGoRedis(RedisClient, RedisDB)
	if err != nil {
		return err
	}

	for _, item := range collaborator_idSlice {

		workerID := item

		var V3JobPlaces interface{}
		V3JobPlaces, err = V3JobPlacesGeneral(workerID, RedisClient)
		if err != nil {
			V3JobPlaces = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		}

		//panic occurred in main: redis: can't marshal store.GetPersonalInfoResponds (implement encoding.BinaryMarshaler)
		byteResult, err := json.Marshal(V3JobPlaces)
		if err != nil {
			return err
		}

		//err = RedisClient.Set(ctxRedis, r.InsuranceNumber, r, 0).Err()
		err = shareRedis.SetLibraryGoRedis(RedisClient, item, byteResult, RedisDB, 0)
		if err != nil {
			return err
		}

		// err = RedisClient.Set(ctxRedis, item+yearFilter, byteResult, 0).Err()
		// if err != nil {
		// 	return err
		// }

	}

	// var AnswerVlasov AnswerVlasov
	// AnswerVlasov.Status = true
	// //AnswerVlasov.Data = ColumnsStructSlice
	// AnswerVlasov.Data = V1BudgetStatGroupRespondsSlice
	// AnswerVlasov.Error = nil
	// c.JSON(http.StatusOK, AnswerVlasov)

	endTime := time.Now()
	elapsedSave := endTime.Sub(BeginTime)

	log.Impl.Infof("Время заполнения кэша Redis по области %d = Старт: %s Конец: %s Продолжительность: %-8v (-) %v\n", RedisDB, BeginTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), elapsedSave, utilityShare.ShortDur(elapsedSave))

	return nil

}
