package entity

// Table db tables info
type Table struct {
	TableName    string `json:"tableName"`
	TableComment string `json:"tableComment"`

	ClassName          string    `gorm:"-" json:"className"`
	LowerCaseClassName string    `gorm:"-" json:"lowerCaseClassName"`
	HasTime            bool      `gorm:"-" json:"HasTime"`
	Columns            []*Column `gorm:"-" json:"columns"`
}
