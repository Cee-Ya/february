package core

import (
	"february/gen/entity"
	"february/gen/pkg/tools"
	"gorm.io/gorm"
	"strings"
)

// postgres type map
var postgresType = map[string]string{
	"int2":      "int",
	"int4":      "int4",
	"int8":      "int64",
	"float4":    "float64",
	"float8":    "float64",
	"numeric":   "float64",
	"char":      "string",
	"varchar":   "string",
	"text":      "string",
	"bytea":     "string",
	"date":      "time.Time",
	"time":      "time.Time",
	"timestamp": "time.Time",
	"bigint":    "int64",
	"integer":   "int",
	"bool":      "bool",
}

type postgres struct {
	db *gorm.DB
}

func init() {
	DBData["postgres"] = &postgres{}
}

func (m *postgres) InitConn(db *gorm.DB) {
	m.db = db
}

func (m *postgres) InitColumn(columns []*entity.Column) []*entity.Column {
	for i := range columns {
		columns[i].ColType = postgresType[columns[i].ColumnKey]
		// 名称改为首字母大写的驼峰命名
		columns[i].ColName = tools.FormatStructName("", columns[i].ColumnName)
		// 格式化json tag
		columns[i].JsonTag = tools.FormatJsonColumn("", columns[i].ColumnName)
	}
	return columns
}

// GetDataBaseName get database name by dsn
// dsn format: postgres: host=%s port=%s user=%s dbname=%s password=%s sslmode=%s
func (m *postgres) GetDataBaseName(dsn string) string {
	start := strings.Index(dsn, "dbname=") + 1
	end := strings.Index(dsn, " ")
	if start > 0 && end > 0 && start < end {
		dsn = dsn[start:end]
	} else {
		dsn = ""
	}
	return dsn
}

// GetTable get tables
// desc pg 数据库目前还没找到通过数据库名称查询表的方法  不过在基础配置中已经确定了数据库  可以直接查询
func (p *postgres) GetTable(db string) (tables []*entity.Table, err error) {
	sql := `SELECT
			relname AS TableName,
			obj_description ( relfilenode ) AS TableComment
		FROM
			pg_class
		WHERE
			relkind = 'r'
			AND relnamespace = ( SELECT oid FROM pg_namespace WHERE nspname = 'public' )
			ORDER BY relname`
	//pg 数据库 9.4通过上面代码无法查询表名，可用以下替代
	//	sql := `SELECT table_name
	//FROM information_schema.tables
	//WHERE table_schema = 'public'`
	q := p.db.Raw(sql).Scan(&tables)
	if q.Error != nil {
		return nil, q.Error
	}
	return tables, nil
}

// GetColumn get columns by table name
func (p *postgres) GetColumn(table string) (column []*entity.Column, err error) {
	sql := `SELECT a.attname AS column_name,
			   pg_catalog.format_type(a.atttypid, NULL) AS column_key,
			   CASE
				   WHEN a.atttypid IN (1042, 1043) THEN a.atttypmod - 4
				   WHEN a.atttypid IN (1005, 1016) THEN 0
				   ELSE a.attlen
			   END AS column_len,
			   col_description(a.attrelid, a.attnum) AS column_comment,
			   CASE WHEN pk.constraint_name IS NOT NULL
				   THEN 'Yes'
				   ELSE ''
			   END AS data_type
		FROM pg_attribute a
		LEFT JOIN (
			SELECT kcu.column_name, tc.constraint_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
			WHERE tc.table_name = ? AND tc.constraint_type = 'PRIMARY KEY'
		) pk ON a.attname = pk.column_name
		WHERE a.attrelid = (SELECT oid FROM pg_class WHERE relname = ? LIMIT 1) AND a.attnum > 0 AND NOT a.attisdropped`
	q := p.db.Raw(sql, table, table).Scan(&column)
	if q.Error != nil {
		return nil, q.Error
	}
	return column, nil
}

func (p *postgres) GetColumnLen(table string) (column []*entity.Column, err error) {
	sql := `SELECT a.attlen AS column_len,
			   pg_catalog.format_type(a.atttypid, NULL) AS column_key,
			   CASE
				   WHEN a.atttypid IN (1042, 1043) THEN a.atttypmod - 4
				   WHEN a.atttypid IN (1005, 1016) THEN 0
				   ELSE a.attlen
			   END AS column_len,
			   col_description(a.attrelid, a.attnum) AS column_comment,
			   CASE WHEN pk.constraint_name IS NOT NULL
				   THEN 'Yes'
				   ELSE 'No'
			   END AS data_type
		FROM pg_attribute a
		LEFT JOIN (
			SELECT kcu.column_name, tc.constraint_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name
			WHERE tc.table_name = ? AND tc.constraint_type = 'PRIMARY KEY'
		) pk ON a.attname = pk.column_name
		WHERE a.attrelid = (SELECT oid FROM pg_class WHERE relname = ? LIMIT 1) AND a.attnum > 0 AND NOT a.attisdropped`
	q := p.db.Raw(sql, table, table).Scan(&column)
	if q.Error != nil {
		return nil, q.Error
	}
	return column, nil
}
