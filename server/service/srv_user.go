package service

import (
	"ai-report/common"
	"ai-report/common/tools"
	"ai-report/entity"
	"ai-report/pkg/logx"
	"ai-report/server/vo"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	*BaseService[entity.User]
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{NewService[entity.User](ctx, &entity.User{})}
}

// Create 创建用户
func (u *UserService) Create(add vo.UserAddVo) error {
	// 判断用户是否存在
	var temp *entity.User
	var err error
	if temp, err = u.FindOne(func(where *gorm.DB) {
		where.Where("username = ?", add.Username)
	}); err != nil {
		logx.ErrorF(u.ctx, "user find one err:: ", zap.Error(err))
		return err
	}
	if temp.ID > 0 {
		return errors.New("user already exists")
	}

	return common.Ormx.WithContext(u.ctx).Transaction(func(tx *gorm.DB) error {
		var pass string
		pass, err = tools.HashPassword(add.Password)
		if err != nil {
			logx.ErrorF(u.ctx, "user password hash err:: ", zap.Error(err))
			return err
		}
		user := entity.User{
			Username: add.Username,
			Password: pass,
			Mobile:   add.Mobile,
			Backup:   add.Backup,
		}
		if err = u.Insert(user, tx); err != nil {
			logx.ErrorF(u.ctx, "user insert err:: ", zap.Error(err))
			return err
		}
		return nil
	})
}
