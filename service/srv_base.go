package service

import (
	"ai-report/common"
	"ai-report/config/log"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseService[T any] struct {
	ctx context.Context
}

func NewService[T any](ctx context.Context) *BaseService[T] {
	return &BaseService[T]{ctx: ctx}
}

// Insert 插入
func (b *BaseService[T]) Insert(entity T, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Create(entity).Error; err != nil {
		log.ErrorF(b.ctx, "insert err:: ", zap.Error(err))
		return err
	}
	return nil
}

// Update 更新
func (b *BaseService[T]) Update(entity T, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Save(entity).Error; err != nil {
		log.ErrorF(b.ctx, "update err:: ", zap.Error(err))
		return err
	}
	return nil
}

//// DeleteById 根据id删除
//func (b *BaseService[T]) DeleteById(id uint64, tx *gorm.DB) error {
//	if tx == nil {
//		tx = common.Ormx
//	}
//	if err := tx.Where("id = ?", id).Delete(T{}).Error; err != nil {
//		log.ErrorF(b.ctx, "delete by id err: %v", zap.Error(err))
//		return err
//	}
//	return nil
//}

// FindById 根据id查询
func (b *BaseService[T]) FindById(id uint64) (*T, error) {
	var t T
	if err := common.Ormx.WithContext(b.ctx).Where("id = ?", id).First(&t).Error; err != nil {
		log.ErrorF(b.ctx, "find by id err:: ", zap.Error(err))
		return nil, err
	}
	return &t, nil
}

// FindList 查询列表
func (b *BaseService[T]) FindList(condition func(where ...interface{}) *gorm.DB) ([]T, error) {
	var list []T
	db := common.Ormx.WithContext(b.ctx)
	if condition != nil {
		db = condition(db)
	}
	if err := db.Find(&list).Error; err != nil {
		log.ErrorF(b.ctx, "find list err:: ", zap.Error(err))
		return nil, err
	}
	return list, nil
}
