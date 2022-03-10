package repository

import (
	"github.com/go-redis/redis"
	"grpc/internal/entity"
	"strconv"
	"time"
)

type UserRedis struct {
	rdb *redis.Client
}

func NewUserRedis(rdb *redis.Client) *UserRedis {
	return &UserRedis{rdb: rdb}
}

func (r *UserRedis) CacheUsers(users []entity.User) {
	for i := range users {
		userId := strconv.FormatUint(users[i].Id, 10)
		r.rdb.Set(ctx, "users:"+userId+":Id", users[i].Id, time.Minute)
		r.rdb.Set(ctx, "users:"+userId+":Name", users[i].Name, time.Minute)
		r.rdb.Set(ctx, "users:"+userId+":Email", users[i].Email, time.Minute)
		r.rdb.Set(ctx, "users:"+userId+":Password", users[i].Password, time.Minute)
	}
	r.rdb.Set(ctx, "users:cached", true, time.Minute)
}

func (r *UserRedis) GetCachedUsers() ([]entity.User, error) {
	var users []entity.User

	usersKeysSlice, err := r.getAllExistingIds()
	if err != nil {
		return users, err
	}

	for i := range usersKeysSlice {
		strUserId, _ := r.rdb.Get(ctx, usersKeysSlice[i].(string)).Result()
		user := r.getCachedUserById(strUserId)
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRedis) HasCachedUsers() bool {
	b, err := r.rdb.Get(ctx, "users:cached").Bool()
	if err != nil {
		return false
	}
	return b
}
func (r *UserRedis) FlushCachedUsers() {
	r.rdb.Set(ctx, "users:cached", false, time.Minute)
	usersKeysSlice, err := r.getAllExistingIds()
	if err != nil {
		return
	}
	for i := range usersKeysSlice {
		strUserId, _ := r.rdb.Get(ctx, usersKeysSlice[i].(string)).Result()
		r.deleteCachedUserById(strUserId)
	}

}

func (r *UserRedis) getCachedUserById(strUserId string) entity.User {
	var user entity.User
	user.Id, _ = r.rdb.Get(ctx, "users:"+strUserId+":Id").Uint64()
	user.Name, _ = r.rdb.Get(ctx, "users:"+strUserId+":Name").Result()
	user.Email, _ = r.rdb.Get(ctx, "users:"+strUserId+":Email").Result()
	user.Password, _ = r.rdb.Get(ctx, "users:"+strUserId+":Password").Result()
	return user
}

func (r *UserRedis) deleteCachedUserById(strUserId string) {
	r.rdb.Del(ctx, "users:"+strUserId+":Id")
	r.rdb.Del(ctx, "users:"+strUserId+":Name")
	r.rdb.Del(ctx, "users:"+strUserId+":Email")
	r.rdb.Del(ctx, "users:"+strUserId+":Password")
}

func (r *UserRedis) getAllExistingIds() ([]interface{}, error) {
	usersKeys, err := r.rdb.Do(ctx, "KEYS", "users:*:Id").Result()
	return usersKeys.([]interface{}), err
}
