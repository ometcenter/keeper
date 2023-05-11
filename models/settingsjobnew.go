package models

import (
	"database/sql"
	"database/sql/driver"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"

	"encoding/json"

	"gorm.io/datatypes"
)

type SettingsJobsV2 struct {
	JobID        string `json:"jobID" gorm:"primaryKey"` //ИдентификаторЗадания
	CodeExternal string `json:"codeExternal"`            //Код1С
	NameExternal string `json:"nameExternal"`            //Наименование1С
	TableName    string `json:"tableName"`               //ИмяТаблицы
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
	QueryDetails []QueryV2 `json:"queryDetails"`
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
	SelectionFields                  []string                 `json:"selectionFields"`                  //ПоляОтбора
	ComparionFields                  []string                 `json:"comparionFields"`                  //ПоляСравнения
	UseCleaningFieldsBeforeLoading   string                   `json:"useCleaningFieldsBeforeLoading"`   // Использовать фильтр очистки данных
	CleaningFieldsBeforeLoading      string                   `json:"cleaningFieldsBeforeLoading"`      // Фильр очистки данных

}

func (S *SettingsJobsAllV2) Scan(value interface{}) (err error) {
	switch value.(type) {
	case string:
		err = json.Unmarshal([]byte(value.(string)), &S)
	case []byte:
		err = json.Unmarshal(value.([]byte), &S)
	default:
		return errors.New("incompatible type for skills")
	}
	if err != nil {
		return
	}

	return nil
}

func (S SettingsJobsAllV2) Value() (driver.Value, error) {
	return json.Marshal(S)
}

func (S *SettingsJobsAllV2) LoadSettingsFromPgByJobID(DB *sql.DB, JobID string) error {

	var argsquery []interface{}
	argsquery = append(argsquery, JobID)

	var LoadValue SettingsJobsAllV2
	err := DB.QueryRow("SELECT settings_jobs_json_byte FROM settings_jobs_v2 WHERE job_id = $1", argsquery...).Scan(&LoadValue)
	if err != nil {
		return err
	}

	// TODO: Переделать по нормальному эту конструкцию
	*S = LoadValue

	return nil
}

func (S *SettingsJobsAllV2) LoadSettingsFromPgByFileds(DB *sql.DB, FieldName string, Value interface{}) error {

	var argsquery []interface{}
	argsquery = append(argsquery, Value)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.Select("settings_jobs_json_byte").From("settings_jobs_v2").Where(sq.Eq{FieldName: Value})

	queryText, _, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	//fmt.Println(queryText)

	var LoadValue SettingsJobsAllV2
	err = DB.QueryRow(queryText, argsquery...).Scan(&LoadValue)
	if err != nil {
		return err
	}

	// TODO: Переделать по нормальному эту конструкцию
	*S = LoadValue

	return nil
}

type QueryV2 struct {
	QueryText string `json:"queryText"`
	//Base                                 string `json:"База"`
	Areas         string `json:"areas"`
	ExchangeJobID string `json:"exchangeJobID"`
	//UsedCalculatedFieldsInQueryParametrs bool   `json:"ИспользуетсяВычисляемыеПоляВПараметрахЗапроса"`
	QueryParams []QueryParams `json:"queryParams"`
}

type QueryParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func (SettingsJobsAllV2 *SettingsJobsAllV2) TranformToOldSettings() (QueryToBI, error) {

	var Q QueryToBI

	var AddParam AdditionParam
	AddParam.HashAnswer = SettingsJobsAllV2.HashAnswer
	AddParam.ZipAnswer = SettingsJobsAllV2.ZipAnswer

	var Options Options
	Options.TableName = SettingsJobsAllV2.TableName
	//SettingsJobsAllV2.CodeExternal
	Options.DSNconnection = SettingsJobsAllV2.DSNconnection
	Options.SelectionFields = SettingsJobsAllV2.SelectionFields
	Options.ComparionFields = SettingsJobsAllV2.ComparionFields

	AddParam.Options = &Options
	Q.AddParam = AddParam

	Q.DataUploadMethod = SettingsJobsAllV2.DataUploadMethod

	Q.InternalProcessingExternalSource = SettingsJobsAllV2.InternalProcessingExternalSource
	Q.JobID = SettingsJobsAllV2.JobID
	Q.ListDataProcessingAlgorithms = SettingsJobsAllV2.ListDataProcessingAlgorithms
	Q.ListHandleAfterLoadAlgorithms = SettingsJobsAllV2.ListHandleAfterLoadAlgorithms
	copier.Copy(&Q.MappingForExcelArray, &SettingsJobsAllV2.MappingForExcelArray)
	//SettingsJobsAllV2.NameExternal
	Q.PublishTableToAPI = SettingsJobsAllV2.PublishTableToAPI
	Q.RuleExternalSource = SettingsJobsAllV2.RuleExternalSource
	Q.SaveToDataVisualizationSystem = SettingsJobsAllV2.SaveToDataVisualizationSystem
	Q.TypeDataGetting = SettingsJobsAllV2.TypeDataGetting
	Q.UseDataProcessingAlgorithms = SettingsJobsAllV2.UseDataProcessingAlgorithms
	Q.UseHandleAfterLoadAlgorithms = SettingsJobsAllV2.UseHandleAfterLoadAlgorithms
	Q.Webhooks = SettingsJobsAllV2.Webhooks

	return Q, nil
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
