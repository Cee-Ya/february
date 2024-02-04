package entity

// Column db tables column info
type Column struct {
	ColumnName    string `json:"columnName"`
	ColumnKey     string `json:"columnKey"`
	ColumnLen     int    `json:"columnLen"`
	DataType      string `json:"dataType"`
	ColumnComment string `json:"columnComment"`

	ColName string `gorm:"-" json:"colName"`
	ColType string `gorm:"-" json:"colType"`
	JsonTag string `gorm:"-" json:"jsonTag"`
}
