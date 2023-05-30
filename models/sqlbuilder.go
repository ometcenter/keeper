package models

type TableDescriptionDSN struct {
	TableName string   `json:"tableName"`
	DSN       string   `json:"dsn"`
	Fields    []Fields `json:"fields"`
}

type Fields struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	NotNull    bool   `json:"notNull"`
	PrimaryKey bool   `json:"primaryKey"`
	TypeChange string `json:"typeChangeEvent"`
}
