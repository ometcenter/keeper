package models

type SettingsJobsV2 struct {
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
