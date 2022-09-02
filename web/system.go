package web

import (
	"database/sql"

	"github.com/ometcenter/keeper/models"
)

func WriteSettingsJob(JobId string, DB *sql.DB, SettingsJobs models.SettingsJobs) error {

	//fmt.Println(SettingsJobs)

	// TODO: Удалить после тестирования, сохранения данных в поле pg JSON типа
	// На данный момент не требуется т.к мы просто записываем byte в поле с JSON
	// var SettingsJobSliceQueryToBI models.SettingsJobSliceQueryToBI
	// //err = json.Unmarshal([]byte(SettingsJobs.JSONString), &SettingsJobSliceQueryToBI)
	// err = json.Unmarshal([]byte(SettingsJobs.JSONString), &SettingsJobSliceQueryToBI)
	// if err != nil {
	// 	c.String(http.StatusBadRequest, err.Error())
	// 	log.Impl.Error(err.Error())
	// 	return
	// }

	//// Вариант 1
	// var argsUpdate []interface{}
	// argsUpdate = append(argsUpdate, SettingsJobs.JobID)
	// argsUpdate = append(argsUpdate, SettingsJobs.JSONString)

	// result, err := store.DB.Exec(`UPDATE settings_jobs SET job_id = $1, json_string = $2
	// 	WHERE job_id = $1;`, argsUpdate...)

	// if err != nil {
	// 	c.String(http.StatusBadRequest, err.Error())
	// 	log.Impl.Error(err.Error())
	// 	return
	// }

	// LastInsertId, _ := result.LastInsertId()
	// RowsAffected, _ := result.RowsAffected()

	// fmt.Println("LastInsertId: ", LastInsertId)
	// fmt.Println("RowsAffected: ", RowsAffected)

	// // Если не обновленно не одной записи, значит это новая запись и ее надо добавить
	// if RowsAffected == 0 {
	// 	var argsInsert []interface{}
	// 	argsInsert = append(argsInsert, SettingsJobs.JobID)
	// 	argsInsert = append(argsInsert, SettingsJobs.JSONString)

	// 	_, err := store.DB.Exec(`INSERT INTO settings_jobs (job_id, json_string)
	// 	VALUES($1, $2);`, argsInsert...)

	// 	if err != nil {
	// 		c.String(http.StatusBadRequest, err.Error())
	// 		log.Impl.Error(err.Error())
	// 		return
	// 	}

	// }

	// Вариант 2
	var argsquery []interface{}
	argsquery = append(argsquery, JobId)

	//queryAllColumns := `SELECT * FROM _jobs WHERE job_id = $1`

	//rows, err := store.DB.Query(queryAllColumns, argsquery...)
	//if err != nil {
	//	c.String(http.StatusBadRequest, err.Error())
	//	log.Impl.Error(err.Error())
	//	return
	//}

	//defer rows.Close()

	// flag := false
	// for rows.Next() {
	// 	flag = true
	// 	break
	// }

	var counter int
	DB.QueryRow("SELECT count(*) FROM settings_jobs WHERE job_id = $1", argsquery...).Scan(&counter)

	if counter != 0 {
		//if flag == true {

		var argsUpdate []interface{}
		argsUpdate = append(argsUpdate, SettingsJobs.JobID)
		argsUpdate = append(argsUpdate, SettingsJobs.JSONString)
		argsUpdate = append(argsUpdate, SettingsJobs.CodeExternal)
		argsUpdate = append(argsUpdate, SettingsJobs.NameExternal)
		argsUpdate = append(argsUpdate, SettingsJobs.TableName)
		argsUpdate = append(argsUpdate, SettingsJobs.UseRemoteCollection)
		argsUpdate = append(argsUpdate, SettingsJobs.ConfigName)
		//argsUpdate = append(argsUpdate, SettingsJobSliceQueryToBI)
		argsUpdate = append(argsUpdate, []byte(SettingsJobs.JSONString))

		// _, err := store.DB.Exec(`UPDATE settings_jobs SET job_id = $1, json_string = $2, code_external = $3,
		// name_external = $4, table_name = $5, use_remote_collection = $6, config_name = $7 WHERE job_id = $1;`, argsUpdate...)
		_, err := DB.Exec(`UPDATE settings_jobs SET job_id = $1, json_string = $2, code_external = $3, 
		 name_external = $4, table_name = $5, use_remote_collection = $6, config_name = $7, json_byte = $8 WHERE job_id = $1;`, argsUpdate...)

		if err != nil {
			return err
		}

		//LastInsertId, _ := result.LastInsertId()
		//RowsAffected, _ := result.RowsAffected()

		//fmt.Println("LastInsertId: ", LastInsertId)
		//fmt.Println("RowsAffected: ", RowsAffected)

	} else {

		var argsInsert []interface{}
		argsInsert = append(argsInsert, SettingsJobs.JobID)
		argsInsert = append(argsInsert, SettingsJobs.JSONString)
		argsInsert = append(argsInsert, SettingsJobs.CodeExternal)
		argsInsert = append(argsInsert, SettingsJobs.NameExternal)
		argsInsert = append(argsInsert, SettingsJobs.TableName)
		argsInsert = append(argsInsert, SettingsJobs.UseRemoteCollection)
		argsInsert = append(argsInsert, SettingsJobs.ConfigName)
		//argsInsert = append(argsInsert, SettingsJobSliceQueryToBI)
		argsInsert = append(argsInsert, []byte(SettingsJobs.JSONString))

		// _, err := store.DB.Exec(`INSERT INTO settings_jobs (job_id, json_string, code_external, name_external, table_name, use_remote_collection, config_name)
		// VALUES($1, $2, $3, $4, $5, $6, $7);`, argsInsert...)

		_, err := DB.Exec(`INSERT INTO settings_jobs (job_id, json_string, code_external, name_external, table_name, 
			use_remote_collection, config_name, json_byte)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8);`, argsInsert...)

		if err != nil {
			return err
		}

	}

	// if config.Conf.UseRedis {
	// 	err := store.SaveOneSettingsToRedis(SettingsJobs.JobID, SettingsJobs.JSONString)
	// 	if err != nil {
	// 		c.String(http.StatusBadRequest, err.Error())
	// 		log.Impl.Error(err.Error())
	// 		return
	// 	}
	// }

	return nil

}
