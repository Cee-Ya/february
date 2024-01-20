package service

import (
	"context"
	"errors"
	"february/common"
	"february/common/tools"
	"february/entity"
	"february/pkg/logx"
	"february/server/vo"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	*BaseService[entity.User]
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{NewService[entity.User](ctx, &entity.User{})}
}

func (u *UserService) PageList(page *entity.Page) (*entity.PageResult[vo.UserPageVo], error) {
	temp, err := u.FindPageList(func(where *gorm.DB) {
		where.Order("id desc")
	}, page)
	if err != nil {
		logx.ErrorF(u.ctx, "user page list err:: ", zap.Error(err))
		return nil, err
	}
	var res = make([]vo.UserPageVo, 0)
	for _, v := range temp.Rows {
		var userPageVo vo.UserPageVo
		if err = copier.Copy(&userPageVo, &v); err != nil {
			logx.ErrorF(u.ctx, "user copy err:: ", zap.Error(err))
			return nil, err
		}
		userPageVo.CreateTime = v.CreateTime.ToString()
		userPageVo.ModifyTime = v.ModifyTime.ToString()
		res = append(res, userPageVo)
	}
	return &entity.PageResult[vo.UserPageVo]{Total: temp.Total, Rows: res}, nil
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

// Update 更新用户
func (u *UserService) Update(update vo.UserUpdateVo) error {
	if update.Id == 0 {
		return errors.New("id is required")
	}
	return common.Ormx.WithContext(u.ctx).Transaction(func(tx *gorm.DB) error {
		var (
			user *entity.User
			err  error
		)
		// 判断用户是否存在
		if user, err = u.FindById(update.Id); err != nil {
			logx.ErrorF(u.ctx, "user find by id err:: ", zap.Error(err))
			return err
		}
		if user.ID == 0 {
			return errors.New("user not exists")
		}
		var temp entity.User
		if err = copier.Copy(&temp, &update); err != nil {
			logx.ErrorF(u.ctx, "user copy err:: ", zap.Error(err))
			return err
		}
		if update.Password != "" {
			var pass string
			if pass, err = tools.HashPassword(update.Password); err != nil {
				logx.ErrorF(u.ctx, "user password hash err:: ", zap.Error(err))
				return err
			}
			temp.Password = pass
		}
		if err = u.ModifyNotNull(&temp, tx); err != nil {
			logx.ErrorF(u.ctx, "user update err:: ", zap.Error(err))
			return err
		}
		return nil
	})
}
