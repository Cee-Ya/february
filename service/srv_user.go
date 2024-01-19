package service

import (
	"ai-report/entity"
	"context"
)

type UserService struct {
	BaseService[entity.User]
}

func NewUserService(ctx context.Context) *UserService {
	userSvr := &UserService{*NewService[entity.User](ctx)}
	return userSvr
}
