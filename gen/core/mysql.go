package core

import (
	"february/gen/entity"
	"february/gen/pkg/tools"
	"gorm.io/gorm"
	"strings"
)

// mysql type map
var mysqlType = map[string]string{
	"tinyint":    "int",
	"smallint":   "int",
	"mediumint":  "int",
	"int":        "int",
	"integer":    "int",
	"bigint":     "int64",
	"float":      "float64",
	"double":     "float64",
	"decimal":    "float64",
	"char":       "string",
	"varchar":    "string",
	"tinytext":   "string",
	"text":       "string",
	"mediumtext": "string",
	"longtext":   "string",
	"tinyblob":   "string",
	"blob":       "string",
	"mediumblob": "string",
	"longblob":   "string",
	"date":       "time.Time",
	"time":       "time.Time",
	"year":       "time.Time",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
}

type mysql struct {
	db *gorm.DB
}

func init() {
	DBData["mysql"] = &mysql{}
}

func (m *mysql) InitConn(db *gorm.DB) {
	m.db = db
}

func (m *mysql) InitColumn(columns []*entity.Column) []*entity.Column {
	for i := range columns {
		columns[i].ColType = mysqlType[columns[i].DataType]
		// 名称改为首字母大写的驼峰命名
		columns[i].ColName = tools.FormatStructName("", columns[i].ColumnName)
		// 格式化json tag
		columns[i].JsonTag = tools.FormatJsonColumn("", columns[i].ColumnName)
	}
	return columns
}

// GetDataBaseName get database name by dsn
// dsn format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func (m *mysql) GetDataBaseName(dsn string) string {
	start := strings.Index(dsn, "/") + 1
	end := strings.Index(dsn, "?")
	if start > 0 && end > 0 && start < end {
		dsn = dsn[start:end]
	} else {
		dsn = ""
	}
	return dsn
}

// GetTable get tables by database name
func (m *mysql) GetTable(db string) (tables []*entity.Table, err error) {
	q := m.db.Raw("select table_name as TableName, table_comment as TableComment from information_schema.tables where table_schema = ?", db).Scan(&tables)
	if q.Error != nil {
		return nil, q.Error
	}
	return tables, nil
}

// GetColumn get columns by table name
func (m *mysql) GetColumn(table string) (column []*entity.Column, err error) {
	if err = m.db.Raw("SELECT column_name as ColumnName, column_key as ColumnKey, data_type as DataType, IFNULL(character_maximum_length,0) as ColumnLen, column_comment as ColumnComment FROM information_schema.columns WHERE table_name = ?", table).Scan(&column).Error; err != nil {
		return nil, err
	}
	return column, nil
}
