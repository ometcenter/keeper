package models

type SnilsValidityCheckingRequestKeeper struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	BirthDate  string `json:"birthDate"`
	SnilsInfo  string `json:"snilsInfo"`
}

type SnilsValidityCheckingRequestKafka struct {
	SnilsValidityCheckingRequest SnilsValidityCheckingRequest `json:"snils_validity_checking_request"`
}

type SnilsValidityCheckingRequest struct {
	RequestID               string      `json:"request_id"`
	RequestingSystem        string      `json:"requesting_system"`
	RequestedMethod         string      `json:"requested_method"`
	StateServiceVarietyCode string      `json:"state_service_variety_code"`
	Errors                  interface{} `json:"errors"`
	PersonInfo              PersonInfo  `json:"person_info"`
	SnilsInfo               string      `json:"snils_info"`
}

type SnilsGettingRequestKeeper struct {
	LastName           string `json:"lastName"`
	FirstName          string `json:"firstName"`
	MiddleName         string `json:"middleName"`
	BirthDate          string `json:"birthDate"`
	GenderCode         string `json:"genderCode"`
	DocumentTypeCode   string `json:"documentTypeCode"`
	DocumentSeries     string `json:"documentSeries"`
	DocumentNumber     string `json:"documentNumber"`
	DocumentIssueDate  string `json:"documentIssueDate"`
	DocumentIssuerName string `json:"documentIssuerName"`
}

type SnilsGettingRequestKafka struct {
	SnilsGettingRequest SnilsGettingRequest `json:"snils_getting_request"`
}

type DocumentInfoSnilsGetting struct {
	DocumentTypeCode   string `json:"document_type_code"`
	DocumentSeries     string `json:"document_series"`
	DocumentNumber     string `json:"document_number"`
	DocumentIssueDate  string `json:"document_issue_date"`
	DocumentIssuerName string `json:"document_issuer_name"`
}
type SnilsGettingRequest struct {
	RequestID               string                   `json:"request_id"`
	RequestingSystem        string                   `json:"requesting_system"`
	RequestedMethod         string                   `json:"requested_method"`
	StateServiceVarietyCode string                   `json:"state_service_variety_code"`
	PersonInfo              PersonInfoSnilsGetting   `json:"person_info"`
	BirthPlaceInfo          interface{}              `json:"birth_place_info"`
	DocumentInfo            DocumentInfoSnilsGetting `json:"document_info"`
}

type SnilsValidityCheckingResponseKafka struct {
	SnilsValidityCheckingResponse SnilsValidityCheckingResponse `json:"snils_validity_checking_response"`
}
type Errors struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}
type PersonInfo struct {
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	BirthDate  string `json:"birth_date"`
}
type SnilsValidityInfo struct {
	SnilsStatus string `json:"snils_status"`
}
type SnilsValidityCheckingResponse struct {
	RequestID               string     `json:"request_id"`
	RequestingSystem        string     `json:"requesting_system"`
	RequestedMethod         string     `json:"requested_method"`
	StateServiceVarietyCode string     `json:"state_service_variety_code"`
	Errors                  []Errors   `json:"errors"`
	PersonInfo              PersonInfo `json:"person_info"`
	SnilsInfo               string     `json:"snils_info"`
	//SnilsValidityInfo       string     `json:"snils_validity_info"`
	SnilsValidityInfo SnilsValidityInfo `json:"snils_validity_info"`
}

type SnilsGettingResponseResponseKafka struct {
	SnilsGettingResponse SnilsGettingResponse `json:"snils_getting_response"`
}
type PersonInfoSnilsGetting struct {
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	BirthDate  string `json:"birth_date"`
	GenderCode string `json:"gender_code"`
}
type DocumentInfo struct {
	DocumentTypeCode   string `json:"document_type_code"`
	DocumentSeries     string `json:"document_series"`
	DocumentNumber     string `json:"document_number"`
	DocumentIssueDate  string `json:"document_issue_date"`
	DocumentIssuerName string `json:"document_issuer_name"`
}
type SnilsGettingInfo struct {
	SnilsSuccessfulFinding string `json:"snils_successful_finding"`
	Snils                  string `json:"snils"`
}
type SnilsGettingResponse struct {
	RequestID               string           `json:"request_id"`
	RequestingSystem        string           `json:"requesting_system"`
	RequestedMethod         string           `json:"requested_method"`
	StateServiceVarietyCode string           `json:"state_service_variety_code"`
	Errors                  []Errors         `json:"errors"`
	PersonInfo              PersonInfo       `json:"person_info"`
	DocumentInfo            DocumentInfo     `json:"document_info"`
	SnilsGettingInfo        SnilsGettingInfo `json:"snils_getting_info"`
}
