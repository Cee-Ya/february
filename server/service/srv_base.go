package service

import (
	"context"
	"february/common"
	"february/entity"
	"february/pkg/logx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var errStr = "BaseService err:: "

type BaseServiceInterface interface {
	TableName() string
}

// BaseService 基础服务
// T: 实体类型
// 用于实现基础的增删改查
type BaseService[T any] struct {
	BaseServiceInterface
	ctx context.Context
}

func NewService[T any](ctx context.Context, entity BaseServiceInterface) *BaseService[T] {
	return &BaseService[T]{ctx: ctx, BaseServiceInterface: entity}
}

// Insert 插入
func (b *BaseService[T]) Insert(entity T, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Create(&entity).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"insert err:: ", zap.Error(err))
		return err
	}
	return nil
}

// Update 更新
// entity中的所有字段都会更新，所以如果需要修改某个字段，需要先查询出来，再修改
func (b *BaseService[T]) Update(entity T, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Save(entity).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"update not null err:: ", zap.Error(err))
		return err
	}
	return nil
}

// UpdateAttr 更新
func (b *BaseService[T]) UpdateAttr(id uint64, attrs map[string]interface{}, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Table(b.TableName()).Where("id = ?", id).Updates(attrs).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"update err:: ", zap.Error(err))
		return err
	}
	return nil

}

// DeleteById 根据id删除
func (b *BaseService[T]) DeleteById(id uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Where("id = ?", id).Delete(b.TableName()).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"delete by id err:: ", zap.Error(err))
		return err
	}
	return nil
}

// Delete 删除
func (b *BaseService[T]) Delete(entity T, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Delete(entity).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"delete err: %v", zap.Error(err))
		return err
	}
	return nil
}

// FindById 根据id查询
func (b *BaseService[T]) FindById(id uint64) (t *T, err error) {
	if err = common.Ormx.WithContext(b.ctx).Where("id = ?", id).Limit(1).Find(&t).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"find by id err:: ", zap.Error(err))
		return
	}
	return
}

// FindOne 查询单个
func (b *BaseService[T]) FindOne(condition func(where *gorm.DB)) (t *T, err error) {
	db := common.Ormx.WithContext(b.ctx).Table(b.TableName())
	if condition != nil {
		condition(db)
	}
	if err = db.Limit(1).Find(&t).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"find one err:: ", zap.Error(err))
		return
	}
	return
}

// FindList 查询列表
func (b *BaseService[T]) FindList(condition func(where *gorm.DB)) (ts []T, err error) {
	db := common.Ormx.WithContext(b.ctx).Table(b.TableName())
	if condition != nil {
		condition(db)
	}
	if err = db.Find(&ts).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"find list err:: ", zap.Error(err))
		return
	}
	return
}

// FindPageList 分页查询
func (b *BaseService[T]) FindPageList(condition func(where *gorm.DB), page *entity.Page) (res *entity.PageResult[T], err error) {
	res = &entity.PageResult[T]{}
	db := common.Ormx.WithContext(b.ctx).Table(b.TableName())
	if condition != nil {
		condition(db)
	} else {
		//全局默认按id倒序
		db.Order("id desc")
	}
	if err = db.Count(&res.Total).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"find page count err:: ", zap.Error(err))
		return
	}
	if err = db.Offset((page.PageNo - 1) * page.PageSize).Limit(page.PageSize).Find(&res.Row).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"find page list err:: ", zap.Error(err))
		return
	}
	return
}
