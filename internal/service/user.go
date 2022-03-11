package service

import (
	"crypto/sha1"
	"fmt"
	"grpc/internal/entity"
	"grpc/internal/repository"
)
const (
	salt = "alskdoaskdpal"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Create(user *entity.User) (uint64, error) {
	user.Password = s.generatePasswordHash(user.Password)
	userId, err := s.repo.User.Create(user)
	if err == nil {
		s.repo.User.FlushCachedUsers()
	}
	err = s.repo.LogUser(user)
	if err != nil {
		return 0, err
	}
	return userId, err

}

func (s *UserService) GetAll() ([]entity.User, error) {

	hasCache := s.repo.User.HasCachedUsers()
	if hasCache {
		users, err := s.repo.User.GetCachedUsers()
		return users, err
	}

	users, err := s.repo.User.GetAll()
	s.repo.User.CacheUsers(users)
	return users, err
}

func (s *UserService) Delete(userId uint64) error {
	return s.repo.User.Delete(userId)
}

func (s *UserService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
