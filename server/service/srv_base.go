package service

import (
	"context"
	"errors"
	"february/common"
	"february/common/tools"
	"february/entity"
	"february/pkg/logx"
	"february/pkg/redis/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var errStr = "BaseService err:: "

type BaseServiceInterface interface {
	TableName() string //表名
	EnableRedis() bool //是否开启缓存
	CacheKey() string  //缓存key
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

// Modify  更新
// entity中的所有字段都会更新，所以如果需要修改某个字段，需要先查询出来，再修改
func (b *BaseService[T]) Modify(entity *T, tx *gorm.DB) error {
	if b.EnableRedis() {
		return b.ModifyByCache(entity, tx)
	} else {
		return b.ModifyBase(entity, tx)
	}
}

func (b *BaseService[T]) ModifyBase(entity *T, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Save(entity).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"update not null err:: ", zap.Error(err))
		return err
	}
	return nil
}

func (b *BaseService[T]) ModifyByCache(entity *T, tx *gorm.DB) error {
	if err := b.ModifyBase(entity, tx); err != nil {
		logx.ErrorF(b.ctx, errStr+"ModifyBase err:: ", zap.Error(err))
		return err
	}
	if err := b.putCacheByEntity(entity); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCache err:: ", zap.Error(err))
		return err
	}
	return nil
}

// ModifyNotNull 更新
// entity中的所有字段都会更新，所以如果需要修改某个字段，需要先查询出来，再修改
func (b *BaseService[T]) ModifyNotNull(entity *T, tx *gorm.DB) error {
	if b.EnableRedis() {
		return b.ModifyNotNullByCache(entity, tx)
	} else {
		return b.ModifyNotNullBase(entity, tx)
	}
}

func (b *BaseService[T]) ModifyNotNullBase(entity *T, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Updates(entity).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"update not null err:: ", zap.Error(err))
		return err
	}
	return nil
}

func (b *BaseService[T]) ModifyNotNullByCache(entity *T, tx *gorm.DB) error {
	if err := b.ModifyNotNullBase(entity, tx); err != nil {
		logx.ErrorF(b.ctx, errStr+"ModifyNotNullBase err:: ", zap.Error(err))
		return err
	}
	if err := b.putCacheByEntity(entity); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCache err:: ", zap.Error(err))
		return err
	}
	return nil
}

// ModifyAttr 更新
func (b *BaseService[T]) ModifyAttr(id uint64, attrs map[string]interface{}, tx *gorm.DB) error {
	if b.EnableRedis() {
		return b.ModifyAttrByCache(id, attrs, tx)
	} else {
		return b.ModifyAttrBase(id, attrs, tx)
	}
}

func (b *BaseService[T]) ModifyAttrBase(id uint64, attrs map[string]interface{}, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Table(b.TableName()).Where("id = ?", id).Updates(attrs).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"update err:: ", zap.Error(err))
		return err
	}
	return nil
}

func (b *BaseService[T]) ModifyAttrByCache(id uint64, attrs map[string]interface{}, tx *gorm.DB) error {
	if err := b.ModifyAttrBase(id, attrs, tx); err != nil {
		logx.ErrorF(b.ctx, errStr+"ModifyAttrBase err:: ", zap.Error(err))
		return err
	}
	if err := b.putCache(id); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCache err:: ", zap.Error(err))
		return err
	}
	return nil
}

// DeleteById 根据id删除
func (b *BaseService[T]) DeleteById(id uint64, tx *gorm.DB) error {
	if b.EnableRedis() {
		return b.DeleteByIdAndCache(id, tx)
	} else {
		return b.DeleteByIdBase(id, tx)
	}
}

func (b *BaseService[T]) DeleteByIdBase(id uint64, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Where("id = ?", id).Delete(b.TableName()).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"delete by id err:: ", zap.Error(err))
		return err
	}
	return nil
}

func (b *BaseService[T]) DeleteByIdAndCache(id uint64, tx *gorm.DB) error {
	if err := b.DeleteByIdBase(id, tx); err != nil {
		logx.ErrorF(b.ctx, errStr+"DeleteByIdBase err:: ", zap.Error(err))
		return err
	}
	if err := b.delCache(id); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCache err:: ", zap.Error(err))
		return err
	}
	return nil
}

// Delete 删除
func (b *BaseService[T]) Delete(entity T, tx *gorm.DB) error {
	if b.EnableRedis() {
		return b.DeleteByCache(entity, tx)
	} else {
		return b.DeleteBase(entity, tx)
	}
}

func (b *BaseService[T]) DeleteBase(entity T, tx *gorm.DB) error {
	if tx == nil {
		tx = common.Ormx
	}
	if err := tx.WithContext(b.ctx).Delete(entity).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"delete err: %v", zap.Error(err))
		return err
	}
	return nil
}

func (b *BaseService[T]) DeleteByCache(entity T, tx *gorm.DB) error {
	if err := b.DeleteBase(entity, tx); err != nil {
		logx.ErrorF(b.ctx, errStr+"DeleteBase err:: ", zap.Error(err))
		return err
	}
	if err := b.delCacheByEntity(&entity); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCache err:: ", zap.Error(err))
		return err
	}
	return nil
}

// FindById 根据id查询
func (b *BaseService[T]) FindById(id uint64) (t *T, err error) {
	if b.EnableRedis() {
		return b.FindByIdCache(id)
	} else {
		return b.FindByIdBase(id)
	}
}

func (b *BaseService[T]) FindByIdBase(id uint64) (t *T, err error) {
	if err = common.Ormx.WithContext(b.ctx).Where("id = ?", id).Limit(1).Find(&t).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"find by id err:: ", zap.Error(err))
		return
	}
	return
}

func (b *BaseService[T]) FindByIdCache(id uint64) (t *T, err error) {
	if t, err = b.getCacheByCheck(id); err == nil && t != nil {
		return
	}
	if err != nil {
		return
	}
	if t, err = b.FindByIdBase(id); err != nil {
		logx.ErrorF(b.ctx, errStr+"FindByIdBase err:: ", zap.Error(err))
		return
	}
	if err = b.initCache(id, t); err != nil {
		logx.ErrorF(b.ctx, errStr+"initCache err:: ", zap.Error(err))
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
	if err = db.Offset((page.PageNo - 1) * page.PageSize).Limit(page.PageSize).Find(&res.Rows).Error; err != nil {
		logx.ErrorF(b.ctx, errStr+"find page list err:: ", zap.Error(err))
		return
	}
	return
}

// ############################################################################################################
// ################################################ 缓存部分  ###################################################
// ############################################################################################################
func (b *BaseService[T]) cacheCheck() (err error) {
	if !b.EnableRedis() {
		return errors.New("redis not enable")
	}
	if b.CacheKey() == "" {
		logx.ErrorF(b.ctx, errStr+"CacheKey is empty")
		return
	}
	return

}

func (b *BaseService[T]) initCache(key any, t *T) (err error) {
	if err = b.cacheCheck(); err != nil {
		if err.Error() == "redis not enable" {
			err = nil
		} else {
			return
		}
	}
	// 如果已存在缓存，则直接返回
	var temp *T
	if temp, err = b.getCache(key); err == nil && temp != nil {
		return
	}
	if err = b.setCache(key, t); err != nil {
		logx.ErrorF(b.ctx, errStr+"initCache err:: ", zap.Error(err))
		return
	}
	return
}

func (b *BaseService[T]) putCache(key any) (err error) {
	var t *T
	if t, err = b.FindByIdBase(key.(uint64)); err != nil {
		logx.ErrorF(b.ctx, errStr+"putCache FindByIdBase err:: ", zap.Error(err))
		return
	}
	if err = b.setCache(key, t); err != nil {
		logx.ErrorF(b.ctx, errStr+"putCache err:: ", zap.Error(err))
		return
	}
	return
}

func (b *BaseService[T]) putCacheByEntity(t *T) (err error) {
	if err = b.cacheCheck(); err != nil {
		if err.Error() == "redis not enable" {
			err = nil
		} else {
			return
		}
	}
	// 通过反射获取id
	var id uint64
	var tempKey interface{}
	if tempKey, err = tools.GetStructField(*t, "ID"); err != nil {
		logx.ErrorF(b.ctx, errStr+"putCacheByEntity GetStructField err:: ", zap.Error(err))
		return
	}
	if id, err = tools.Any2Uint64(tempKey); err != nil {
		logx.ErrorF(b.ctx, errStr+"putCacheByEntity Any2Uint64 err:: ", zap.Error(err))
		return
	}
	if err = b.putCache(id); err != nil {
		logx.ErrorF(b.ctx, errStr+"putCacheByEntity err:: ", zap.Error(err))
		return
	}
	return
}

func (b *BaseService[T]) delCacheByEntity(t *T) (err error) {
	if err = b.cacheCheck(); err != nil {
		if err.Error() == "redis not enable" {
			err = nil
		} else {
			return
		}
	}
	// 通过反射获取id
	var id uint64
	var tempKey interface{}
	if tempKey, err = tools.GetStructField(*t, "ID"); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCacheByEntity GetStructField err:: ", zap.Error(err))
		return
	}
	if id, err = tools.Any2Uint64(tempKey); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCacheByEntity Any2Uint64 err:: ", zap.Error(err))
		return
	}
	if err = b.delCache(id); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCacheByEntity err:: ", zap.Error(err))
		return
	}
	return

}

func (b *BaseService[T]) setCache(key any, t *T) (err error) {
	var str string
	str, err = tools.ToJson(t)
	if err != nil {
		logx.ErrorF(b.ctx, errStr+"setCache ToJson err:: ", zap.Error(err))
		return
	}
	var myKey string
	if myKey, err = tools.ToString(key); err != nil {
		logx.ErrorF(b.ctx, errStr+"setCache ToString err:: ", zap.Error(err))
		return
	}
	if err = utils.NewRedisUtils(b.ctx).MustSet(b.CacheKey()+myKey, str); err != nil {
		logx.ErrorF(b.ctx, errStr+"setCache MustSet err:: ", zap.Error(err))
		return
	}
	return
}

func (b *BaseService[T]) getCacheByCheck(key any) (t *T, err error) {
	if err = b.cacheCheck(); err != nil {
		if err.Error() == "redis not enable" {
			err = nil
		} else {
			return
		}
	}
	if t, err = b.getCache(key); err != nil {
		logx.ErrorF(b.ctx, errStr+"getCache err:: ", zap.Error(err))
		return
	}
	return
}

func (b *BaseService[T]) getCache(key any) (t *T, err error) {
	var myKey string
	if myKey, err = tools.ToString(key); err != nil {
		logx.ErrorF(b.ctx, errStr+"getCache ToString err:: ", zap.Error(err))
		return
	}
	var str string
	if str = utils.NewRedisUtils(b.ctx).MustGet(b.CacheKey() + myKey); str == "" {
		return
	}
	if err = tools.Str2Struct(str, &t); err != nil {
		logx.ErrorF(b.ctx, errStr+"getCache ToStruct err:: ", zap.Error(err))
		return
	}
	return
}

func (b *BaseService[T]) delCache(key any) (err error) {
	if err = b.cacheCheck(); err != nil {
		if err.Error() == "redis not enable" {
			err = nil
		} else {
			return
		}
	}
	var myKey string
	if myKey, err = tools.ToString(key); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCache ToString err:: ", zap.Error(err))
		return
	}
	if err = utils.NewRedisUtils(b.ctx).MustDel(b.CacheKey() + myKey); err != nil {
		logx.ErrorF(b.ctx, errStr+"delCache MustDel err:: ", zap.Error(err))
		return
	}
	return
}
