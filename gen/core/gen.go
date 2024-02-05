package core

import (
	"bytes"
	"february/gen/entity"
	"february/gen/pkg/conf"
	"february/gen/pkg/ormx"
	"february/gen/pkg/tools"
	"fmt"
	"gorm.io/gorm"
	"html/template"
)

var (
	DBData = map[string]GenInterface{}
)

type GenInterface interface {
	InitConn(db *gorm.DB)
	GetDataBaseName(dsn string) string
	GetTable(db string) (tables []*entity.Table, err error)
	GetColumn(table string) (column []*entity.Column, err error)
	InitColumn(columns []*entity.Column) []*entity.Column
}

// GenCode generate code to file
func GenCode(tableNameAttr []string) {
	// 1. connect to db
	db, err := ormx.New(conf.C.Database)
	if err != nil {
		panic(fmt.Sprintf("connect to db failed, err: %v", err))
	}
	// 2. get tables
	genSvr := DBData[conf.C.Database.DBType]
	genSvr.InitConn(db)
	tables, err := genSvr.GetTable(genSvr.GetDataBaseName(conf.C.Database.DSN))
	if err != nil {
		panic(fmt.Sprintf("get tables failed, err: %v", err))
	}
	// 3. get columns
	for _, s := range tableNameAttr {
		for index := range tables {
			if tables[index].TableName == s {
				tables[index].InitClassName(conf.C.Database.TablePrefix)
				tables[index].InitLowerCaseClassName(conf.C.Database.TablePrefix)
				tables[index].InitCacheName(conf.C.Gen.EnableCache)
				columns, err := genSvr.GetColumn(tables[index].TableName)
				if err != nil {
					panic(fmt.Sprintf("get columns failed, err: %v", err))
				}
				tables[index].Columns = genSvr.InitColumn(columns)
				// if columns type constains time.Time, add import time
				tables[index].InitHasTime()
				// 4. generate code
				// model
				t1, err := template.ParseFiles(conf.C.Gen.AbsPath + "entity.go.template")
				if err != nil {
					panic(fmt.Sprintf("parse model template failed, err: %v", err))
				}
				// service
				t2, err := template.ParseFiles(conf.C.Gen.AbsPath + "service.go.template")
				if err != nil {
					panic(fmt.Sprintf("parse service template failed, err: %v", err))
				}
				// api
				t3, err := template.ParseFiles(conf.C.Gen.AbsPath + "router.go.template")
				if err != nil {
					panic(fmt.Sprintf("parse router template failed, err: %v", err))
				}

				_ = tools.PathCreate(conf.C.Gen.DomainTargetPath)
				_ = tools.PathCreate(conf.C.Gen.ServiceTargetPath)
				_ = tools.PathCreate(conf.C.Gen.ApiTargetPath)
				var b1 bytes.Buffer
				err = t1.Execute(&b1, tables[index])
				if err != nil {
					panic(fmt.Sprintf("generate model code to buf failed, err: %v", err))
				}

				var b2 bytes.Buffer
				err = t2.Execute(&b2, tables[index])
				if err != nil {
					panic(fmt.Sprintf("generate services code to buf failed, err: %v", err))
				}

				var b3 bytes.Buffer
				err = t3.Execute(&b3, tables[index])
				if err != nil {
					panic(fmt.Sprintf("generate api code to buf failed, err: %v", err))
				}

				// 5. write to file
				tools.FileCreate(b1, conf.C.Gen.DomainTargetPath+"ent_"+tables[index].LowerCaseClassName+".go")
				tools.FileCreate(b2, conf.C.Gen.ServiceTargetPath+"srv_"+tables[index].LowerCaseClassName+".go")
				tools.FileCreate(b3, conf.C.Gen.ApiTargetPath+"router_"+tables[index].LowerCaseClassName+".go")
			}
		}
	}
}
