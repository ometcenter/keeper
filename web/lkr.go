package web

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	libraryGoRedis "github.com/go-redis/redis/v8"
	"github.com/ometcenter/keeper/config"
	log "github.com/ometcenter/keeper/logging"
	redis "github.com/ometcenter/keeper/redis"
	store "github.com/ometcenter/keeper/store"
	tree "github.com/ometcenter/keeper/tree"
	utilityShare "github.com/ometcenter/keeper/utility"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type V1HolidayStatResponds struct {
	PersonId  string `json:"personId"`
	Status    string `json:"status"`
	DataStart string `json:"dataStart"`
	DataEnd   string `json:"dataEnd"`
}

type V1HolidayAllStatResponds struct {
	Status       string `json:"status"`
	DateStart    string `json:"dateStart"`
	DateEnd      string `json:"dateEnd"`
	DocumentBase string `json:"documentBase"`
}

type V1HolidayAllStatRespondsForColleagues struct {
	CollaboratorId string `json:"collaboratorId"`
	FullName       string `json:"fullName"`
	Status         string `json:"status"`
	DateStart      string `json:"dateStart"`
	DateEnd        string `json:"dateEnd"`
	DocumentBase   string `json:"documentBase"`
}

type AverageSalary struct {
	Months       int     `json:"months"`
	Summ         float32 `json:"summ"`
	SummGross    float32 `json:"summGross"`
	Average      float32 `json:"average"`
	AverageGross float32 `json:"averageGross"`
	DaySum       float32 `json:"daySum"`
	DaySumGross  float32 `json:"daySumGross"`
}

type V1VacationSchedule struct {
	DataStart    string `json:"dataStart"`
	DataEnd      string `json:"dataEnd"`
	DaysNumber   int    `json:"daysNumber"`
	TypeVacation string `json:"typeVacation"`
}

type GetPersonalInfoResponds struct {
	PersonId         string
	InsuranceNumber  string
	Inn              string
	FullName         string
	Position         string
	OrganizationName string
	Status           string
}

type V1ActiveWorkers struct {
	PersonId          string    `json:"personId"`
	CollaboratorId    string    `json:"collaboratorId"`
	InsuranceNumber   string    `json:"insuranceNumber"`
	Inn               string    `json:"inn"`
	FullName          string    `json:"fullName"`
	Position          string    `json:"position"`
	OrganizationName  string    `json:"organizationName"`
	Status            string    `json:"status"`
	Email             string    `json:"email"`
	EmailEPS          string    `json:"emailEPS"`
	MobilePhone       string    `json:"mobilePhone"`
	WorkPhone         string    `json:"workPhone"`
	EmailArray        string    `json:"emailArray"`
	DateBirth         string    `json:"dateBirth"`
	BranchID          string    `json:"branchID"`
	BranchName        string    `json:"branchName"`
	LargeGroupOfPosts string    `json:"large_group_of_posts"`
	Position_tag      string    `json:"position_tag"`
	UpdatedAt         time.Time `json:"updatedAt"`
	CreatedAt         time.Time `json:"createdAt"`
	DateDismissals    time.Time `json:"dateDismissals"`
}

func (V1ActiveWorkers *V1ActiveWorkers) Scan(value interface{}) (err error) {
	switch value.(type) {
	case string:
		err = json.Unmarshal([]byte(value.(string)), &V1ActiveWorkers)
	case []byte:
		err = json.Unmarshal(value.([]byte), &V1ActiveWorkers)
	default:
		return errors.New("incompatible type for skills")
	}
	if err != nil {
		return
	}

	return nil
}

// type V1JobPlaces struct {
// 	PersonId          string    `json:"personId"`
// 	CollaboratorId    string    `json:"collaboratorId"`
// 	InsuranceNumber   string    `json:"insuranceNumber"`
// 	Inn               string    `json:"inn"`
// 	FullName          string    `json:"fullName"`
// 	Position          string    `json:"position"`
// 	OrganizationName  string    `json:"organizationName"`
// 	Status            string    `json:"status"`
// 	Email             string    `json:"email"`
// 	MobilePhone       string    `json:"mobilePhone"`
// 	WorkPhone         string    `json:"workPhone"`
// 	EmailArray        string    `json:"emailArray"`
// 	DateBirth         string    `json:"dateBirth"`
// 	BranchID          string    `json:"branchID"`
// 	BranchName        string    `json:"branchName"`
// 	LargeGroupOfPosts string    `json:"large_group_of_posts"`
// 	Position_tag      string    `json:"position_tag"`
// 	UpdatedAt         time.Time `json:"updatedAt"`
// }

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
	ExpSec                         int64
	Role                           string
	InsuranceNumber                string
	Email                          string
	Status                         string //Уволен и т.д
	DateDismissals                 time.Time
	Blocked                        bool
	Source                         string
	PersonJSONByte                 datatypes.JSON
	Person                         V1ActiveWorkers `gorm:"-"`
	AdditionalSettingsUserJSONByte datatypes.JSON
	AdditionalSettingsUser         AdditionalSettingsUser `gorm:"-"`
	Notes                          string
}

func (LkUsers *LkUsers) HashAndSalt() ([]byte, error) {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(LkUsers.Password), bcrypt.DefaultCost) //bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	LkUsers.HashPassword = string(hash)

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return hash, nil
}

func (LkUsers *LkUsers) ComparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false
	}
	//if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {

	return true
}

type AdditionalSettingsUser struct {
	AccessToSystemTables bool `json:"accessToSystemTables"`
}

func (AdditionalSettingsUser *AdditionalSettingsUser) Scan(value interface{}) (err error) {
	switch value.(type) {
	case string:
		err = json.Unmarshal([]byte(value.(string)), &AdditionalSettingsUser)
	case []byte:
		err = json.Unmarshal(value.([]byte), &AdditionalSettingsUser)
	default:
		return errors.New("incompatible type for skills")
	}
	if err != nil {
		return
	}

	return nil
}

type SalaryResponds struct {
	Area             string
	Snils            string
	DateRegistration string
	SettlementGroup  string
	CalculationType  string
	FullName         string
	DaysWorked       string
	HoursWorked      string
	OrganizationId   string
	Summ             int
}

type V1BudgetStatResponds struct {
	DateRegistration string `json:"dateRegistration"`
	SettlementGroup  string `json:"settlementGroup"`
	CalculationType  string `json:"calculationType"`
	//DaysWorked       int     `json:"daysWorked"`
	//HoursWorked      float32 `json:"hoursWorked"`
	Summ float32 `json:"summ"`
}

type V1BudgetStatGroupResponds struct {
	Total          float32                `json:"total"`
	TotalGross     float32                `json:"totalGross"`
	TotalDeduction float32                `json:"totalDeduction"`
	Month          int                    `json:"month"`
	DaysWorked     int                    `json:"daysWorked"`
	HoursWorked    float32                `json:"hoursWorked"`
	Items          []V1BudgetStatResponds `json:"items"`
}

type AllInformationV1Answer struct {
	HolidayStat              interface{} `json:"holidayStat"`
	BudgetStat               interface{} `json:"budgetStat"`
	JobPlaces                interface{} `json:"jobPlaces"`
	HolidayStatForColleagues interface{} `json:"holidayStatForColleagues"`
	GetBranchTree            interface{} `json:"getBranchTree"`
	AverageSalary            interface{} `json:"averageSalary"`
}

func AllInformationV1General(workerID string, UseYearFilter bool, yearFilter, yearFilterFrom, yearFilterTo string, RedisClient *libraryGoRedis.Client) (interface{}, error) {

	JSONString, err := redis.GetLibraryGoRedis(RedisClient, workerID+yearFilterFrom+yearFilterTo, 4)
	//if err != nil {
	if JSONString == "" {
		//log.Impl.Error(err.Error())
		// JSONString, err = store.GetSettingsByIdJobPg(JobIdParam)
		// if err != nil {
		// 	log.Impl.Error(err)
		// }

		// AnswerWebV1 := AnswerWebV1{false, store.DataAuthorizatioAnswer{}, ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		// c.JSON(http.StatusBadRequest, AnswerWebV1)
	} else {
		//c.Data(http.StatusOK, "application/json", []byte(JSONString))
		//c.JSON(http.StatusOK, JSONString)
		//return []byte(JSONString), nil

		var AnswerWebV1 AnswerWebV1
		if err := json.Unmarshal([]byte(JSONString), &AnswerWebV1); err != nil {
			return nil, err
		}

		return AnswerWebV1, nil
	}

	var AllInformationV1Answer AllInformationV1Answer

	var HolidayStat interface{}
	HolidayStat, err = V1HolidayStatGeneral(workerID, UseYearFilter, yearFilterFrom, yearFilterTo, RedisClient)
	if err != nil {
		HolidayStat = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
	}

	var BudgetStat interface{}
	BudgetStat, err = V1BudgetStatGeneral(workerID, UseYearFilter, yearFilter, RedisClient)
	if err != nil {
		BudgetStat = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
	}

	var JobPlaces interface{}
	JobPlaces, err = V3JobPlacesGeneral(workerID, RedisClient)
	if err != nil {
		JobPlaces = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
	}

	// startDateFilter := time.Now()
	// endDateFilter := startDateFilter.AddDate(0, 0, 7)
	startDateFilter := utilityShare.StartOfWeek(time.Now())
	//startDateFilter = startDateFilter.Truncate(24 * time.Hour)
	endDateFilter := startDateFilter.AddDate(0, 0, 7)
	//fmt.Println(startDateFilter, " - ", endDateFilter)
	var HolidayStatForColleagues interface{}
	HolidayStatForColleagues, err = V1HolidayStatForColleaguesGeneral(workerID, startDateFilter, endDateFilter)
	if err != nil {
		HolidayStatForColleagues = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
	}

	var GetBranchTree interface{}
	GetBranchTree, err = GetBranchTreeGeneral(workerID)
	if err != nil {
		GetBranchTree = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
	}

	UseYearFilter = false

	var AverageSalary interface{}
	AverageSalary, err = V1AverageSalaryGeneral(workerID, UseYearFilter, yearFilter)
	if err != nil {
		AverageSalary = AnswerWebV1{false, nil, &ErrorWebV1{http.StatusInternalServerError, err.Error()}}
	}

	UseYearFilter = true

	AllInformationV1Answer.HolidayStat = HolidayStat
	AllInformationV1Answer.BudgetStat = BudgetStat
	AllInformationV1Answer.JobPlaces = JobPlaces
	AllInformationV1Answer.HolidayStatForColleagues = HolidayStatForColleagues
	AllInformationV1Answer.GetBranchTree = GetBranchTree
	AllInformationV1Answer.AverageSalary = AverageSalary

	var AnswerWebV1 AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = AllInformationV1Answer
	AnswerWebV1.Error = nil
	//c.JSON(http.StatusOK, AnswerWebV1)

	return AnswerWebV1, nil
}

func V1HolidayStatGeneral(WorkerID string, UseYearFilter bool, yearFilterFrom, yearFilterTo string, RedisClient *libraryGoRedis.Client) (interface{}, error) {

	JSONString, err := redis.GetLibraryGoRedis(RedisClient, WorkerID+yearFilterFrom+yearFilterTo, 3)
	//if err != nil {
	if JSONString == "" {
		//log.Impl.Error(err.Error())
		// JSONString, err = store.GetSettingsByIdJobPg(JobIdParam)
		// if err != nil {
		// 	log.Impl.Error(err)
		// }

		// AnswerWebV1 := AnswerWebV1{false, store.DataAuthorizatioAnswer{}, ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		// c.JSON(http.StatusBadRequest, AnswerWebV1)
	} else {
		//c.Data(http.StatusOK, "application/json", []byte(JSONString))
		//c.JSON(http.StatusOK, JSONString)
		//return []byte(JSONString), nil

		var AnswerWebV1 AnswerWebV1
		if err := json.Unmarshal([]byte(JSONString), &AnswerWebV1); err != nil {
			return nil, err
		}

		return AnswerWebV1, nil
	}

	DB, err := store.GetDB(config.Conf.DatabaseURLMainAnalytics)
	if err != nil {
		return nil, err
	}

	var argsquery []interface{}
	argsquery = append(argsquery, WorkerID)
	//queryAllColumns := "SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = $1;"

	queryAllColumns := `select
		period,
		valid_until,
		status,
		document_base,
		year
	from
		lkr_otsutstviy_all
	where
		collaborator_id = $1
		and status <> 'Работа'
	union all
	select
		case
			when otpuska.moved = 'Да' then otpuska.moved_data_start
			else otpuska.data_start
		end as data_start,
		case
			when otpuska.moved = 'Да' then otpuska.moved_data_end
			else otpuska.data_end
		end as data_start,
		'Отпуск по графику',
		case
			when moving_doc is null then moving_doc
			else planing_doc
		end,
		replace(year, ' ', '') as year
	from
		lkr_otpuska as otpuska
	where
		collaborator_id = $1`

	rows, err := DB.Query(queryAllColumns, argsquery...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mapCheckDoubles := make(map[string]bool)

	ColumnsStructSlice := []V1HolidayAllStatResponds{}
	for rows.Next() {
		var r V1HolidayAllStatResponds
		var year string
		err = rows.Scan(&r.DateStart, &r.DateEnd, &r.Status, &r.DocumentBase, &year)
		if err != nil {
			return nil, err
		}

		_, ok := mapCheckDoubles[r.DateStart+r.DateEnd]
		if ok {
			continue
		}
		mapCheckDoubles[r.DateStart+r.DateEnd] = true

		if UseYearFilter {

			YearAccruals := 0

			if r.Status == "Отпуск по графику" {

				year = strings.Replace(year, " ", "", -1)
				YearAccruals, err = strconv.Atoi(year)
				if err != nil {
					return nil, err
				}

			} else {
				re := regexp.MustCompile(`\d{2}.\d{2}.\d{4}`)
				date_from_subjectArray := re.FindAllString(year, -1)
				if len(date_from_subjectArray) == 0 {
					log.Impl.Errorf("Ошибка парсинга даты для метода: %s collaborator_id: %s\n", "V1HolidayStatGeneral", WorkerID)
					continue
				}
				//fmt.Printf("%q\n", date_from_subjectArray)

				date_from_subject, err := time.Parse("02.01.2006", date_from_subjectArray[0])
				if err != nil {
					return nil, err
				}
				//fmt.Println(year)

				//yearArg, monthArg, dayArg := time.Now().Date()
				YearAccruals = date_from_subject.Year()
			}

			yearFilterFromInt, err := strconv.Atoi(yearFilterFrom)
			if err != nil {
				return nil, err
			}

			yearFilterToInt, err := strconv.Atoi(yearFilterTo)
			if err != nil {
				return nil, err
			}

			if yearFilterFromInt > int(YearAccruals) || int(YearAccruals) > yearFilterToInt {
				continue
			}

			// Ok := strings.Contains(year, yearFilter)
			// if !Ok {
			// 	continue
			// }
		}

		ColumnsStructSlice = append(ColumnsStructSlice, r)
	}

	var AnswerWebV1 AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = ColumnsStructSlice
	AnswerWebV1.Error = nil
	//c.JSON(http.StatusOK, AnswerWebV1)

	// byteResult, err := json.Marshal(AnswerWebV1)
	// if err != nil {
	// 	return nil, err
	// }
	return AnswerWebV1, nil

}

func V1BudgetStatGeneral(WorkerID string, UseYearFilter bool, yearFilter string, RedisClient *libraryGoRedis.Client) (interface{}, error) {

	JSONString, err := redis.GetLibraryGoRedis(RedisClient, WorkerID+yearFilter, 2)
	//if err != nil {
	if JSONString == "" {
		//log.Impl.Error(err.Error())
		// JSONString, err = store.GetSettingsByIdJobPg(JobIdParam)
		// if err != nil {
		// 	log.Impl.Error(err)
		// }

		// AnswerWebV1 := AnswerWebV1{false, store.DataAuthorizatioAnswer{}, ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		// c.JSON(http.StatusBadRequest, AnswerWebV1)
	} else {
		// c.Data(http.StatusOK, "application/json", []byte(JSONString))
		// //c.JSON(http.StatusOK, JSONString)
		// return
		var AnswerWebV1 AnswerWebV1
		if err := json.Unmarshal([]byte(JSONString), &AnswerWebV1); err != nil {
			return nil, err
		}

		return AnswerWebV1, nil
	}

	DB, err := store.GetDB(config.Conf.DatabaseURLMainAnalytics)
	if err != nil {
		return nil, err
	}

	var argsquery []interface{}
	argsquery = append(argsquery, WorkerID)
	//queryAllColumns := "SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = $1;"

	// queryAllColumns := `select
	// 	date_registration,
	// 	settlement_group,
	// 	calculation_type,
	// 	days_worked,
	// 	hours_worked,
	// 	replace(summa, ' ', '')
	// from
	// 	lkr_nachisleniy_zp
	// where
	// 	collaborator_id = $1
	// union all
	// select
	// 	date_registration,
	// 	settlement_group,
	// 	calculation_type,
	// 	days_worked,
	// 	hours_worked,
	// 	replace(summa, ' ', '')
	// from
	// 	lkr_nachisleniy_zp2022
	// where
	// 	collaborator_id = $1
	// order by
	// 	1`

	queryAllColumns := `select
		date_registration,
		settlement_group,
		calculation_type,
		days_worked,
		hours_worked,
		replace(summa, ' ', '')
	from
		lkr_nachisleniy_zp2020
	where
		collaborator_id = $1
	union all
	select
		date_registration,
		settlement_group,
		calculation_type,
		days_worked,
		hours_worked,
		replace(summa, ' ', '')
	from
		lkr_nachisleniy_zp2021
	where
		collaborator_id = $1
	union all
	select
		date_registration,
		settlement_group,
		calculation_type,
		days_worked,
		hours_worked,
		replace(summa, ' ', '')
	from
		lkr_nachisleniy_zp2022
	where
		collaborator_id = $1
	order by
		1`

	rows, err := DB.Query(queryAllColumns, argsquery...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//ColumnsStructSlice := []store.V1BudgetStatResponds{}
	//V1BudgetStatGroupResponds := store.V1BudgetStatGroupResponds{}

	MapV1BudgetStatGroupResponds := make(map[int]V1BudgetStatGroupResponds)

	// Создать карту  с V1BudgetStatGroupResponds по месяцу.
	// Каждый раз доставать эту карту и обогащать ее данные + добавлять новые значения в массив
	// провести сортировку записей по дате + провести сортироку карты через пакет сорт или через слайз, после сортировки
	// Сделать конечный массив с группировками и выести его в итоги.

	for rows.Next() {
		var r V1BudgetStatResponds

		var DaysWorked, HoursWorked, Summ string
		err = rows.Scan(&r.DateRegistration, &r.SettlementGroup, &r.CalculationType, &DaysWorked, &HoursWorked, &Summ)
		if err != nil {
			return nil, err
		}

		if Summ == "0" && DaysWorked == "0" && HoursWorked == "0" {
			continue
		}

		re := regexp.MustCompile(`\d{2}.\d{2}.\d{4}`)
		date_from_subjectArray := re.FindAllString(r.DateRegistration, -1)
		if len(date_from_subjectArray) == 0 {
			log.Impl.Errorf("Ошибка парсинга даты для метода: %s collaborator_id: %s\n", "V1BudgetStatGeneral", WorkerID)
			continue
		}

		//fmt.Printf("%q\n", date_from_subjectArray)

		date_from_subject, err := time.Parse("02.01.2006", date_from_subjectArray[0])
		if err != nil {
			return nil, err
		}
		//fmt.Println(date_from_subject)

		//yearArg, monthArg, dayArg := time.Now().Date()
		MonthAccruals := date_from_subject.Month()
		//fmt.Println(int(month))
		//V1BudgetStatGroupResponds.Month = int(MonthAccruals)

		var V1BudgetStatGroupResponds V1BudgetStatGroupResponds
		V1BudgetStatGroupResponds, _ = MapV1BudgetStatGroupResponds[int(MonthAccruals)]
		// if !ok {
		// 	var V1BudgetStatGroupResponds store.V1BudgetStatGroupResponds
		// 	MapV1BudgetStatGroupResponds[int(MonthAccruals)] = V1BudgetStatGroupResponds
		// }

		V1BudgetStatGroupResponds.Month = int(MonthAccruals)

		Summ = strings.Replace(Summ, ",", ".", -1)
		Summ = strings.Replace(Summ, " ", "", -1)

		//HoursWorked = strings.Replace(HoursWorked, "\n", "", -1)
		SummFloat, err := strconv.ParseFloat(Summ, 32)
		if err != nil {
			return nil, err
		}

		r.Summ = float32(SummFloat)

		if r.SettlementGroup == "Начислено" {
			//r.Summ = 10
		} else {
			//r.Summ = -10
			r.Summ = -r.Summ
		}

		V1BudgetStatGroupResponds.Total = V1BudgetStatGroupResponds.Total + r.Summ
		if r.SettlementGroup == "Начислено" {
			V1BudgetStatGroupResponds.TotalGross = V1BudgetStatGroupResponds.TotalGross + r.Summ
		} else {
			V1BudgetStatGroupResponds.TotalDeduction = V1BudgetStatGroupResponds.TotalDeduction + -r.Summ
		}

		if UseYearFilter {
			// TODO: Подключить по возможности к регулярному выражению ниже.
			Ok := strings.Contains(r.DateRegistration, yearFilter)
			if !Ok {
				continue
			}
		}

		DaysWorkedInt, err := strconv.Atoi(DaysWorked)
		if err != nil {
			return nil, err
		}
		V1BudgetStatGroupResponds.DaysWorked = DaysWorkedInt

		//fmt.Println(HoursWorked)
		//fmt.Println(777)
		//HoursWorked = strings.Replace(HoursWorked, `\`, "", -1)
		//HoursWorked = strings.ReplaceAll(HoursWorked, string([]byte{92, 114, 92, 110}), "")
		//re := regexp.MustCompile(`\r?\n`)
		//HoursWorked = re.ReplaceAllString(HoursWorked, "")
		//HoursWorked = strings.TrimSpace(HoursWorked)

		HoursWorked = strings.Replace(HoursWorked, ",", ".", -1)

		//HoursWorked = strings.Replace(HoursWorked, "\n", "", -1)
		HoursWorkedFloat, err := strconv.ParseFloat(HoursWorked, 32)
		if err != nil {
			return nil, err
		}

		V1BudgetStatGroupResponds.HoursWorked = float32(HoursWorkedFloat)
		//ColumnsStructSlice = append(ColumnsStructSlice, r)
		V1BudgetStatGroupResponds.Items = append(V1BudgetStatGroupResponds.Items, r)

		MapV1BudgetStatGroupResponds[int(MonthAccruals)] = V1BudgetStatGroupResponds
	}

	var V1BudgetStatGroupRespondsSlice []V1BudgetStatGroupResponds

	keys := make([]int, 0, len(MapV1BudgetStatGroupResponds))
	for k := range MapV1BudgetStatGroupResponds {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		V1BudgetStatGroupRespondsSlice = append(V1BudgetStatGroupRespondsSlice, MapV1BudgetStatGroupResponds[k])
	}

	var AnswerWebV1 AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = V1BudgetStatGroupRespondsSlice
	AnswerWebV1.Error = nil
	//c.JSON(http.StatusOK, AnswerWebV1)

	return AnswerWebV1, nil

}

func V2JobPlacesGeneral(WorkerID string, RedisClient *libraryGoRedis.Client) (interface{}, error) {

	DB, err := store.GetDB(config.Conf.DatabaseURLMainAnalytics)
	if err != nil {
		return nil, err
	}

	var argsquery []interface{}
	argsquery = append(argsquery, WorkerID)

	queryText := `select
		collaborators_posle.person_id as person_id,
		collaborators_posle.collaborator_id as collaborator_id,
		collaborators_posle.insurance_number as insurance_number,
		collaborators_posle.inn as inn,
		collaborators_posle.full_name as full_name,
		collaborators_posle.position as position,
		organizations_zkgu.name as organization_name,
		collaborators_posle.status as status,
		coalesce(contact_inf_pochta_posle.email, '') as email,
		coalesce(contact_inf_telephone_posle.mobile, '') as mobile_phone,
		coalesce(contact_inf_telephone_posle."work", '') as work_phone,
		collaborators_posle.date_birth as date_birth,
		collaborators_posle.podrazdelenie as podrazdelenie,
		coalesce(collaborators_posle.podrazdelenie_id, '') as podrazdelenie_id,
		coalesce(dit_gruppirovka_dolzhnostey.large_group_of_posts, '') as large_group_of_posts,
		coalesce(dit_gruppirovka_dolzhnostey.position_tag, '') as position_tag,
		coalesce(collaborators_posle.updated_at, DATE '0001-01-01') as updated_at,
		coalesce(collaborators_posle.date_dismissals_as_date, DATE '0001-01-01') as date_dismissals_as_date
	from
		collaborators_posle as collaborators_posle
	left join dit_gruppirovka_dolzhnostey as dit_gruppirovka_dolzhnostey on
		collaborators_posle.position = dit_gruppirovka_dolzhnostey.position
	left join organizations_zkgu as organizations_zkgu on
		collaborators_posle.organization_id = organizations_zkgu.organization_id
	left join contact_inf_pochta_posle as contact_inf_pochta_posle on
		collaborators_posle.person_id = contact_inf_pochta_posle.person_id
	left join contact_inf_telephone_posle as contact_inf_telephone_posle on
		collaborators_posle.person_id = contact_inf_telephone_posle.person_id
	where
		collaborator_id = $1`

	rows, err := DB.Query(queryText, argsquery...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ColumnsStructSlice := []V1ActiveWorkers{}
	for rows.Next() {
		var r V1ActiveWorkers
		err = rows.Scan(&r.PersonId, &r.CollaboratorId, &r.InsuranceNumber, &r.Inn, &r.FullName, &r.Position, &r.OrganizationName, &r.Status,
			&r.Email, &r.MobilePhone, &r.WorkPhone, &r.DateBirth, &r.BranchName, &r.BranchID, &r.LargeGroupOfPosts, &r.Position_tag, &r.UpdatedAt, &r.DateDismissals)
		if err != nil {
			return nil, err
		}

		// JSONString, err := redis.GetLibraryGoRedis(RedisClient, r.InsuranceNumber, 1)
		// if err != nil {
		// 	return err, nil
		// }

		// r.EmailArray = JSONString

		ColumnsStructSlice = append(ColumnsStructSlice, r)
	}

	var AnswerWebV1 AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = ColumnsStructSlice
	AnswerWebV1.Error = nil
	//c.JSON(http.StatusOK, AnswerWebV1)

	return AnswerWebV1, nil

}

func V3JobPlacesGeneral(WorkerID string, RedisClient *libraryGoRedis.Client) (interface{}, error) {

	JSONString, err := redis.GetLibraryGoRedis(RedisClient, WorkerID, 5)
	//if err != nil {
	if JSONString == "" {
		//log.Impl.Error(err.Error())
		// JSONString, err = store.GetSettingsByIdJobPg(JobIdParam)
		// if err != nil {
		// 	log.Impl.Error(err)
		// }

		// AnswerWebV1 := AnswerWebV1{false, store.DataAuthorizatioAnswer{}, ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		// c.JSON(http.StatusBadRequest, AnswerWebV1)
	} else {
		// c.Data(http.StatusOK, "application/json", []byte(JSONString))
		// //c.JSON(http.StatusOK, JSONString)
		// return
		var AnswerWebV1 AnswerWebV1
		if err := json.Unmarshal([]byte(JSONString), &AnswerWebV1); err != nil {
			return nil, err
		}

		return AnswerWebV1, nil
	}

	DB, err := store.GetDB(config.Conf.DatabaseURLMainAnalytics)
	if err != nil {
		return nil, err
	}

	var argsquery []interface{}
	argsquery = append(argsquery, WorkerID)

	//coalesce(collaborators_posle.updated_at, DATE '1900-01-01') as updated_at,
	//coalesce(collaborators_posle.created_at, DATE '1900-01-01') as created_at,

	queryText := `select
		collaborators_posle.person_id as person_id,
		collaborators_posle.collaborator_id as collaborator_id,
		collaborators_posle.insurance_number as insurance_number,
		collaborators_posle.inn as inn,
		collaborators_posle.full_name as full_name,
		collaborators_posle.position as position,
		organizations_zkgu.name as organization_name,
		collaborators_posle.status as status,
		coalesce(contact_inf_pochta_posle.email, '') as email,
		coalesce(contact_inf_pochta_posle."emailEPS", '') as emailEPS,
		coalesce(contact_inf_telephone_posle.mobile, '') as mobile_phone,
		coalesce(contact_inf_telephone_posle."work", '') as work_phone,
		collaborators_posle.date_birth as date_birth,
		collaborators_posle.podrazdelenie as podrazdelenie,
		coalesce(collaborators_posle.podrazdelenie_id, '') as podrazdelenie_id,
		coalesce(dit_gruppirovka_dolzhnostey.large_group_of_posts, '') as large_group_of_posts,
		coalesce(dit_gruppirovka_dolzhnostey.position_tag, '') as position_tag,
		case
			when (collaborators_posle.created_at >contact_inf_pochta_posle.created_at
			and collaborators_posle.created_at >contact_inf_telephone_posle.created_at) then coalesce(collaborators_posle.created_at, DATE '1900-01-01')
			when (contact_inf_pochta_posle.created_at>collaborators_posle.created_at
			and contact_inf_pochta_posle.created_at>contact_inf_telephone_posle.created_at) then coalesce(contact_inf_pochta_posle.created_at, DATE '1900-01-01')
			else coalesce(contact_inf_telephone_posle.created_at, DATE '1900-01-01')
		end as created_at,
		case
			when (collaborators_posle.updated_at>contact_inf_pochta_posle.updated_at
			and collaborators_posle.updated_at>contact_inf_telephone_posle.created_at) then coalesce(collaborators_posle.updated_at, DATE '1900-01-01')
			when (contact_inf_pochta_posle.created_at>collaborators_posle.updated_at
			and contact_inf_pochta_posle.created_at>contact_inf_telephone_posle.created_at) then coalesce(contact_inf_pochta_posle.created_at, DATE '1900-01-01')
			else coalesce(contact_inf_telephone_posle.updated_at, DATE '1900-01-01')
		end as updated_at,
		coalesce(collaborators_posle.date_dismissals_as_date, DATE '0001-01-01') as date_dismissals_as_date
	from
		collaborators_posle as collaborators_posle
	left join dit_gruppirovka_dolzhnostey as dit_gruppirovka_dolzhnostey on
		collaborators_posle.position = dit_gruppirovka_dolzhnostey.position
	left join organizations_zkgu as organizations_zkgu on
		collaborators_posle.organization_id = organizations_zkgu.organization_id
		and collaborators_posle.area = organizations_zkgu.area
	left join contact_inf_pochta_posle as contact_inf_pochta_posle on
		collaborators_posle.person_id = contact_inf_pochta_posle.person_id
	left join contact_inf_telephone_posle as contact_inf_telephone_posle on
		collaborators_posle.person_id = contact_inf_telephone_posle.person_id
	where
		collaborator_id = $1`

	rows, err := DB.Query(queryText, argsquery...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//ColumnsStructSlice := []V1ActiveWorkers{}
	ColumnsStruct := V1ActiveWorkers{}
	for rows.Next() {
		var r V1ActiveWorkers
		err = rows.Scan(&r.PersonId, &r.CollaboratorId, &r.InsuranceNumber, &r.Inn, &r.FullName, &r.Position, &r.OrganizationName, &r.Status,
			&r.Email, &r.EmailEPS, &r.MobilePhone, &r.WorkPhone, &r.DateBirth, &r.BranchName, &r.BranchID, &r.LargeGroupOfPosts, &r.Position_tag, &r.CreatedAt, &r.UpdatedAt, &r.DateDismissals)
		if err != nil {
			return nil, err
		}

		//JSONString, err := redis.GetLibraryGoRedis(RedisClient, r.InsuranceNumber, 1)
		//if err != nil {
		//	return err, nil
		//}

		//r.EmailArray = JSONString
		//r.EmailEPS = ""

		//ColumnsStructSlice = append(ColumnsStructSlice, r)

		ColumnsStruct = r
	}

	var AnswerWebV1 AnswerWebV1
	AnswerWebV1.Status = true
	if ColumnsStruct.CollaboratorId == "" {
		//AnswerWebV1.Data = nil
	} else {
		AnswerWebV1.Data = ColumnsStruct
	}
	AnswerWebV1.Error = nil
	//c.JSON(http.StatusOK, AnswerWebV1)

	return AnswerWebV1, nil

}

func V1JobPlacesGeneral(WorkerID string, RedisClient *libraryGoRedis.Client) (interface{}, error) {

	DB, err := store.GetDB(config.Conf.DatabaseURLMainAnalytics)
	if err != nil {
		return nil, err
	}

	var argsquery []interface{}
	argsquery = append(argsquery, WorkerID)
	//queryAllColumns := "SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = $1;"

	queryAllColumns := `select
		lkr_kadrovie_dannie.person_id,
		lkr_kadrovie_dannie.collaborator_id,
		lkr_kadrovie_dannie.insurance_number,
		lkr_kadrovie_dannie.inn,
		lkr_kadrovie_dannie.full_name,
		lkr_kadrovie_dannie.position,
		lkr_kadrovie_dannie.organization_name,
		lkr_kadrovie_dannie.status,
		lkr_kadrovie_dannie.email,
		lkr_kadrovie_dannie.mobile_phone,
		lkr_kadrovie_dannie.work_phone,
		lkr_kadrovie_dannie.date_birth,
		lkr_kadrovie_dannie.podrazdelenie,
		coalesce(lkr_kadrovie_dannie.guid_podrazdelenie, ''),
		coalesce(dit_gruppirovka_dolzhnostey.large_group_of_posts, '') as large_group_of_posts,
		coalesce(dit_gruppirovka_dolzhnostey.position_tag, '') as position_tag,
		COALESCE(lkr_kadrovie_dannie.updated_at, DATE '0001-01-01')
	from
		lkr_kadrovie_dannie as lkr_kadrovie_dannie
	left join dit_gruppirovka_dolzhnostey as dit_gruppirovka_dolzhnostey on
		lkr_kadrovie_dannie.position = dit_gruppirovka_dolzhnostey.position
	where
		collaborator_id = $1`
	//where
	//	insurance_number = $1`

	rows, err := DB.Query(queryAllColumns, argsquery...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ColumnsStructSlice := []V1ActiveWorkers{}
	for rows.Next() {
		var r V1ActiveWorkers
		err = rows.Scan(&r.PersonId, &r.CollaboratorId, &r.InsuranceNumber, &r.Inn, &r.FullName, &r.Position, &r.OrganizationName, &r.Status,
			&r.Email, &r.MobilePhone, &r.WorkPhone, &r.DateBirth, &r.BranchName, &r.BranchID, &r.LargeGroupOfPosts, &r.Position_tag, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// JSONString, err := redis.GetLibraryGoRedis(RedisClient, r.InsuranceNumber, 1)
		// if err != nil {
		// 	return err, nil
		// }

		// r.EmailArray = JSONString

		ColumnsStructSlice = append(ColumnsStructSlice, r)
	}

	var AnswerWebV1 AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = ColumnsStructSlice
	AnswerWebV1.Error = nil
	//c.JSON(http.StatusOK, AnswerWebV1)

	return AnswerWebV1, nil

}

func V1HolidayStatForColleaguesGeneral(WorkerID string, startDateFilter time.Time, endDateFilter time.Time) (interface{}, error) {

	// Добрый день,
	// Для блока "ближайшие отсутствия коллег" нужен метод со следующими входящими параметрами:
	// from - дата от (формат dd.MM.yyyy)
	// to - дата до (формат dd.MM.yyyy)
	// branch - необязательный, id подразделения, если не указан, отдавать всё за указанный промежуток времени, если указан, то за указнный промежуток времени для подразделения

	//Так же нужно добавить еще два обязательных параметра: ИНН и КПП организации, чтобы можно было отфильтровать отсутствия по организации без подразделения

	//Возможно к возврату нужно добавить еще и ФИО?

	// Получить родителя подразделение
	// 	select
	// 	table_podrazdelenie_id.roditel,
	// 	table_podrazdelenie_id.roditel_guid
	// from
	// 	(
	// 	select
	// 		lkr_podrazdelenie_branch.area,
	// 		lkr_podrazdelenie_branch.unit_guid,
	// 		lkr_podrazdelenie_branch.unit_name,
	// 		lkr_podrazdelenie_branch.roditel,
	// 		lkr_podrazdelenie_branch.roditel_guid,
	// 		collaborators_posle.podrazdelenie_id
	// 	from
	// 		public.lkr_podrazdelenie_branch as lkr_podrazdelenie_branch
	// 	inner join collaborators_posle as collaborators_posle on
	// 		lkr_podrazdelenie_branch.area = collaborators_posle.area
	// 		and collaborators_posle.collaborator_id = '8f4d1e85-fde4-11eb-9113-005056a2fd67'
	// 	where
	// 		lkr_podrazdelenie_branch.unit_guid = collaborators_posle.podrazdelenie_id) as table_podrazdelenie_id

	// Получить просто подразделение
	// select
	// 	collaborators_posle.collaborator_id
	// from collaborators_posle
	// where collaborators_posle.podrazdelenie_id in
	// 	(select
	// 		collaborators_posle.podrazdelenie_id
	// 	from
	// 		collaborators_posle where  collaborators_posle.collaborator_id = '8f4d1e85-fde4-11eb-9113-005056a2fd67')

	//?from=2020&to=2023

	//UseYearFilter := false
	// yearFilter := c.Query("year")
	// if yearFilter != "" {
	// 	UseYearFilter = true
	// }

	// yearFilterFrom := c.Query("from")
	// yearFilterTo := c.Query("to")
	// if yearFilterFrom != "" && yearFilterTo != "" {
	// 	UseYearFilter = true
	// }

	DB, err := store.GetDB(config.Conf.DatabaseURLMainAnalytics)
	if err != nil {
		return err, nil
	}

	var argsquery []interface{}
	argsquery = append(argsquery, WorkerID)

	// Выбираем коллег по текущему подразделению, в которм находится сотрудник
	// 	queryAllColumns := `select
	// 	lkr_otsutstviy_all.collaborator_id,
	// 	collaborators_posle.full_name,
	// 	lkr_otsutstviy_all.period,
	// 	lkr_otsutstviy_all.valid_until,
	// 	lkr_otsutstviy_all.status,
	// 	lkr_otsutstviy_all.document_base,
	// 	lkr_otsutstviy_all.year
	// from
	// 	lkr_otsutstviy_all
	// left join collaborators_posle as collaborators_posle on
	// 	lkr_otsutstviy_all.collaborator_id = collaborators_posle.collaborator_id
	// where
	// 	lkr_otsutstviy_all.collaborator_id in (
	// 	select
	// 		collaborators_posle.collaborator_id
	// 	from
	// 		collaborators_posle
	// 	where
	// 		collaborators_posle.podrazdelenie_id in (
	// 		select
	// 			collaborators_posle.podrazdelenie_id
	// 		from
	// 			collaborators_posle
	// 		where
	// 			collaborators_posle.collaborator_id = $1))
	// 	and lkr_otsutstviy_all.status <> 'Работа'
	// 	and lkr_otsutstviy_all.collaborator_id <> $1
	// union all
	// select
	// 	otpuska.collaborator_id,
	// 	collaborators_posle.full_name,
	// 	case
	// 		when otpuska.moved = 'Да' then otpuska.moved_data_start
	// 		else otpuska.data_start
	// 	end as data_start,
	// 	case
	// 		when otpuska.moved = 'Да' then otpuska.moved_data_end
	// 		else otpuska.data_end
	// 	end as data_end,
	// 	'Отпуск по графику',
	// 	case
	// 		when moving_doc is null then moving_doc
	// 		else planing_doc
	// 	end,
	// 	replace(otpuska.year, ' ', '') as year
	// from
	// 	otpuska as otpuska
	// left join collaborators_posle as collaborators_posle on
	// 	otpuska.collaborator_id = collaborators_posle.collaborator_id
	// where
	// 	otpuska.collaborator_id in (
	// 	select
	// 		collaborators_posle.collaborator_id
	// 	from
	// 		collaborators_posle
	// 	where
	// 		collaborators_posle.podrazdelenie_id in (
	// 		select
	// 			collaborators_posle.podrazdelenie_id
	// 		from
	// 			collaborators_posle
	// 		where
	// 			collaborators_posle.collaborator_id = $1))
	// 		and otpuska.collaborator_id <> $1
	// order by 2`

	// Выбираем коллег их одной организации.
	queryAllColumns := `select
	lkr_otsutstviy_all.collaborator_id,
	collaborators_posle.full_name,
	lkr_otsutstviy_all.period,
	lkr_otsutstviy_all.valid_until,
	lkr_otsutstviy_all.status,
	lkr_otsutstviy_all.document_base,
	lkr_otsutstviy_all.year
from
	lkr_otsutstviy_all
left join collaborators_posle as collaborators_posle on
	lkr_otsutstviy_all.collaborator_id = collaborators_posle.collaborator_id
where
	lkr_otsutstviy_all.collaborator_id in (
		select
			collaborators_posle.collaborator_id
		from
			collaborators_posle
		where
			collaborators_posle.area in (
			select
				collaborators_posle.area
			from
				collaborators_posle
			where
				collaborators_posle.collaborator_id = $1))
	and lkr_otsutstviy_all.status <> 'Работа'
	and collaborators_posle.status <> 'Увольнение'
union all
select
	otpuska.collaborator_id,
	collaborators_posle.full_name,
	case
		when otpuska.moved = 'Да' then otpuska.moved_data_start
		else otpuska.data_start
	end as data_start,
	case
		when otpuska.moved = 'Да' then otpuska.moved_data_end
		else otpuska.data_end
	end as data_end,
	'Отпуск по графику',
	case
		when moving_doc is null then moving_doc
		else planing_doc
	end,
	replace(otpuska.year, ' ', '') as year
from
	lkr_otpuska as otpuska
left join collaborators_posle as collaborators_posle on
	otpuska.collaborator_id = collaborators_posle.collaborator_id
where
	otpuska.collaborator_id in (
		select
			collaborators_posle.collaborator_id
		from
			collaborators_posle
		where
			collaborators_posle.area in (
			select
				collaborators_posle.area
			from
				collaborators_posle
			where
				collaborators_posle.collaborator_id = $1))
		and collaborators_posle.status <> 'Увольнение'
order by 2`

	rows, err := DB.Query(queryAllColumns, argsquery...)
	if err != nil {
		return err, nil
	}

	//currentTime := time.Now()

	defer rows.Close()

	mapCheckDoubles := make(map[string]bool)

	ColumnsStructSlice := []V1HolidayAllStatRespondsForColleagues{}
	for rows.Next() {
		var r V1HolidayAllStatRespondsForColleagues
		var year string
		err = rows.Scan(&r.CollaboratorId, &r.FullName, &r.DateStart, &r.DateEnd, &r.Status, &r.DocumentBase, &year)
		if err != nil {
			return nil, err
		}

		_, ok := mapCheckDoubles[r.CollaboratorId+r.DateStart+r.DateEnd]
		if ok {
			continue
		}
		mapCheckDoubles[r.CollaboratorId+r.DateStart+r.DateEnd] = true

		if r.Status == "Увольнение" {
			continue
		}

		re := regexp.MustCompile(`\d{2}.\d{2}.\d{4}`)
		//compareData := false

		DateStartArray := re.FindAllString(r.DateStart, -1)
		if len(DateStartArray) == 0 {
			log.Impl.Errorf("Ошибка парсинга даты для метода: %s collaborator_id: %s\n", "V1HolidayStatForColleaguesGeneral", r.CollaboratorId)
			continue
		}
		DateStart, err := time.Parse("02.01.2006", DateStartArray[0])
		if err != nil {
			return nil, err
		}

		DateEndArray := re.FindAllString(r.DateEnd, -1)
		if len(DateEndArray) == 0 {
			log.Impl.Errorf("Ошибка парсинга даты для метода: %s collaborator_id: %s\n", "V1HolidayStatForColleaguesGeneral", r.CollaboratorId)
			continue
		}
		DateEnd, err := time.Parse("02.01.2006", DateEndArray[0])
		if err != nil {
			return nil, err
		}

		// //yearArg, monthArg, dayArg := time.Now().Date()
		// filer1 := false
		// filer11 := false

		// filer2 := false
		// filer22 := false

		// Если ДатаОтбораНачало > ОтпускНачало и ДатаОтбораНачало < ОтпускКонец = +
		// Если ДатаОтбораКонец > ОтпускНачало и ДатаОтбораКонец < ОтпускКонец = -

		// compareData = startDateFilter.After(DateStart)
		// if !compareData {
		// 	compareDataNested := startDateFilter.Equal(DateStart)
		// 	if compareDataNested {
		// 		filer1 = true
		// 	}
		// } else {
		// 	filer1 = true
		// }

		// compareData = startDateFilter.Before(DateEnd)
		// if !compareData {
		// 	compareDataNested := startDateFilter.Equal(DateEnd)
		// 	if compareDataNested {
		// 		filer11 = true
		// 	}
		// } else {
		// 	filer11 = true
		// }

		// compareData = endDateFilter.After(DateStart)
		// if !compareData {
		// 	compareDataNested := endDateFilter.Equal(DateStart)
		// 	if compareDataNested {
		// 		filer2 = true
		// 	}
		// } else {
		// 	filer2 = true
		// }

		// compareData = endDateFilter.Before(DateEnd)
		// if !compareData {
		// 	compareDataNested := endDateFilter.Equal(DateEnd)
		// 	if compareDataNested {
		// 		filer22 = true
		// 	}
		// } else {
		// 	filer22 = true
		// }

		// Если ДатаОтбораКонец => ОтпускНачало и ОтпускНачало => ДатаОтбораНачало
		// Если ДатаОтбораКонец => ОтпускКонец и ОтпускКонец => ДатаОтбораНачало

		// compareData = endDateFilter.After(DateStart)
		// if !compareData {
		// 	compareDataNested := endDateFilter.Equal(DateStart)
		// 	if compareDataNested {
		// 		filer1 = true
		// 	}
		// } else {
		// 	filer1 = true
		// }

		// compareData = DateStart.After(startDateFilter)
		// if !compareData {
		// 	compareDataNested := DateStart.Equal(startDateFilter)
		// 	if compareDataNested {
		// 		filer11 = true
		// 	}
		// } else {
		// 	filer11 = true
		// }

		// compareData = endDateFilter.After(DateEnd)
		// if !compareData {
		// 	compareDataNested := endDateFilter.Equal(DateEnd)
		// 	if compareDataNested {
		// 		filer2 = true
		// 	}
		// } else {
		// 	filer2 = true
		// }

		// compareData = DateEnd.After(startDateFilter)
		// if !compareData {
		// 	compareDataNested := DateEnd.Equal(startDateFilter)
		// 	if compareDataNested {
		// 		filer22 = true
		// 	}
		// } else {
		// 	filer22 = true
		// }

		// if filer1 && filer11 || filer2 && filer22 {

		// } else {
		// 	continue
		// }

		// if filer1 && filer11 || filer2 && filer22 {

		// } else {
		// 	continue
		// }

		filer1 := utilityShare.InTimeSpan(DateStart, DateEnd, startDateFilter)
		filer2 := utilityShare.InTimeSpan(DateStart, DateEnd, endDateFilter)

		filer3 := utilityShare.InTimeSpan(startDateFilter, endDateFilter, DateStart)
		filer4 := utilityShare.InTimeSpan(startDateFilter, endDateFilter, DateEnd)

		if filer1 || filer2 || filer3 || filer4 {

		} else {
			continue
		}

		// if UseYearFilter {

		// 	YearAccruals := 0

		// 	if r.Status == "Отпуск по графику" {

		// 		year = strings.Replace(year, " ", "", -1)
		// 		YearAccruals, err = strconv.Atoi(year)
		// 		if err != nil {
		// 			AnswerWebV1 := AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		// 			c.JSON(http.StatusBadRequest, AnswerWebV1)
		// 			log.Impl.Error(err.Error())
		// 			return
		// 		}

		// 	} else {
		// 		re := regexp.MustCompile(`\d{2}.\d{2}.\d{4}`)
		// 		date_from_subjectArray := re.FindAllString(year, -1)
		// 		//fmt.Printf("%q\n", date_from_subjectArray)

		// 		date_from_subject, err := time.Parse("02.01.2006", date_from_subjectArray[0])
		// 		if err != nil {
		// 			AnswerWebV1 := AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		// 			c.JSON(http.StatusBadRequest, AnswerWebV1)
		// 			log.Impl.Error(err.Error())
		// 			return
		// 		}
		// 		//fmt.Println(year)

		// 		//yearArg, monthArg, dayArg := time.Now().Date()
		// 		YearAccruals = date_from_subject.Year()
		// 	}

		// 	yearFilterFromInt, err := strconv.Atoi(yearFilterFrom)
		// 	if err != nil {
		// 		AnswerWebV1 := AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		// 		c.JSON(http.StatusBadRequest, AnswerWebV1)
		// 		log.Impl.Error(err.Error())
		// 		return
		// 	}

		// 	yearFilterToInt, err := strconv.Atoi(yearFilterTo)
		// 	if err != nil {
		// 		AnswerWebV1 := AnswerWebV1{false, nil, &web.ErrorWebV1{http.StatusInternalServerError, err.Error()}}
		// 		c.JSON(http.StatusBadRequest, AnswerWebV1)
		// 		log.Impl.Error(err.Error())
		// 		return
		// 	}

		// 	if yearFilterFromInt > int(YearAccruals) || int(YearAccruals) > yearFilterToInt {
		// 		continue
		// 	}

		// 	// Ok := strings.Contains(year, yearFilter)
		// 	// if !Ok {
		// 	// 	continue
		// 	// }
		// }

		ColumnsStructSlice = append(ColumnsStructSlice, r)
	}

	var AnswerWebV1 AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = ColumnsStructSlice
	AnswerWebV1.Error = nil
	//c.JSON(http.StatusOK, AnswerWebV1)

	return AnswerWebV1, nil

}

func GetBranchTreeGeneral(WorkerID string) (interface{}, error) {

	DB, err := store.GetDB(config.Conf.DatabaseURLMainAnalytics)
	if err != nil {
		return nil, err
	}

	var argsquery []interface{}
	argsquery = append(argsquery, WorkerID)

	// queryAllColumns := `SELECT area, unit_guid, unit_name, roditel, roditel_guid
	// FROM public.lkr_podrazdelenie_branch where area = '6083';`

	//lkr_kadrovie_dannie.collaborator_id
	queryAllColumns := `select
	lkr_podrazdelenie_branch.area,
	lkr_podrazdelenie_branch.unit_guid,
	lkr_podrazdelenie_branch.unit_name,
	lkr_podrazdelenie_branch.roditel,
	lkr_podrazdelenie_branch.roditel_guid,
	collaborators_posle.podrazdelenie_id
from
	public.lkr_podrazdelenie_branch as lkr_podrazdelenie_branch
inner join collaborators_posle as collaborators_posle on
	lkr_podrazdelenie_branch.area = collaborators_posle.area
	and collaborators_posle.collaborator_id = $1
order by
	lkr_podrazdelenie_branch.roditel, lkr_podrazdelenie_branch.unit_name`
	// TODO: Смена сортировки ведет к потери узлов в дереве

	rows, err := DB.Query(queryAllColumns, argsquery...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ColumnsStructSlice := []tree.BranchTree{}
	for rows.Next() {
		var r tree.BranchTree
		err = rows.Scan(&r.Area, &r.BranchID, &r.BranchName, &r.PatentName, &r.PatentID, &r.CurrectBranchId)
		if err != nil {
			return nil, err
		}
		ColumnsStructSlice = append(ColumnsStructSlice, r)
	}

	Note := tree.AssembleTreeHandler(ColumnsStructSlice)

	var AnswerWebV1 AnswerWebV1
	AnswerWebV1.Status = true
	AnswerWebV1.Data = Note
	AnswerWebV1.Error = nil

	// byteTest, _ := json.Marshal(&Note)
	// fmt.Println(string(byteTest))

	//c.JSON(http.StatusOK, AnswerWebV1)

	return AnswerWebV1, nil

}

func V1AverageSalaryGeneral(WorkerID string, UseYearFilter bool, yearFilter string) (interface{}, error) {

	// Начисленные суммировать за месяцы и разделить на 29,3
	// Получать количество отработанных дней за период.

	//UseYearFilter = false
	// yearFilter = "2021"

	currentTime := time.Now()
	//fmt.Println("Today:", currentTime)

	subtractYear := currentTime.AddDate(-1, 0, 0)
	//	fmt.Println("Subtract 1 Year:", subtractYear)

	// JSONString, err := GetDataRedisByInsuranceNumber(InsuranceNumber+yearFilter, 2)
	// //if err != nil {
	// if JSONString == "" {
	// 	//log.Impl.Error(err.Error())
	// 	// JSONString, err = store.GetSettingsByIdJobPg(JobIdParam)
	// 	// if err != nil {
	// 	// 	log.Impl.Error(err)
	// 	// }

	// 	// AnswerWebV1 := AnswerWebV1{false, store.DataAuthorizatioAnswer{}, ErrorWebV1{http.StatusInternalServerError, err.Error()}}
	// 	// c.JSON(http.StatusBadRequest, AnswerWebV1)
	// } else {
	// 	c.Data(http.StatusOK, "application/json", []byte(JSONString))
	// 	//c.JSON(http.StatusOK, JSONString)
	// 	return
	// }

	DB, err := store.GetDB(config.Conf.DatabaseURLMainAnalytics)
	if err != nil {
		return nil, err
	}

	var argsquery []interface{}
	argsquery = append(argsquery, WorkerID)
	//queryAllColumns := "SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = $1;"

	// queryAllColumns := `select
	// 	date_registration,
	// 	settlement_group,
	// 	calculation_type,
	// 	days_worked,
	// 	hours_worked,
	// 	replace(summa, ' ', '')
	// from
	// 	lkr_nachisleniy_zp
	// where
	// 	collaborator_id = $1
	// union all
	// select
	// 	date_registration,
	// 	settlement_group,
	// 	calculation_type,
	// 	days_worked,
	// 	hours_worked,
	// 	replace(summa, ' ', '')
	// from
	// 	lkr_nachisleniy_zp2022
	// where
	// 	collaborator_id = $1
	// order by
	// 	1`

	queryAllColumns := `select
		date_registration,
		settlement_group,
		calculation_type,
		days_worked,
		hours_worked,
		replace(summa, ' ', '')
	from
		lkr_nachisleniy_zp2020
	where
		collaborator_id = $1
	union all
	select
		date_registration,
		settlement_group,
		calculation_type,
		days_worked,
		hours_worked,
		replace(summa, ' ', '')
	from
		lkr_nachisleniy_zp2021
	where
		collaborator_id = $1
	union all
	select
		date_registration,
		settlement_group,
		calculation_type,
		days_worked,
		hours_worked,
		replace(summa, ' ', '')
	from
		lkr_nachisleniy_zp2022
	where
		collaborator_id = $1
	order by
		1`

	rows, err := DB.Query(queryAllColumns, argsquery...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//ColumnsStructSlice := []store.V1BudgetStatResponds{}
	//V1BudgetStatGroupResponds := store.V1BudgetStatGroupResponds{}

	MapV1BudgetStatGroupResponds := make(map[int]V1BudgetStatGroupResponds)

	// Создать карту  с V1BudgetStatGroupResponds по месяцу.
	// Каждый раз доставать эту карту и обогащать ее данные + добавлять новые значения в массив
	// провести сортировку записей по дате + провести сортироку карты через пакет сорт или через слайз, после сортировки
	// Сделать конечный массив с группировками и выести его в итоги.

	for rows.Next() {
		var r V1BudgetStatResponds

		var DaysWorked, HoursWorked, Summ string
		err = rows.Scan(&r.DateRegistration, &r.SettlementGroup, &r.CalculationType, &DaysWorked, &HoursWorked, &Summ)
		if err != nil {
			return nil, err
		}

		if Summ == "0" && DaysWorked == "0" && HoursWorked == "0" {
			continue
		}

		re := regexp.MustCompile(`\d{2}.\d{2}.\d{4}`)
		date_from_subjectArray := re.FindAllString(r.DateRegistration, -1)
		if len(date_from_subjectArray) == 0 {
			log.Impl.Errorf("Ошибка парсинга даты для метода: %s collaborator_id: %s\n", "V1AverageSalaryGeneral", WorkerID)
			continue
		}
		//fmt.Printf("%q\n", date_from_subjectArray)

		date_from_subject, err := time.Parse("02.01.2006", date_from_subjectArray[0])
		if err != nil {
			return nil, err
		}
		//fmt.Println(date_from_subject)

		compareData := date_from_subject.Before(subtractYear)
		if compareData {
			continue
		}

		//yearArg, monthArg, dayArg := time.Now().Date()
		MonthAccruals := date_from_subject.Month()
		//fmt.Println(int(month))
		//V1BudgetStatGroupResponds.Month = int(MonthAccruals)

		var V1BudgetStatGroupResponds V1BudgetStatGroupResponds
		V1BudgetStatGroupResponds, _ = MapV1BudgetStatGroupResponds[int(MonthAccruals)]
		// if !ok {
		// 	var V1BudgetStatGroupResponds store.V1BudgetStatGroupResponds
		// 	MapV1BudgetStatGroupResponds[int(MonthAccruals)] = V1BudgetStatGroupResponds
		// }

		V1BudgetStatGroupResponds.Month = int(MonthAccruals)

		Summ = strings.Replace(Summ, ",", ".", -1)
		Summ = strings.Replace(Summ, " ", "", -1)

		//HoursWorked = strings.Replace(HoursWorked, "\n", "", -1)
		SummFloat, err := strconv.ParseFloat(Summ, 32)
		if err != nil {
			return nil, err
		}

		r.Summ = float32(SummFloat)

		if r.SettlementGroup == "Начислено" {
			//r.Summ = 10
		} else {
			//r.Summ = -10
			r.Summ = -r.Summ
		}

		V1BudgetStatGroupResponds.Total = V1BudgetStatGroupResponds.Total + r.Summ
		if r.SettlementGroup == "Начислено" {
			V1BudgetStatGroupResponds.TotalGross = V1BudgetStatGroupResponds.TotalGross + r.Summ
		} else {
			V1BudgetStatGroupResponds.TotalDeduction = V1BudgetStatGroupResponds.TotalDeduction + -r.Summ
		}

		if UseYearFilter {
			// TODO: Подключить по возможности к регулярному выражению ниже.
			Ok := strings.Contains(r.DateRegistration, yearFilter)
			if !Ok {
				continue
			}
		}

		DaysWorkedInt, err := strconv.Atoi(DaysWorked)
		if err != nil {
			return nil, err
		}
		V1BudgetStatGroupResponds.DaysWorked = DaysWorkedInt

		//fmt.Println(HoursWorked)
		//fmt.Println(777)
		//HoursWorked = strings.Replace(HoursWorked, `\`, "", -1)
		//HoursWorked = strings.ReplaceAll(HoursWorked, string([]byte{92, 114, 92, 110}), "")
		//re := regexp.MustCompile(`\r?\n`)
		//HoursWorked = re.ReplaceAllString(HoursWorked, "")
		//HoursWorked = strings.TrimSpace(HoursWorked)

		HoursWorked = strings.Replace(HoursWorked, ",", ".", -1)

		//HoursWorked = strings.Replace(HoursWorked, "\n", "", -1)
		HoursWorkedFloat, err := strconv.ParseFloat(HoursWorked, 32)
		if err != nil {
			return nil, err
		}

		V1BudgetStatGroupResponds.HoursWorked = float32(HoursWorkedFloat)
		//ColumnsStructSlice = append(ColumnsStructSlice, r)
		V1BudgetStatGroupResponds.Items = append(V1BudgetStatGroupResponds.Items, r)

		MapV1BudgetStatGroupResponds[int(MonthAccruals)] = V1BudgetStatGroupResponds
	}

	var V1BudgetStatGroupRespondsSlice []V1BudgetStatGroupResponds

	countMonth := 0
	var countSum float32
	var countSumGross float32

	keys := make([]int, 0, len(MapV1BudgetStatGroupResponds))
	for k := range MapV1BudgetStatGroupResponds {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		V1BudgetStatGroupRespondsSlice = append(V1BudgetStatGroupRespondsSlice, MapV1BudgetStatGroupResponds[k])

		countMonth++
		countSum = countSum + MapV1BudgetStatGroupResponds[k].Total
		countSumGross = countSumGross + MapV1BudgetStatGroupResponds[k].TotalGross
	}

	var AverageSalary AverageSalary
	if countMonth != 0 {
		AverageSalary.Months = countMonth
		AverageSalary.Summ = countSum
		AverageSalary.SummGross = countSumGross
		AverageSalary.Average = countSum / float32(countMonth)
		AverageSalary.AverageGross = countSumGross / float32(countMonth)
		AverageSalary.DaySum = (countSum / float32(countMonth)) / 29.3
		AverageSalary.DaySum = float32(math.Ceil(float64(AverageSalary.DaySum)*100) / 100)

		AverageSalary.DaySumGross = (countSumGross / float32(countMonth)) / 29.3
		AverageSalary.DaySumGross = float32(math.Ceil(float64(AverageSalary.DaySumGross)*100) / 100)
	}

	//fmt.Printf("Month = %d, Summ = %f Average = %f\n", AverageSalary.Months, AverageSalary.Summ, AverageSalary.Average)

	var AnswerWebV1 AnswerWebV1
	AnswerWebV1.Status = true
	//AnswerWebV1.Data = ColumnsStructSlice
	AnswerWebV1.Data = AverageSalary
	AnswerWebV1.Error = nil
	//c.JSON(http.StatusOK, AnswerWebV1)
	return AnswerWebV1, nil
}
