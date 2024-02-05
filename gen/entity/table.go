package entity

import "february/gen/pkg/tools"

// Table db tables info
type Table struct {
	TableName    string `json:"tableName"`
	TableComment string `json:"tableComment"`

	CacheName          string    `gorm:"-" json:"cacheName"`
	ClassName          string    `gorm:"-" json:"className"`
	LowerCaseClassName string    `gorm:"-" json:"lowerCaseClassName"`
	HasTime            bool      `gorm:"-" json:"HasTime"`
	Columns            []*Column `gorm:"-" json:"columns"`
}

func (t *Table) InitCacheName(enableCache bool) {
	if enableCache {
		t.CacheName = tools.CreateCacheName(t.ClassName)
	}
}

func (t *Table) InitClassName(tablePrefix string) {
	t.ClassName = tools.FormatStructName(tablePrefix, t.TableName)
}

func (t *Table) InitLowerCaseClassName(tablePrefix string) {
	t.LowerCaseClassName = tools.FormatJsonColumn(tablePrefix, t.TableName)
}

func (t *Table) InitHasTime() {
	for _, c := range t.Columns {
		if c.ColType == "time.Time" {
			t.HasTime = true
			break
		}
	}
}
