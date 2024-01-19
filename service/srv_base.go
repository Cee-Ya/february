package service

import (
	"ai-report/common"
	"ai-report/config/log"
	"ai-report/entity"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var errStr = "BaseService err:: "

type BaseServiceInterface interface {
	TableName() string
}

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
	if err := tx.WithContext(b.ctx).Create(entity).Error; err != nil {
		log.ErrorF(b.ctx, errStr+"insert err:: ", zap.Error(err))
		return err
	}
	return nil
}

// UpdateNotNull 更新
func (b *BaseService[T]) UpdateNotNull(entity T, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Save(entity).Error; err != nil {
		log.ErrorF(b.ctx, errStr+"update not null err:: ", zap.Error(err))
		return err
	}
	return nil
}

// Update 更新
func (b *BaseService[T]) Update(id uint64, attrs map[string]interface{}, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Table(b.TableName()).Where("id = ?", id).Updates(attrs).Error; err != nil {
		log.ErrorF(b.ctx, errStr+"update err:: ", zap.Error(err))
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
		log.ErrorF(b.ctx, errStr+"delete by id err:: ", zap.Error(err))
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
		log.ErrorF(b.ctx, errStr+"delete err: %v", zap.Error(err))
		return err
	}
	return nil
}

// FindById 根据id查询
func (b *BaseService[T]) FindById(id uint64) (t *T, err error) {
	if err = common.Ormx.WithContext(b.ctx).Where("id = ?", id).First(&t).Error; err != nil {
		log.ErrorF(b.ctx, errStr+"find by id err:: ", zap.Error(err))
		return
	}
	return
}

// FindOne 查询单个
func (b *BaseService[T]) FindOne(condition func(where ...interface{}) *gorm.DB) (t *T, err error) {
	db := common.Ormx.WithContext(b.ctx)
	if condition != nil {
		db = condition(db)
	}
	if err = db.First(&t).Error; err != nil {
		log.ErrorF(b.ctx, errStr+"find one err:: ", zap.Error(err))
		return
	}
	return
}

// FindList 查询列表
func (b *BaseService[T]) FindList(condition func(where ...interface{}) *gorm.DB) (ts []T, err error) {
	db := common.Ormx.WithContext(b.ctx)
	if condition != nil {
		db = condition(db)
	}
	if err = db.Find(&ts).Error; err != nil {
		log.ErrorF(b.ctx, errStr+"find list err:: ", zap.Error(err))
		return
	}
	return
}

// FindPageList 分页查询
func (b *BaseService[T]) FindPageList(condition func(where ...interface{}) *gorm.DB, page *entity.Page) (res *entity.PageResult[T], err error) {
	db := common.Ormx.WithContext(b.ctx).Table(b.TableName())
	if condition != nil {
		db = condition(db)
	}
	if err = db.Count(&res.Total).Error; err != nil {
		log.ErrorF(b.ctx, errStr+"find page count err:: ", zap.Error(err))
		return
	}
	if err = db.Offset((page.PageNo - 1) * page.PageSize).Limit(page.PageSize).Find(&res.Row).Error; err != nil {
		log.ErrorF(b.ctx, errStr+"find page list err:: ", zap.Error(err))
		return
	}
	return
}
