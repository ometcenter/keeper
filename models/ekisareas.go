package models

import "gorm.io/gorm"

type EkisAreas struct {
	Ekis             string
	Inn              string
	ShortName        string
	FullName         string
	Mrsd             string
	AreaBgu          string
	AreaZkgu         string
	AddressBgu       string
	AddressZkgu      string
	BaseBgu          string
	BaseZkgu         string
	District         string
	Region           string
	Director         string
	Site             string
	MailDirector     string
	MailOrganization string
	Cgu              string
	Updated          string
	AreaUnf          string
	AddressUnf       string
	BaseUnf          string
}

type EkisTokenSession struct {
	Status    bool   `json:"status"`
	ID        string `json:"id"`
	ExpiresIn int    `json:"expiresIn"`
}

type EkisOrganizationDesctiption struct {
	gorm.Model
	EoId                         int    `json:"eo_id"`
	FullName                     string `json:"full_name"`
	ShortName                    string `json:"short_name"`
	Number                       int    `json:"number"`
	TypeId                       int    `json:"type_id"`
	InPreschoolStaffing          int    `json:"in_preschool_staffing"`
	Type                         string `json:"type"`
	Type2Id                      int    `json:"type2_id"`
	Type2                        string `json:"type2"`
	ClassId                      int    `json:"class_id"`
	Class                        string `json:"class"`
	Inn                          string `json:"inn"`
	Ogrn                         string `json:"ogrn"`
	Kpp                          string `json:"kpp"`
	Founder                      string `json:"founder"`
	PropertyTypeId               int    `json:"property_type_id"`
	PropertyType                 string `json:"property_type"`
	LegalOrganizationId          int    `json:"legal_organization_id"`
	LegalOrganization            string `json:"legal_organization"`
	SubordinationId              int    `json:"subordination_id"`
	Subordination                string `json:"subordination"`
	Director                     string `json:"director"`
	XaIsActive                   string `json:"xa_is_active"`
	Post                         string `json:"post"`
	LegalAddress                 string `json:"legal_address"`
	PublicPhone                  string `json:"public_phone"`
	Email                        string `json:"email"`
	Website                      string `json:"website"`
	HasPreschoolEdu              string `json:"has_preschool_edu"`
	DayCareEdu                   string `json:"day_care_edu"`
	EoDistrictId                 int    `json:"eo_district_id"`
	EoDistrictName               string `json:"eo_district_name"`
	Mrsd                         string `json:"mrsd"`
	MrsdId                       string `json:"mrsd_id"`
	MrsdPredsedatel              string `json:"mrsd_predsedatel"`
	Rating                       string `json:"rating"`
	Arhiv                        int    `json:"arhiv"`
	StatusId                     int    `json:"status_id"`
	Status                       string `json:"status"`
	MunicipalName                string `json:"municipal_name"`
	EoInFirstClass               int    `json:"eo_in_first_class"`
	EoInRecordUdo                int    `json:"eo_in_record_udo"`
	EoWorkCityCamp               int    `json:"eo_work_city_camp"`
	LegalOrganizationIdAfter2013 string `json:"legal_organization_id_after2013"`
	IsDonm                       int    `json:"is_donm"`
	IsCgu                        string `json:"is_cgu"`
	//DataUpdate                   time.Time
}

type EkisOrganizationDesctiptionRespons struct {
	Status bool                          `json:"is_cgu"`
	Data   []EkisOrganizationDesctiption `json:"data"`
}

type EkisOrganizationAddresses struct {
	gorm.Model
	EoId                          int    `json:"eo_id"`            // Номер организации ЕКИС
	Unom                          int    `json:"unom"`             // Уникальный номер статкарты БТИ
	Unad                          int    `json:"unad"`             // UNAD
	District                      string `json:"district"`         // Муниципальный округ (Район)
	AreaArea                      string `json:"area"`             // Административный округ
	Address                       string `json:"address"`          // Адрес
	AddressAsur                   string `json:"address_asur"`     // Адрес (другой формат)
	IsMainBuilding                string `json:"is_main_building"` // Признак главного здания
	AdrLng                        string `json:"adr_lng"`          // X_center
	AdrLat                        string `json:"adr_lat"`          // Y_center
	Fias                          string `json:"fias"`
	FiasAddressDadata             string `json:"address_fias_dadata"`
	FiasAddressUnrestrictedDadata string `json:"address_fias_unrestricted_dadata"`
	IsTempAccom                   string `json:"is_temp_accom"` // Временное размещение
	TempEnd                       string `json:"temp_end"`      // Дата окончания временного размещения
	FullName                      string //full_name
	ShortName                     string //short_name
	Inn                           string
	Number                        int    //number
	XaIsActive                    string //xa_is_active
}

type EkisOrganizationAddressesRespons struct {
	Status bool                        `json:"status"`
	Data   []EkisOrganizationAddresses `json:"data"`
}
