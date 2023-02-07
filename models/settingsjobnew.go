package models

import (
	"database/sql"

	"encoding/json"

	"gorm.io/datatypes"
)

type SettingsJobsV2 struct {
	JobID        string `json:"jobID"`        //ИдентификаторЗадания
	CodeExternal string `json:"codeExternal"` //Код1С
	NameExternal string `json:"nameExternal"` //Наименование1С
	TableName    string `json:"tableName"`    //ИмяТаблицы
	//UseRemoteCollection    bool   `json:"УдаленныйСбор"`      //УдаленныйСбор
	//ConfigName             string `json:"ИмяКонфигурации"`    //ИмяКонфигурации
	TypeDataGetting  string `json:"typeDataGetting"`  //ВидПолученияДанных
	DataUploadMethod string `json:"dataUploadMethod"` // CпособЗагрузкиДанных

	//SettingsJobsJSONString string `json:"JSONСтрокаНастроек"`
	SettingsJobsJSONByte datatypes.JSON
}

func (S SettingsJobsAllV2) SaveSettingsToPg(DB *sql.DB) error {

	settingByte, err := json.Marshal(S)
	if err != nil {
		return err
	}

	//// Вариант 1
	// var argsUpdate []interface{}
	// argsUpdate = append(argsUpdate, SettingsJobs.JobID)

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
	argsquery = append(argsquery, S.JobID)

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
	DB.QueryRow("SELECT count(*) FROM settings_jobs_v2 WHERE job_id = $1", argsquery...).Scan(&counter)

	if counter != 0 {
		//if flag == true {

		var argsUpdate []interface{}
		argsUpdate = append(argsUpdate, S.JobID)
		argsUpdate = append(argsUpdate, S.CodeExternal)
		argsUpdate = append(argsUpdate, S.NameExternal)
		argsUpdate = append(argsUpdate, S.TableName)
		argsUpdate = append(argsUpdate, S.TypeDataGetting)
		argsUpdate = append(argsUpdate, S.DataUploadMethod)
		argsUpdate = append(argsUpdate, settingByte)

		// _, err := store.DB.Exec(`UPDATE settings_jobs SET job_id = $1, json_string = $2, code_external = $3,
		// name_external = $4, table_name = $5, use_remote_collection = $6, config_name = $7 WHERE job_id = $1;`, argsUpdate...)
		_, err := DB.Exec(`UPDATE settings_jobs_v2 SET job_id = $1, code_external = $2, 
		 name_external = $3, table_name = $4, type_data_getting = $5, data_upload_method = $6,
		  settings_jobs_json_byte = $7 WHERE job_id = $1;`, argsUpdate...)

		if err != nil {
			return err
		}

		//LastInsertId, _ := result.LastInsertId()
		//RowsAffected, _ := result.RowsAffected()

		//fmt.Println("LastInsertId: ", LastInsertId)
		//fmt.Println("RowsAffected: ", RowsAffected)

	} else {

		var argsInsert []interface{}
		argsInsert = append(argsInsert, S.JobID)
		argsInsert = append(argsInsert, S.CodeExternal)
		argsInsert = append(argsInsert, S.NameExternal)
		argsInsert = append(argsInsert, S.TableName)
		argsInsert = append(argsInsert, S.TypeDataGetting)
		argsInsert = append(argsInsert, S.DataUploadMethod)
		argsInsert = append(argsInsert, settingByte)

		// _, err := store.DB.Exec(`INSERT INTO settings_jobs (job_id, json_string, code_external, name_external, table_name, use_remote_collection, config_name)
		// VALUES($1, $2, $3, $4, $5, $6, $7);`, argsInsert...)

		_, err := DB.Exec(`INSERT INTO settings_jobs_v2 (job_id, code_external, name_external, table_name, type_data_getting, data_upload_method, settings_jobs_json_byte)
		VALUES($1, $2, $3, $4, $5, $6, $7);`, argsInsert...)

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

type SettingsJobsAllV2 struct {
	JobID string `json:"jobID"` // ИдентификаторЗадания
	//SendUseREST                      bool                   `json:"ОтправлятьПоREST"`
	//RemoteCollect                    bool                   `json:"УдаленныйСбор"`
	TypeDataGetting  string `json:"typeDataGetting"`  // ВидПолученияДанных
	DataUploadMethod string `json:"dataUploadMethod"` // CпособЗагрузкиДанных
	//Portions                         int                    `json:"Порции"`
	//Query                            []Query                `json:"Запросы"`
	//AddParam                         AdditionParam          `json:"ДополнительныеПараметрыJSON"`
	//AddParamJSNOString               string                 `json:"JSONСтрокаДополнительныеПараметры"`
	//Connect                          Connect                `json:"ПараметрыПодключения"`
	//ConnectContur                    ConnectContur          `json:"ПараметрыПодключенияКонтура"`
	//ConnectBI1C                      ConnectBI1C            `json:"ПараметрыПодключенияBI1C"`
	//ConnectConturJSNOString          string                 `json:"JSONСтрокаПараметрыПодключенияКонтура"`
	Schedule                         ScheduleV2               `json:"schedule"`                      // РасписаниеПланировщика
	SaveResultToHistory              bool                     `json:"saveResultToHistory"`           //СохранятьРезультатВИсторию
	PublishTableToAPI                bool                     `json:"publishTableToAPI"`             // ПубликоватьТаблицуВAPI
	SaveToDataVisualizationSystem    bool                     `json:"saveToDataVisualizationSystem"` //СохранятьВСистемуВизуализацииДанных
	UseDataProcessingAlgorithms      bool                     `json:"useDataProcessingAlgorithms"`   // ИспользоватьАлгоритмыОбработкиДанных
	ListDataProcessingAlgorithms     []string                 `json:"listDataProcessingAlgorithms"`  // СписокАлгоритмовОбработкиДанных
	UseHandleAfterLoadAlgorithms     bool                     `json:"useHandleAfterLoadAlgorithms"`  //ИспользоватьАлгоритмыОбработкиДанныхПослеЗагрузки
	ListHandleAfterLoadAlgorithms    []string                 `json:"listHandleAfterLoadAlgorithms"` //СписокАлгоритмовОбработкиДанныхПослеЗагрузки
	Webhooks                         []string                 `json:"webhooks"`
	MappingForExcelArray             []MappingForExcelArrayV2 `json:"mappingForExcelArray"`             //СопоставлениеДляExcalМассив
	RuleExternalSource               string                   `json:"ruleExternalSource"`               //ПравилоВнешнийИсточник
	InternalProcessingExternalSource bool                     `json:"internalProcessingExternalSource"` //ВнутренняяОбработкаВнешнегоИсточника
	TableName                        string                   `json:"tableName"`                        //ИмяТаблицы
	ZipAnswer                        bool                     `json:"zipAnswer"`                        //СжиматьОтвет
	DSNconnection                    string                   `json:"dsnConnection"`                    // БазаСУБДDSN
	HashAnswer                       bool                     `json:"hashAnswer"`                       //ХешироватьРезультат
	CodeExternal                     string                   `json:"codeExternal"`                     // Внешний код задания
	NameExternal                     string                   `json:"nameExternal"`                     // Внешнее имя задания
}

type ScheduleV2 struct {
	JobID             string   `json:"jobID"`             //ИдентификаторЗадания
	UseRegulatoryTask bool     `json:"useRegulatoryTask"` //ИспользоватьРегламентноеЗадание
	UseCronSchedule   bool     `json:"useCronSchedule"`   //ИспользоватьCronРасписание
	SliceCronString   []string `json:"sliceCronString"`   //МассивСтрокCron
}

type MappingForExcelArrayV2 struct {
	NumberField int    `json:"numberField"` //НомерСтроки
	Name        string `json:"name"`        //Имя
}
