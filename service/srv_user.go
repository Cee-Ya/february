package service

import (
	"ai-report/entity"
	"context"
)

type UserService struct {
	*BaseService[entity.User]
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{NewService[entity.User](ctx, &entity.User{})}
}
