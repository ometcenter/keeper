package models

type TableDescriptionDSN struct {
	TableName string   `json:"tableName"`
	DSN       string   `json:"dsn"`
	QuerySQL  string   `json:"query"`
	Fields    []Fields `json:"fields"`
}

type Fields struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	NotNull    bool   `json:"notNull"`
	PrimaryKey bool   `json:"primaryKey"`
	TypeChange string `json:"typeChangeEvent"`
}

type ColumnsStruct struct {
	ColumnName string `json:"columnName"`
	DataType   string `json:"dataType"`
	IsNullable string `json:"isNullable"`
	PrimaryKey bool   `json:"primaryKey"`
}

type IndexesDescription struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
