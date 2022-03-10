package service

import (
	"grpc/internal/repository"
	"grpc/transport/handler"
)

type Service struct {
	User handler.UserService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repo),
	}
}
