package service

import (
	"february/common"
	"february/common/tools"
	"february/entity"
	"february/server/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	*BaseService[entity.User]
}

func NewUserService(ctx *gin.Context) *UserService {
	c := common.GetTraceCtx(ctx)
	return &UserService{NewService[entity.User](c, &entity.User{})}
}

func (u *UserService) PageList(page *entity.Page) (*entity.PageResult[vo.UserPageVo], error) {
	temp, err := u.FindPageList(func(where *gorm.DB) {
		where.Order("id desc")
	}, page)
	if err != nil {
		return nil, err
	}
	var res = make([]vo.UserPageVo, 0)
	for _, v := range temp.Rows {
		var userPageVo vo.UserPageVo
		if err = copier.Copy(&userPageVo, &v); err != nil {
			u.Error(zap.Error(errors.Wrap(err, "copy err:")))
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
		return err
	}
	if temp.ID > 0 {
		return errors.New("user already exists")
	}

	return u.orm.Transaction(func(tx *gorm.DB) error {
		var pass string
		pass, err = tools.HashPassword(add.Password)
		if err != nil {
			u.Error(zap.Error(errors.Wrap(err, "hash password err:")))
			return err
		}
		user := entity.User{
			Username: add.Username,
			Password: pass,
			Mobile:   add.Mobile,
			Backup:   add.Backup,
		}
		if err = u.Insert(user, tx); err != nil {
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
	return u.orm.Transaction(func(tx *gorm.DB) error {
		var err error
		// 判断用户是否存在
		if _, err = u.FindById(update.Id); err != nil {
			return err
		}
		var temp entity.User
		if err = copier.Copy(&temp, &update); err != nil {
			u.Error(zap.Error(errors.Wrap(err, "copy err:")))
			return err
		}
		if update.Password != "" {
			var pass string
			if pass, err = tools.HashPassword(update.Password); err != nil {
				u.Error(zap.Error(errors.Wrap(err, "hash password err:")))
				return err
			}
			temp.Password = pass
		}
		if err = u.ModifyNotNull(&temp, tx); err != nil {
			return err
		}
		return nil
	})
}
