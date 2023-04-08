package dbrepo

import (
	"testing"
	"time"

	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/util"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) models.User {
	var err error
	email := util.RandomString(6) + "@" + util.RandomString(2) + "." + util.RandomString(2)
	avatar := "http://" + util.RandomString(6)
	name := util.RandomString(4)
	password := util.RandomString(6)
	password, err = util.GenHashedPassword(password)
	require.NoError(t, err)

	user := models.User{
		Email:    email,
		Avatar:   avatar,
		Name:     name,
		Password: password,
	}

	err = testRepo.CreateUser(&user)
	require.NoError(t, err)
	require.True(t, user.ID > 0)
	require.Equal(t, email, user.Email)
	require.Equal(t, avatar, user.Avatar)
	require.Equal(t, name, user.Name)
	require.Equal(t, password, user.Password)
	require.WithinDuration(t, user.CreatedAt, time.Now(), time.Second)
	require.WithinDuration(t, user.UpdatedAt, time.Now(), time.Second)
	require.False(t, user.DeletedAt.Valid)

	return user
}

func TestCreateUser(t *testing.T) {
	_ = createUser(t)
}

func TestGetUserByID(t *testing.T) {
	user := createUser(t)
	newUser, err := testRepo.GetUserByID(user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, user.Avatar, newUser.Avatar)
	require.Equal(t, user.Name, newUser.Name)
	require.Equal(t, user.Password, newUser.Password)
	require.WithinDuration(t, user.CreatedAt, newUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.UpdatedAt, newUser.UpdatedAt, time.Second)
	require.False(t, user.DeletedAt.Valid)
	require.False(t, newUser.DeletedAt.Valid)
}

// 测试 user 创建了 link ，获取 user 是否返回 links
func TestGetUserByIDWithLinks(t *testing.T) {
	user := createUser(t)

	count := 5
	for i := 0; i < count; i++ {
		_ = createLink(t, user.ID)
	}

	qUser, err := testRepo.GetUserByID(user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, qUser)
	require.NotEmpty(t, qUser.Links)
	require.True(t, len(qUser.Links) > 0)

	for _, link := range qUser.Links {
		require.NotEmpty(t, link)
		require.True(t, link.ID > 0)
		require.True(t, len(link.LongURL) > 0)
		require.True(t, len(link.ShortHash) > 0)
	}
}

func TestGetUserByEmail(t *testing.T) {
	user := createUser(t)
	newUser, err := testRepo.GetUserByEmail(user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, user.Avatar, newUser.Avatar)
	require.Equal(t, user.Name, newUser.Name)
	require.Equal(t, user.Password, newUser.Password)
	require.WithinDuration(t, user.CreatedAt, newUser.CreatedAt, time.Second)
	require.WithinDuration(t, newUser.UpdatedAt, time.Now(), time.Second)
	require.False(t, user.DeletedAt.Valid)
	require.False(t, newUser.DeletedAt.Valid)
}

func TestUpdateUserByID(t *testing.T) {
	user := createUser(t)

	name := util.RandomString(6)
	avatar := util.RandomString(6)

	updateUser := models.User{
		Name:   name,
		Avatar: avatar,
	}

	err := testRepo.UpdateUserByID(user.ID, updateUser)
	require.NoError(t, err)

	updateUser, err = testRepo.GetUserByID(user.ID)

	require.NoError(t, err)
	require.Equal(t, avatar, updateUser.Avatar)
	require.Equal(t, name, updateUser.Name)
	require.Equal(t, user.Email, updateUser.Email)
	require.Equal(t, user.Password, updateUser.Password)
	require.WithinDuration(t, user.CreatedAt, updateUser.CreatedAt, time.Second)
	require.WithinDuration(t, updateUser.UpdatedAt, time.Now(), time.Second)
	require.False(t, user.DeletedAt.Valid)
	require.False(t, updateUser.DeletedAt.Valid)

}
