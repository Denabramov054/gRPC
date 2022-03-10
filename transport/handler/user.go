package handler

import (
	"context"
	"grpc/internal/entity"
	"grpc/pkg/api"
)
type UserHandler struct {
	service UserService
	api.UnimplementedUserServer
}

type UserService interface {
	Create(user *entity.User) (uint64, error)
	GetAll() ([]entity.User, error)
	Delete(userId uint64) error
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (s *UserHandler) Create(ctx context.Context, req *api.CreateRequest) (*api.UserResponse, error) {
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	userId, err := s.service.Create(user)
	return &api.UserResponse{
		Id:       userId,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, err
}

func (s *UserHandler) Users(context.Context, *api.Empty) (*api.UsersResponse, error) {
	users, err := s.service.GetAll()
	response := &api.UsersResponse{Users: []*api.UserResponse{}}
	for i := range users {
		response.Users = append(response.Users, &api.UserResponse{
			Id:       users[i].Id,
			Name:     users[i].Name,
			Email:    users[i].Email,
			Password: users[i].Password,
		})
	}
	return response, err
}

func (s *UserHandler) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	userId := req.GetId()
	err := s.service.Delete(userId)
	response := &api.DeleteResponse{
		Message: "OK",
	}
	if err != nil {
		response.Message = "Err"
	}
	return response, err
}
