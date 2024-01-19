package common

import (
	"ai-report/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var GlobalConfig *entity.Config
var Logger *zap.Logger
var Ormx *gorm.DB
