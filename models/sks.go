package models

type HRInfo struct {
	MessageType      string             `json:"MessageType"`
	RequestExecuted  bool               `json:"RequestExecuted"`
	ErrorDescription string             `json:"ErrorDescription"`
	RequestType      string             `json:"RequestType"`
	ResponceBody     ResponceBodyHRInfo `json:"ResponceBody"`
}

type HRInfoRecord struct {
	OrgID             string `json:"OrgID"`
	PersonID          string `json:"PersonID"`
	Name              string `json:"Name"`
	Snils             string `json:"SNILS"`
	HRCaseType        string `json:"HRCaseType"`
	CaseDate          string `json:"CaseDate"`
	Subdivision       string `json:"Subdivision"`
	Position          string `json:"Position"`
	NamePosition      string `json:"NamePosition"`
	NumberOfPositions string `json:"NumberOfPositions"`
	Registrator       string `json:"Registrator"`
	EmploymentTypeID  string `json:"EmploymentTypeID"`
	NumContract       string `json:"NumContract"`
	DataContract      string `json:"DataContract"`
	//MobilePhone       []string `json:"MobilePhone"`
	EmailEPS string `json:"EmailEPS"`
}

type ResponceBodyHRInfo struct {
	HRInfoRecord []HRInfoRecord `json:"HRInfo"`
}

type HRInfoBodyReduest struct {
	OrgIDArray []string `json:"OrgIDArray"`
}

type SksReferenceInfo struct {
	MessageType      string                       `json:"MessageType"`
	RequestExecuted  bool                         `json:"RequestExecuted"`
	ErrorDescription string                       `json:"ErrorDescription"`
	RequestType      string                       `json:"RequestType"`
	ResponceBody     SksReferenceInfoResponceBody `json:"ResponceBody"`
}
type SksReferenceInfoOrganizations struct {
	ID        string `json:"ID"`
	Name      string `json:"Name"`
	Inn       string `json:"INN"`
	Kpp       string `json:"KPP"`
	ChiefFIO  string `json:"ChiefFIO"`
	ShortName string `json:"ShortName"`
}
type SksReferenceInfoResponceBody struct {
	Organizations []SksReferenceInfoOrganizations `json:"Organizations"`
}

type PersonsInfoRequestBody struct {
	PersonIDArray []string                 `json:"PersonIDArray"`
	AttributeList PersonsInfoAttributeList `json:"AttributeList"`
}
type PersonsInfoFilter struct {
	Include []string `json:"Include"`
}
type PersonsInfoAttributeList struct {
	PersonsInfo PersonsInfoFilter `json:"PersonsInfo"`
}

type SksPersonsInfoResponse struct {
	MessageType     string       `json:"MessageType"`
	RequestExecuted bool         `json:"RequestExecuted"`
	RequestType     string       `json:"RequestType"`
	ResponceBody    ResponceBody `json:"ResponceBody"`
}
type SksPersonsInfo struct {
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	FirstName  string `json:"FirstName"`
	MiddleName string `json:"MiddleName"`
	LastName   string `json:"LastName"`
	BirthDate  string `json:"BirthDate"`
	Gender     string `json:"Gender"`
	Inn        string `json:"INN"`
	Snils      string `json:"SNILS"`
	EmailEPS   string `json:"EmailEPS"`
}
type ResponceBody struct {
	PersonsInfo []SksPersonsInfo `json:"PersonsInfo"`
}
