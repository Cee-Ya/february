package common

import (
	"february/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var GlobalConfig *entity.Config
var Logger *zap.Logger
var Ormx *gorm.DB
var Redisx entity.Redis
